package p_uniquity_finance_invoices

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
	finance_accounts "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_accounts"
	finance_customer "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_customer"
	finance_products "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_products"
	finance_taxes "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_taxes"
)

const (
	financeAccountsMainMenuInvoicesLinkKey            = "finance_invoices.FinanceAccountsMainMenuLink"
	financeAccountsMainMenuInvoicePaymentTermsLinkKey = "finance_invoices.FinanceAccountsMainMenuPaymentTermsLink"
	financeAccountsMainMenuPostedInvoicesLinkKey      = "finance_invoices.FinanceAccountsMainMenuPostedLink"
	financeAccountsMainMenuCancelledInvoicesLinkKey   = "finance_invoices.FinanceAccountsMainMenuCancelledLink"
)

func sidebarMenuHasChildKey(menu *components.SidebarMenu, key string) bool {
	for _, ch := range menu.Children {
		if item, ok := ch.(*components.SidebarMenuItem); ok && item.GetKey() == key {
			return true
		}
	}
	return false
}

func patchFinanceAccountsMainMenuForInvoices(page components.PageInterface) components.PageInterface {
	menu, ok := page.(*components.SidebarMenu)
	if !ok {
		panic("p_uniquity_finance_invoices: finance_accounts.MainMenu must be *components.SidebarMenu")
	}
	newChildren := append([]components.PageInterface{}, menu.Children...)
	add := func(key, title, routeKey, icon string) {
		if !sidebarMenuHasChildKey(menu, key) {
			newChildren = append(newChildren, &components.SidebarMenuItem{
				Page:  components.Page{Key: key, Roles: []string{"superuser"}},
				Title: getters.Static(title),
				Url:   lamu.RoutePath(routeKey, nil),
				Icon:  icon,
			})
		}
	}
	add(financeAccountsMainMenuInvoicesLinkKey, "Draft invoices", "finance_invoices.DefaultRoute", "document-text")
	add(financeAccountsMainMenuPostedInvoicesLinkKey, "Posted invoices", "finance_invoices.PostedInvoiceListRoute", "document-check")
	add(financeAccountsMainMenuCancelledInvoicesLinkKey, "Cancelled invoices", "finance_invoices.CancelledInvoiceListRoute", "document-minus")
	add(financeAccountsMainMenuInvoicePaymentTermsLinkKey, "Payment terms", "finance_invoices.PaymentTermListRoute", "calendar-days")
	if len(newChildren) == len(menu.Children) {
		return menu
	}
	cloned := *menu
	cloned.Children = newChildren
	return &cloned
}

func pluginPages() lamu.PluginFeatures[components.PageInterface] {
	e := pageEntriesDraftInvoicePages()
	e = append(e, pageEntriesPostedInvoicePages()...)
	e = append(e, pageEntriesCancelledInvoicePages()...)
	e = append(e, pageEntriesInvoiceFilterPage()...)
	e = append(e, pageEntriesPaymentTermPages()...)
	e = append(e, pageEntriesPaymentTermFkSelectPages()...)
	return lamu.PluginFeatures[components.PageInterface]{
		Entries: e,
		Patches: []registry.Pair[string, func(components.PageInterface) components.PageInterface]{
			{Key: "finance_accounts.MainMenu", Value: patchFinanceAccountsMainMenuForInvoices},
		},
	}
}

func invoiceDatetimeStringGetter(ctxKey string) getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		t, err := getters.Key[time.Time](ctxKey)(ctx)
		if err != nil {
			return "", err
		}
		if t.IsZero() {
			return "", nil
		}
		return t.Format(time.RFC3339), nil
	}
}

func invoicePaymentTermFKDisplayGetter() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		id, err := getters.Key[uint]("$in.ID")(ctx)
		if err != nil || id == 0 {
			return "", nil
		}
		typ, err := getters.Key[string]("$in.Type")(ctx)
		if err != nil {
			return fmt.Sprintf("#%d", id), nil
		}
		bid, _ := getters.Key[uint]("$in.BackingID")(ctx)
		inst, err := ResolvePaymentTermInstance(ctx, typ, bid)
		if err != nil {
			return fmt.Sprintf("#%d", id), nil
		}
		return fmt.Sprintf("#%d — %s", id, inst.Summary()), nil
	}
}

func invoiceCreateDatetimeGetter() getters.Getter[time.Time] {
	return func(ctx context.Context) (time.Time, error) {
		t, err := getters.Key[time.Time]("$in.Datetime")(ctx)
		if err != nil {
			return time.Time{}, nil
		}
		return t, nil
	}
}

func draftNumberGetter(ctxKey string) getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		v, err := getters.Key[*string](ctxKey)(ctx)
		if err != nil || v == nil {
			return "", nil
		}
		return *v, nil
	}
}

func invoiceDetailPaymentTermSummaryGetter() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		typ, err := getters.Key[string]("$in.PaymentTerm.Type")(ctx)
		if err != nil {
			return "", err
		}
		bid, err := getters.Key[uint]("$in.PaymentTerm.BackingID")(ctx)
		if err != nil {
			return "", err
		}
		id, err := getters.Key[uint]("$in.PaymentTerm.ID")(ctx)
		if err != nil {
			return "", err
		}
		inst, err := ResolvePaymentTermInstance(ctx, typ, bid)
		if err != nil {
			return fmt.Sprintf("#%d", id), nil
		}
		return fmt.Sprintf("#%d — %s", id, inst.Summary()), nil
	}
}

func invoiceDetailTaxesNamesGetter() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		m, ok := ctx.Value("$in").(map[string]any)
		if !ok || m == nil {
			return "—", nil
		}
		raw, ok := m["Taxes"]
		if !ok || raw == nil {
			return "—", nil
		}
		taxes, ok := raw.([]finance_taxes.Tax)
		if !ok || len(taxes) == 0 {
			return "—", nil
		}
		names := make([]string, 0, len(taxes))
		for _, t := range taxes {
			names = append(names, t.Name)
		}
		return strings.Join(names, ", "), nil
	}
}

func invoiceProductChoices() getters.Getter[[]registry.Pair[uint, string]] {
	return func(ctx context.Context) ([]registry.Pair[uint, string], error) {
		db, err := getters.DBFromContext(ctx)
		if err != nil {
			return nil, err
		}
		var products []finance_products.Product
		if err := db.WithContext(ctx).Order("name asc").Find(&products).Error; err != nil {
			return nil, err
		}
		out := make([]registry.Pair[uint, string], 0, len(products))
		for _, p := range products {
			out = append(out, registry.Pair[uint, string]{Key: p.ID, Value: p.Name})
		}
		return out, nil
	}
}

func invoiceLinesDraftJSONGetter() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		if v := ctx.Value(getters.ContextKeyIn); v != nil {
			if m, ok := v.(map[string]any); ok {
				if raw, ok := m["InvoiceLinesJSON"].(string); ok && strings.TrimSpace(raw) != "" {
					return raw, nil
				}
				if raw, ok := m["PendingLines"]; ok && raw != nil {
					b, err := json.Marshal(raw)
					if err == nil && len(b) > 2 && string(b) != "null" {
						return string(b), nil
					}
				}
			}
		}
		return `[{"product_id":0,"quantity":"1","rate":""}]`, nil
	}
}

func invoiceDraftUpdateLinesDefaultsGetter() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		d, err := getters.Key[DraftInvoice]("draft_invoice")(ctx)
		if err != nil || len(d.Lines) == 0 {
			return `[{"product_id":0,"quantity":"1","rate":""}]`, nil
		}
		pending := make([]DraftLinePending, 0, len(d.Lines))
		for _, ln := range d.Lines {
			pending = append(pending, DraftLinePending{
				ProductID: ln.ProductID,
				Quantity:  ln.Quantity.String(),
				Rate:      ln.Rate.String(),
			})
		}
		b, err := json.Marshal(pending)
		if err != nil {
			return `[{"product_id":0,"quantity":"1","rate":""}]`, nil
		}
		return string(b), nil
	}
}

func draftInvoiceLinesDisplayGetter() getters.Getter[[]InvoiceLineDisplay] {
	return func(ctx context.Context) ([]InvoiceLineDisplay, error) {
		m, ok := ctx.Value("$in").(map[string]any)
		if !ok || m == nil {
			return nil, nil
		}
		raw, ok := m["Lines"]
		if !ok || raw == nil {
			return nil, nil
		}
		lines, ok := raw.([]DraftInvoiceLine)
		if !ok || len(lines) == 0 {
			return nil, nil
		}
		out := make([]InvoiceLineDisplay, 0, len(lines))
		for _, ln := range lines {
			name := ln.Product.Name
			if name == "" {
				name = fmt.Sprintf("#%d", ln.ProductID)
			}
			out = append(out, InvoiceLineDisplay{
				Product:  name,
				Quantity: ln.Quantity.String(),
				Rate:     ln.Rate.String(),
			})
		}
		return out, nil
	}
}

func postedInvoiceLinesDisplayGetter() getters.Getter[[]InvoiceLineDisplay] {
	return func(ctx context.Context) ([]InvoiceLineDisplay, error) {
		m, ok := ctx.Value("$in").(map[string]any)
		if !ok || m == nil {
			return nil, nil
		}
		raw, ok := m["Lines"]
		if !ok || raw == nil {
			return nil, nil
		}
		lines, ok := raw.([]PostedInvoiceLine)
		if !ok || len(lines) == 0 {
			return nil, nil
		}
		out := make([]InvoiceLineDisplay, 0, len(lines))
		for _, ln := range lines {
			name := ln.Product.Name
			if name == "" {
				name = fmt.Sprintf("#%d", ln.ProductID)
			}
			out = append(out, InvoiceLineDisplay{
				Product:  name,
				Quantity: ln.Quantity.String(),
				Rate:     ln.Rate.String(),
			})
		}
		return out, nil
	}
}

func cancelledInvoiceLinesDisplayGetter() getters.Getter[[]InvoiceLineDisplay] {
	return func(ctx context.Context) ([]InvoiceLineDisplay, error) {
		m, ok := ctx.Value("$in").(map[string]any)
		if !ok || m == nil {
			return nil, nil
		}
		raw, ok := m["Lines"]
		if !ok || raw == nil {
			return nil, nil
		}
		lines, ok := raw.([]CancelledInvoiceLine)
		if !ok || len(lines) == 0 {
			return nil, nil
		}
		out := make([]InvoiceLineDisplay, 0, len(lines))
		for _, ln := range lines {
			name := ln.Product.Name
			if name == "" {
				name = fmt.Sprintf("#%d", ln.ProductID)
			}
			out = append(out, InvoiceLineDisplay{
				Product:  name,
				Quantity: ln.Quantity.String(),
				Rate:     ln.Rate.String(),
			})
		}
		return out, nil
	}
}

func draftInvoiceCreateEditInputs() []components.PageInterface {
	return []components.PageInterface{
		&components.ContainerError{
			Error: getters.Key[error]("$error.Number"),
			Children: []components.PageInterface{
				&components.InputText{
					Name:     "Number",
					Label:    "Invoice number (optional; assigned on post if empty)",
					Required: false,
					Getter:   draftNumberGetter("$in.Number"),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Datetime"),
			Children: []components.PageInterface{
				&components.InputDatetime{Label: "Invoice date & time", Name: "Datetime", Required: true, Getter: invoiceCreateDatetimeGetter()},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.CustomerID"),
			Children: []components.PageInterface{
				&components.InputForeignKey[finance_customer.Customer]{
					Label:       "Customer",
					Name:        "CustomerID",
					Required:    true,
					Url:         lamu.RoutePath("finance_customers.CustomerFkSelectRoute", nil),
					Display:     getters.Key[string]("$in.Name"),
					Placeholder: "Select customer…",
					Getter:      getters.Association[finance_customer.Customer, uint](getters.Key[uint]("$in.CustomerID")),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.PaymentTermID"),
			Children: []components.PageInterface{
				&components.InputForeignKey[PaymentTerm]{
					Label:       "Payment term",
					Name:        "PaymentTermID",
					Required:    true,
					Url:         lamu.RoutePath("finance_invoices.PaymentTermFkSelectRoute", nil),
					Display:     invoicePaymentTermFKDisplayGetter(),
					Placeholder: "Select payment term…",
					Getter:      getters.Association[PaymentTerm, uint](getters.Key[uint]("$in.PaymentTermID")),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.AccountReceivableID"),
			Children: []components.PageInterface{
				&components.InputForeignKey[finance_accounts.Account]{
					Label:       "Accounts receivable",
					Name:        "AccountReceivableID",
					Required:    true,
					Url:         lamu.RoutePath("finance_accounts.AccountSelectRoute", nil),
					Display:     getters.Format("%s (#%d)", getters.Any(getters.Key[string]("$in.Name")), getters.Any(getters.Key[uint]("$in.ID"))),
					Placeholder: "Select…",
					Getter:      getters.Association[finance_accounts.Account, uint](getters.Key[uint]("$in.AccountReceivableID")),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.AccountRevenueID"),
			Children: []components.PageInterface{
				&components.InputForeignKey[finance_accounts.Account]{
					Label:       "Revenue account",
					Name:        "AccountRevenueID",
					Required:    true,
					Url:         lamu.RoutePath("finance_accounts.AccountSelectRoute", nil),
					Display:     getters.Format("%s (#%d)", getters.Any(getters.Key[string]("$in.Name")), getters.Any(getters.Key[uint]("$in.ID"))),
					Placeholder: "Select…",
					Getter:      getters.Association[finance_accounts.Account, uint](getters.Key[uint]("$in.AccountRevenueID")),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.AccountTaxPayableID"),
			Children: []components.PageInterface{
				&components.InputForeignKey[finance_accounts.Account]{
					Label:       "Tax payable",
					Name:        "AccountTaxPayableID",
					Required:    true,
					Url:         lamu.RoutePath("finance_accounts.AccountSelectRoute", nil),
					Display:     getters.Format("%s (#%d)", getters.Any(getters.Key[string]("$in.Name")), getters.Any(getters.Key[uint]("$in.ID"))),
					Placeholder: "Select…",
					Getter:      getters.Association[finance_accounts.Account, uint](getters.Key[uint]("$in.AccountTaxPayableID")),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.JournalID"),
			Children: []components.PageInterface{
				&components.InputForeignKey[finance_accounts.Journal]{
					Label:       "Journal",
					Name:        "JournalID",
					Required:    true,
					Url:         lamu.RoutePath("finance_accounts.JournalListRoute", nil),
					Display:     getters.Key[string]("$in.Name"),
					Placeholder: "Select journal…",
					Getter:      getters.Association[finance_accounts.Journal, uint](getters.Key[uint]("$in.JournalID")),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Taxes"),
			Children: []components.PageInterface{
				&components.InputManyToMany[finance_taxes.Tax]{
					Label:       "Taxes",
					Name:        "Taxes",
					Display:     getters.Key[string]("$in.Name"),
					Url:         lamu.RoutePath("finance_taxes.TaxMultiSelectRoute", nil),
					Placeholder: "Select taxes…",
					Classes:     "w-full",
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.InvoiceLinesJSON"),
			Children: []components.PageInterface{
				&InputInvoiceLinesDraft{
					Page:     components.Page{Key: "finance_invoices.DraftInvoiceCreateForm.Lines"},
					Label:    "Lines",
					Name:     "InvoiceLinesJSON",
					Choices:  invoiceProductChoices(),
					Defaults: invoiceLinesDraftJSONGetter(),
					Classes:  "w-full",
				},
			},
		},
	}
}

func draftInvoiceUpdateFormInputs() []components.PageInterface {
	inputs := draftInvoiceCreateEditInputs()
	for i := range inputs {
		if box, ok := inputs[i].(*components.ContainerError); ok {
			if inp, ok2 := box.Children[0].(*InputInvoiceLinesDraft); ok2 {
				inp.Defaults = invoiceDraftUpdateLinesDefaultsGetter()
				break
			}
		}
	}
	return inputs
}

func invoiceListPaymentTermSummaryGetter() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		pt, err := getters.Key[PaymentTerm]("$row.PaymentTerm")(ctx)
		if err != nil {
			return "", err
		}
		inst, err := ResolvePaymentTermInstanceFromTerm(ctx, &pt)
		if err != nil {
			return pt.Type, nil
		}
		return inst.Summary(), nil
	}
}

func pageEntriesDraftInvoicePages() []registry.Pair[string, components.PageInterface] {
	createName := getters.Static("finance_invoices.DraftInvoiceCreateForm")
	updateName := getters.Static("finance_invoices.DraftInvoiceUpdateForm")
	deleteName := getters.Static("finance_invoices.DraftInvoiceDeleteForm")

	return []registry.Pair[string, components.PageInterface]{
		{Key: "finance_invoices.DraftInvoiceTable", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_accounts.MainMenu"}},
			Children: []components.PageInterface{
				&components.ContainerColumn{
					Classes: "gap-2 mb-2 w-full",
					Children: []components.PageInterface{
						&components.Environment[uint]{
							Label:   "Fiscal year",
							Key:     getters.Static(FinanceInvoicesEnvironmentFiscalYearKey),
							Options: fiscalYearsEnvironmentOptionsGetter,
							Default: invoiceFiscalYearEnvironmentDefault,
							Classes: "w-full",
						},
					},
				},
				&components.DataTable[DraftInvoice]{
					UID:     "finance-draft-invoice-table",
					Classes: "w-full",
					Data:    getters.Key[components.ObjectList[DraftInvoice]]("draft_invoices"),
					Actions: []components.PageInterface{
						&components.TableButtonFilter{Child: lamu.DynamicPage{Name: "finance_invoices.InvoiceFilter"}},
						&components.TableButtonCreate{
							Link: lamu.RoutePath("finance_invoices.DraftInvoiceCreateRoute", nil),
							Page: components.Page{Roles: []string{"superuser"}},
						},
					},
					RowAttr: getters.RowAttrNavigate(lamu.RoutePath("finance_invoices.DraftInvoiceDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("$row.ID")),
					})),
					Columns: []components.TableColumn{
						{Label: "Number", Name: "Number", Children: []components.PageInterface{
							&components.FieldText{Getter: draftNumberOrPlaceholderRow("$row.Number")},
						}},
						{Label: "Datetime", Name: "Datetime", Children: []components.PageInterface{
							&components.FieldText{Getter: invoiceDatetimeStringGetter("$row.Datetime")},
						}},
						{Label: "Customer", Name: "Customer", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Customer.Name")},
						}},
						{Label: "Payment term", Name: "PaymentTerm", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Format("#%d — %s",
								getters.Any(getters.Key[uint]("$row.PaymentTermID")),
								getters.Any(invoiceListPaymentTermSummaryGetter()),
							)},
						}},
					},
				},
			},
		}},
		{Key: "finance_invoices.DraftInvoiceDetail", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_invoices.DraftInvoiceDetailMenu"}},
			Children: []components.PageInterface{
				&components.ContainerError{
					Error: getters.Key[error]("$error._global"),
					Children: []components.PageInterface{
						&components.Detail[DraftInvoice]{
							Getter: getters.Key[DraftInvoice]("draft_invoice"),
							Children: []components.PageInterface{
								&components.ContainerColumn{
									Classes: "p-4",
									Children: []components.PageInterface{
										&components.LabelInline{Title: "Number", Children: []components.PageInterface{
											&components.FieldText{Getter: draftNumberOrPlaceholderDetail("$in.Number")},
										}},
										&components.LabelInline{Title: "Invoice date", Children: []components.PageInterface{
											&components.FieldText{Getter: invoiceDatetimeStringGetter("$in.Datetime")},
										}},
										&components.LabelInline{Title: "Customer", Children: []components.PageInterface{
											&components.FieldText{Getter: getters.Key[string]("$in.Customer.Name")},
										}},
										&components.LabelInline{Title: "Payment term", Children: []components.PageInterface{
											&components.FieldText{Getter: invoiceDetailPaymentTermSummaryGetter()},
										}},
										&components.LabelInline{Title: "Taxes", Children: []components.PageInterface{
											&components.FieldText{Getter: invoiceDetailTaxesNamesGetter()},
										}},
										&components.LabelInline{Title: "Lines", Children: []components.PageInterface{
											&FieldInvoiceLines{Getter: draftInvoiceLinesDisplayGetter()},
										}},
									},
								},
							},
						},
					},
				},
			},
		}},
		{Key: "finance_invoices.DraftInvoiceDetailMenu", Value: &components.SidebarMenu{
			Title: getters.Format("Draft %s", getters.Any(draftNumberOrPlaceholderMenu("draft_invoice.Number"))),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("Draft invoices"),
				Url:   lamu.RoutePath("finance_invoices.DefaultRoute", nil),
			},
			Children: []components.PageInterface{
				&components.SidebarMenuItem{
					Title: getters.Static("Detail"),
					Url: lamu.RoutePath("finance_invoices.DraftInvoiceDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("draft_invoice.ID")),
					}),
				},
				&components.SidebarMenuItem{
					Page:  components.Page{Roles: []string{"superuser"}},
					Title: getters.Static("Edit"),
					Url: lamu.RoutePath("finance_invoices.DraftInvoiceUpdateRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("draft_invoice.ID")),
					}),
				},
				&components.SidebarMenuItem{
					Page:  components.Page{Roles: []string{"superuser"}},
					Title: getters.Static("Post invoice"),
					Url: lamu.RoutePath("finance_invoices.DraftInvoicePostRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("draft_invoice.ID")),
					}),
				},
			},
		}},
		{Key: "finance_invoices.DraftInvoicePostForm", Value: &components.ShellScaffold{
			Page:    components.Page{Roles: []string{"superuser"}},
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_invoices.DraftInvoiceDetailMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name: getters.Static("finance_invoices.DraftInvoicePostFormInner"),
					ActionURL: lamu.RoutePath("finance_invoices.DraftInvoicePostRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("draft_invoice.ID")),
					}),
					Children: []components.PageInterface{
						&components.FormComponent[struct{}]{
							Title:         "Post invoice",
							Subtitle:      "Creates the journal entry and posted invoice. This cannot be undone except by cancellation.",
							ChildrenInput: []components.PageInterface{},
							ChildrenAction: []components.PageInterface{
								&components.ButtonSubmit{Label: "Post"},
							},
						},
					},
				},
			},
		}},
		{Key: "finance_invoices.DraftInvoiceCreateForm", Value: &components.ShellScaffold{
			Page:    components.Page{Roles: []string{"superuser"}},
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_accounts.MainMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name:      createName,
					ActionURL: lamu.RoutePath("finance_invoices.DraftInvoiceCreateRoute", nil),
					Children: []components.PageInterface{
						&components.FormComponent[DraftInvoice]{
							Attr:          getters.FormBubbling(createName),
							Title:         "Create draft invoice",
							Subtitle:      "Customer, accounts, journal, lines, and taxes",
							ChildrenInput: draftInvoiceCreateEditInputs(),
							ChildrenAction: []components.PageInterface{
								&components.ButtonSubmit{Label: "Save"},
							},
						},
					},
				},
			},
		}},
		{Key: "finance_invoices.DraftInvoiceUpdateForm", Value: &components.ShellScaffold{
			Page:    components.Page{Roles: []string{"superuser"}},
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_invoices.DraftInvoiceDetailMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name: updateName,
					ActionURL: lamu.RoutePath("finance_invoices.DraftInvoiceUpdateRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("draft_invoice.ID")),
					}),
					Children: []components.PageInterface{
						&components.FormComponent[DraftInvoice]{
							Getter:        getters.Key[DraftInvoice]("draft_invoice"),
							Attr:          getters.FormBubbling(updateName),
							Title:         "Edit draft invoice",
							Subtitle:      "Update header and lines",
							ChildrenInput: draftInvoiceUpdateFormInputs(),
							ChildrenAction: []components.PageInterface{
								&components.ContainerRow{
									Classes: "flex flex-wrap justify-between gap-2 mt-2 items-center",
									Children: []components.PageInterface{
										&components.ContainerRow{
											Classes: "flex justify-end gap-2",
											Children: []components.PageInterface{
												&components.ButtonSubmit{Label: "Update"},
												&components.ButtonModalForm{
													Page:        components.Page{Roles: []string{"superuser"}},
													Label:       "Delete",
													Icon:        "trash",
													Name:        deleteName,
													Url:         lamu.RoutePath("finance_invoices.DraftInvoiceDeleteRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("draft_invoice.ID"))}),
													FormPostURL: lamu.RoutePath("finance_invoices.DraftInvoiceDeleteRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("draft_invoice.ID"))}),
													ModalUID:    "finance-draft-invoice-delete-modal",
													Classes:     "btn-error",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}},
		{Key: "finance_invoices.DraftInvoiceDeleteForm", Value: &components.Modal{
			Page: components.Page{Roles: []string{"superuser"}},
			UID:  "finance-draft-invoice-delete-modal",
			Children: []components.PageInterface{
				&components.DeleteConfirmation{
					Title:   "Delete draft invoice?",
					Message: "This removes the draft and its lines.",
					Attr:    getters.FormBubbling(getters.Key[string]("$get.name")),
				},
			},
		}},
	}
}

func draftNumberOrPlaceholderRow(ctxKey string) getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		s, err := draftNumberOrDash(ctxKey)(ctx)
		if err != nil {
			return "", err
		}
		if s == "—" {
			return "(auto on post)", nil
		}
		return s, nil
	}
}

func draftNumberOrPlaceholderDetail(ctxKey string) getters.Getter[string] {
	return draftNumberOrPlaceholderRow(ctxKey)
}

func draftNumberOrPlaceholderMenu(ctxKey string) getters.Getter[string] {
	return draftNumberOrPlaceholderRow(ctxKey)
}

func draftNumberOrDash(ctxKey string) getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		v, err := getters.Key[*string](ctxKey)(ctx)
		if err != nil || v == nil || strings.TrimSpace(*v) == "" {
			return "—", nil
		}
		return *v, nil
	}
}

func journalEntryLinkGetter(jeIDKey string) getters.Getter[string] {
	return lamu.RoutePath("finance_accounts.JournalEntryDetailRoute", map[string]getters.Getter[any]{
		"id": getters.Any(getters.Key[uint](jeIDKey)),
	})
}

func optionalTimePointerGetter(ctxKey string) getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		t, err := getters.Key[*time.Time](ctxKey)(ctx)
		if err != nil || t == nil || t.IsZero() {
			return "—", nil
		}
		return t.Format(time.RFC3339), nil
	}
}

func pageEntriesPostedInvoicePages() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "finance_invoices.PostedInvoiceTable", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_accounts.MainMenu"}},
			Children: []components.PageInterface{
				&components.ContainerColumn{
					Classes: "gap-2 mb-2 w-full",
					Children: []components.PageInterface{
						&components.Environment[uint]{
							Label:   "Fiscal year",
							Key:     getters.Static(FinanceInvoicesEnvironmentFiscalYearKey),
							Options: fiscalYearsEnvironmentOptionsGetter,
							Default: invoiceFiscalYearEnvironmentDefault,
							Classes: "w-full",
						},
					},
				},
				&components.DataTable[PostedInvoice]{
					UID:     "finance-posted-invoice-table",
					Classes: "w-full",
					Data:    getters.Key[components.ObjectList[PostedInvoice]]("posted_invoices"),
					Actions: []components.PageInterface{
						&components.TableButtonFilter{Child: lamu.DynamicPage{Name: "finance_invoices.InvoiceFilter"}},
					},
					RowAttr: getters.RowAttrNavigate(lamu.RoutePath("finance_invoices.PostedInvoiceDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("$row.ID")),
					})),
					Columns: []components.TableColumn{
						{Label: "Number", Name: "Number", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Number")},
						}},
						{Label: "Posted at", Name: "PostedAt", Children: []components.PageInterface{
							&components.FieldText{Getter: optionalTimePointerGetter("$row.PostedAt")},
						}},
						{Label: "Datetime", Name: "Datetime", Children: []components.PageInterface{
							&components.FieldText{Getter: invoiceDatetimeStringGetter("$row.Datetime")},
						}},
						{Label: "Customer", Name: "Customer", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Customer.Name")},
						}},
					},
				},
			},
		}},
		{Key: "finance_invoices.PostedInvoiceDetail", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_invoices.PostedInvoiceDetailMenu"}},
			Children: []components.PageInterface{
				&components.ContainerError{
					Error: getters.Key[error]("$error._global"),
					Children: []components.PageInterface{
						&components.Detail[PostedInvoice]{
							Getter: getters.Key[PostedInvoice]("posted_invoice"),
							Children: []components.PageInterface{
								&components.ContainerColumn{
									Classes: "p-4",
									Children: []components.PageInterface{
										&components.LabelInline{Title: "Number", Children: []components.PageInterface{
											&components.FieldText{Getter: getters.Key[string]("$in.Number")},
										}},
										&components.LabelInline{Title: "Posted at", Children: []components.PageInterface{
											&components.FieldText{Getter: optionalTimePointerGetter("$in.PostedAt")},
										}},
										&components.LabelInline{Title: "Invoice date", Children: []components.PageInterface{
											&components.FieldText{Getter: invoiceDatetimeStringGetter("$in.Datetime")},
										}},
										&components.LabelInline{Title: "Customer", Children: []components.PageInterface{
											&components.FieldText{Getter: getters.Key[string]("$in.Customer.Name")},
										}},
										&components.LabelInline{Title: "Payment term", Children: []components.PageInterface{
											&components.FieldText{Getter: invoiceDetailPaymentTermSummaryGetterPosted()},
										}},
										&components.LabelInline{Title: "Journal entry", Children: []components.PageInterface{
											&components.FieldLink{
												Href:  journalEntryLinkGetter("$in.JournalEntryID"),
												Label: getters.Format("Entry #%d", getters.Any(getters.Key[uint]("$in.JournalEntryID"))),
											},
										}},
										&components.LabelInline{Title: "Taxes", Children: []components.PageInterface{
											&components.FieldText{Getter: invoiceDetailTaxesNamesGetterPosted()},
										}},
										&components.LabelInline{Title: "Lines", Children: []components.PageInterface{
											&FieldInvoiceLines{Getter: postedInvoiceLinesDisplayGetter()},
										}},
									},
								},
							},
						},
					},
				},
			},
		}},
		{Key: "finance_invoices.PostedInvoiceDetailMenu", Value: &components.SidebarMenu{
			Title: getters.Format("Posted %s", getters.Any(getters.Key[string]("posted_invoice.Number"))),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("Posted invoices"),
				Url:   lamu.RoutePath("finance_invoices.PostedInvoiceListRoute", nil),
			},
			Children: []components.PageInterface{
				&components.SidebarMenuItem{
					Title: getters.Static("Detail"),
					Url: lamu.RoutePath("finance_invoices.PostedInvoiceDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("posted_invoice.ID")),
					}),
				},
				&components.SidebarMenuItem{
					Page:  components.Page{Roles: []string{"superuser"}},
					Title: getters.Static("Cancel invoice"),
					Url: lamu.RoutePath("finance_invoices.PostedInvoiceCancelRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("posted_invoice.ID")),
					}),
				},
			},
		}},
		{Key: "finance_invoices.PostedInvoiceCancelForm", Value: &components.ShellScaffold{
			Page:    components.Page{Roles: []string{"superuser"}},
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_invoices.PostedInvoiceDetailMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name: getters.Static("finance_invoices.PostedInvoiceCancelInner"),
					ActionURL: lamu.RoutePath("finance_invoices.PostedInvoiceCancelRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("posted_invoice.ID")),
					}),
					Children: []components.PageInterface{
						&components.FormComponent[struct{}]{
							Title:    "Cancel invoice",
							Subtitle: "Creates a credit note reversing the journal entry.",
							ChildrenInput: []components.PageInterface{
								&components.InputText{Name: "Reason", Label: "Reason"},
							},
							ChildrenAction: []components.PageInterface{
								&components.ButtonSubmit{Label: "Cancel invoice"},
							},
						},
					},
				},
			},
		}},
	}
}

func invoiceDetailPaymentTermSummaryGetterPosted() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		typ, err := getters.Key[string]("$in.PaymentTerm.Type")(ctx)
		if err != nil {
			return "", err
		}
		bid, err := getters.Key[uint]("$in.PaymentTerm.BackingID")(ctx)
		if err != nil {
			return "", err
		}
		id, err := getters.Key[uint]("$in.PaymentTerm.ID")(ctx)
		if err != nil {
			return "", err
		}
		inst, err := ResolvePaymentTermInstance(ctx, typ, bid)
		if err != nil {
			return fmt.Sprintf("#%d", id), nil
		}
		return fmt.Sprintf("#%d — %s", id, inst.Summary()), nil
	}
}

func invoiceDetailTaxesNamesGetterPosted() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		m, ok := ctx.Value("$in").(map[string]any)
		if !ok || m == nil {
			return "—", nil
		}
		raw, ok := m["Taxes"]
		if !ok || raw == nil {
			return "—", nil
		}
		taxes, ok := raw.([]finance_taxes.Tax)
		if !ok || len(taxes) == 0 {
			return "—", nil
		}
		names := make([]string, 0, len(taxes))
		for _, t := range taxes {
			names = append(names, t.Name)
		}
		return strings.Join(names, ", "), nil
	}
}

func pageEntriesCancelledInvoicePages() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "finance_invoices.CancelledInvoiceTable", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_accounts.MainMenu"}},
			Children: []components.PageInterface{
				&components.ContainerColumn{
					Classes: "gap-2 mb-2 w-full",
					Children: []components.PageInterface{
						&components.Environment[uint]{
							Label:   "Fiscal year",
							Key:     getters.Static(FinanceInvoicesEnvironmentFiscalYearKey),
							Options: fiscalYearsEnvironmentOptionsGetter,
							Default: invoiceFiscalYearEnvironmentDefault,
							Classes: "w-full",
						},
					},
				},
				&components.DataTable[CancelledInvoice]{
					UID:     "finance-cancelled-invoice-table",
					Classes: "w-full",
					Data:    getters.Key[components.ObjectList[CancelledInvoice]]("cancelled_invoices"),
					Actions: []components.PageInterface{
						&components.TableButtonFilter{Child: lamu.DynamicPage{Name: "finance_invoices.InvoiceFilter"}},
					},
					RowAttr: getters.RowAttrNavigate(lamu.RoutePath("finance_invoices.CancelledInvoiceDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("$row.ID")),
					})),
					Columns: []components.TableColumn{
						{Label: "Number", Name: "Number", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Number")},
						}},
						{Label: "Cancelled at", Name: "CancelledAt", Children: []components.PageInterface{
							&components.FieldText{Getter: optionalTimePointerGetter("$row.CancelledAt")},
						}},
						{Label: "Customer", Name: "Customer", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Customer.Name")},
						}},
					},
				},
			},
		}},
		{Key: "finance_invoices.CancelledInvoiceDetail", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_invoices.CancelledInvoiceDetailMenu"}},
			Children: []components.PageInterface{
				&components.ContainerError{
					Error: getters.Key[error]("$error._global"),
					Children: []components.PageInterface{
						&components.Detail[CancelledInvoice]{
							Getter: getters.Key[CancelledInvoice]("cancelled_invoice"),
							Children: []components.PageInterface{
								&components.ContainerColumn{
									Classes: "p-4",
									Children: []components.PageInterface{
										&components.LabelInline{Title: "Number", Children: []components.PageInterface{
											&components.FieldText{Getter: getters.Key[string]("$in.Number")},
										}},
										&components.LabelInline{Title: "Cancelled at", Children: []components.PageInterface{
											&components.FieldText{Getter: optionalTimePointerGetter("$in.CancelledAt")},
										}},
										&components.LabelInline{Title: "Invoice date", Children: []components.PageInterface{
											&components.FieldText{Getter: invoiceDatetimeStringGetter("$in.Datetime")},
										}},
										&components.LabelInline{Title: "Customer", Children: []components.PageInterface{
											&components.FieldText{Getter: getters.Key[string]("$in.Customer.Name")},
										}},
										&components.LabelInline{Title: "Credit note", Children: []components.PageInterface{
											&components.FieldText{Getter: getters.Format("#%d", getters.Any(getters.Key[uint]("$in.CreditNoteID")))},
										}},
										&components.LabelInline{Title: "Lines", Children: []components.PageInterface{
											&FieldInvoiceLines{Getter: cancelledInvoiceLinesDisplayGetter()},
										}},
									},
								},
							},
						},
					},
				},
			},
		}},
		{Key: "finance_invoices.CancelledInvoiceDetailMenu", Value: &components.SidebarMenu{
			Title: getters.Format("Cancelled %s", getters.Any(getters.Key[string]("cancelled_invoice.Number"))),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("Cancelled invoices"),
				Url:   lamu.RoutePath("finance_invoices.CancelledInvoiceListRoute", nil),
			},
			Children: []components.PageInterface{
				&components.SidebarMenuItem{
					Title: getters.Static("Detail"),
					Url: lamu.RoutePath("finance_invoices.CancelledInvoiceDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("cancelled_invoice.ID")),
					}),
				},
				&components.SidebarMenuItem{
					Page:  components.Page{Roles: []string{"superuser"}},
					Title: getters.Static("New draft from this"),
					Url: lamu.RoutePath("finance_invoices.CancelledInvoiceNewDraftRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("cancelled_invoice.ID")),
					}),
				},
			},
		}},
		{Key: "finance_invoices.CancelledInvoiceNewDraftForm", Value: &components.ShellScaffold{
			Page:    components.Page{Roles: []string{"superuser"}},
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_invoices.CancelledInvoiceDetailMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name: getters.Static("finance_invoices.CancelledNewDraftInner"),
					ActionURL: lamu.RoutePath("finance_invoices.CancelledInvoiceNewDraftRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("cancelled_invoice.ID")),
					}),
					Children: []components.PageInterface{
						&components.FormComponent[struct{}]{
							Title:         "New draft",
							Subtitle:      "Creates a new editable draft copied from this cancellation.",
							ChildrenInput: []components.PageInterface{},
							ChildrenAction: []components.PageInterface{
								&components.ButtonSubmit{Label: "Create draft"},
							},
						},
					},
				},
			},
		}},
	}
}
