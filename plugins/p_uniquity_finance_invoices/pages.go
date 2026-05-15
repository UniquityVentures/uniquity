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
	finance_customer "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_customer"
	finance_products "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_products"
	finance_taxes "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_taxes"
)

func pluginPages() lamu.PluginFeatures[components.PageInterface] {
	e := pageEntriesInvoiceMenus()
	e = append(e, pageEntriesInvoicePages()...)
	e = append(e, pageEntriesPaymentTermPages()...)
	e = append(e, pageEntriesPaymentTermFkSelectPages()...)
	return lamu.PluginFeatures[components.PageInterface]{Entries: e}
}

func pageEntriesInvoiceMenus() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "finance_invoices.MainMenu", Value: &components.SidebarMenu{
			Title: getters.Static("Finance invoices"),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("Back to Home"),
				Url:   lamu.RoutePath("dashboard.AppsPage", nil),
			},
			Children: []components.PageInterface{
				&components.SidebarMenuItem{
					Page:  components.Page{Roles: []string{"superuser"}},
					Title: getters.Static("Invoices"),
					Url:   lamu.RoutePath("finance_invoices.DefaultRoute", nil),
					Icon:  "document-text",
				},
				&components.SidebarMenuItem{
					Page:  components.Page{Roles: []string{"superuser"}},
					Title: getters.Static("Payment terms"),
					Url:   lamu.RoutePath("finance_invoices.PaymentTermListRoute", nil),
					Icon:  "calendar-days",
				},
			},
		}},
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

func invoiceStatusChoices() getters.Getter[[]registry.Pair[string, string]] {
	return getters.Static([]registry.Pair[string, string]{
		{Key: string(InvoiceStatusDraft), Value: "Draft"},
		{Key: string(InvoiceStatusPosted), Value: "Posted"},
		{Key: string(InvoiceStatusCancelled), Value: "Cancelled"},
	})
}

func invoiceCreateStatusPairGetter() getters.Getter[registry.Pair[string, string]] {
	return func(ctx context.Context) (registry.Pair[string, string], error) {
		s, err := getters.Key[string]("$in.Status")(ctx)
		if err != nil || s == "" {
			return registry.Pair[string, string]{Key: string(InvoiceStatusDraft), Value: "Draft"}, nil
		}
		if p, ok := registry.PairFromPairs(s, *invoiceStatusChoicesStatic()); ok {
			return p, nil
		}
		return registry.Pair[string, string]{Key: s, Value: s}, nil
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

func invoiceStatusChoicesStatic() *[]registry.Pair[string, string] {
	s := []registry.Pair[string, string]{
		{Key: string(InvoiceStatusDraft), Value: "Draft"},
		{Key: string(InvoiceStatusPosted), Value: "Posted"},
		{Key: string(InvoiceStatusCancelled), Value: "Cancelled"},
	}
	return &s
}

func invoiceStatusLabelRowGetter() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		s, err := getters.Key[InvoiceStatus]("$row.Status")(ctx)
		if err != nil {
			return "", err
		}
		if p, ok := registry.PairFromPairs(string(s), *invoiceStatusChoicesStatic()); ok {
			return p.Value, nil
		}
		return string(s), nil
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

func invoiceStatusLabelDetailGetter() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		s, err := getters.Key[InvoiceStatus]("$in.Status")(ctx)
		if err != nil {
			return "", err
		}
		if p, ok := registry.PairFromPairs(string(s), *invoiceStatusChoicesStatic()); ok {
			return p.Value, nil
		}
		return string(s), nil
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
		return `[{"product_id":0,"quantity":"1"}]`, nil
	}
}

func invoiceDetailLinesSummaryGetter() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		m, ok := ctx.Value("$in").(map[string]any)
		if !ok || m == nil {
			return "", nil
		}
		raw, ok := m["Lines"]
		if !ok || raw == nil {
			return "—", nil
		}
		lines, ok := raw.([]InvoiceLine)
		if !ok || len(lines) == 0 {
			return "—", nil
		}
		var b strings.Builder
		for i, ln := range lines {
			if i > 0 {
				b.WriteString("; ")
			}
			name := ln.Product.Name
			if name == "" {
				name = fmt.Sprintf("#%d", ln.ProductID)
			}
			fmt.Fprintf(&b, "%s × %s", name, ln.Quantity.String())
		}
		return b.String(), nil
	}
}

func invoiceCreateFormInputs() []components.PageInterface {
	return []components.PageInterface{
		&components.ContainerError{
			Error: getters.Key[error]("$error.Number"),
			Children: []components.PageInterface{
				&components.InputText{Name: "Number", Label: "Invoice number", Required: true, Getter: getters.Key[string]("$in.Number")},
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
			Error: getters.Key[error]("$error.Status"),
			Children: []components.PageInterface{
				&components.InputSelect[string]{
					Label:    "Status",
					Name:     "Status",
					Required: true,
					Choices:  invoiceStatusChoices(),
					Getter:   invoiceCreateStatusPairGetter(),
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
					Page:     components.Page{Key: "finance_invoices.InvoiceCreateForm.Lines"},
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

func pageEntriesInvoicePages() []registry.Pair[string, components.PageInterface] {
	createName := getters.Static("finance_invoices.InvoiceCreateForm")

	return []registry.Pair[string, components.PageInterface]{
		{Key: "finance_invoices.InvoiceTable", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_invoices.MainMenu"}},
			Children: []components.PageInterface{
				&components.DataTable[Invoice]{
					UID:     "finance-invoice-table",
					Classes: "w-full",
					Data:    getters.Key[components.ObjectList[Invoice]]("invoices"),
					Actions: []components.PageInterface{
						&components.TableButtonCreate{
							Link: lamu.RoutePath("finance_invoices.InvoiceCreateRoute", nil),
							Page: components.Page{Roles: []string{"superuser"}},
						},
					},
					RowAttr: getters.RowAttrNavigate(lamu.RoutePath("finance_invoices.InvoiceDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("$row.ID")),
					})),
					Columns: []components.TableColumn{
						{Label: "Number", Name: "Number", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Number")},
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
						{Label: "Status", Name: "Status", Children: []components.PageInterface{
							&components.FieldText{Getter: invoiceStatusLabelRowGetter()},
						}},
					},
				},
			},
		}},
		{Key: "finance_invoices.InvoiceDetail", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_invoices.InvoiceDetailMenu"}},
			Children: []components.PageInterface{
				&components.Detail[Invoice]{
					Getter: getters.Key[Invoice]("invoice"),
					Children: []components.PageInterface{
						&components.ContainerColumn{
							Classes: "p-4",
							Children: []components.PageInterface{
								&components.LabelInline{Title: "Number", Children: []components.PageInterface{
									&components.FieldText{Getter: getters.Key[string]("$in.Number")},
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
								&components.LabelInline{Title: "Status", Children: []components.PageInterface{
									&components.FieldText{Getter: invoiceStatusLabelDetailGetter()},
								}},
								&components.LabelInline{Title: "Taxes", Children: []components.PageInterface{
									&components.FieldText{Getter: invoiceDetailTaxesNamesGetter()},
								}},
								&components.LabelInline{Title: "Lines", Children: []components.PageInterface{
									&components.FieldText{Getter: invoiceDetailLinesSummaryGetter()},
								}},
							},
						},
					},
				},
			},
		}},
		{Key: "finance_invoices.InvoiceDetailMenu", Value: &components.SidebarMenu{
			Title: getters.Format("Invoice %s", getters.Any(getters.Key[string]("invoice.Number"))),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("All invoices"),
				Url:   lamu.RoutePath("finance_invoices.DefaultRoute", nil),
			},
			Children: []components.PageInterface{
				&components.SidebarMenuItem{
					Title: getters.Static("Detail"),
					Url: lamu.RoutePath("finance_invoices.InvoiceDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("invoice.ID")),
					}),
				},
			},
		}},
		{Key: "finance_invoices.InvoiceCreateForm", Value: &components.ShellScaffold{
			Page:    components.Page{Roles: []string{"superuser"}},
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_invoices.MainMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name:      createName,
					ActionURL: lamu.RoutePath("finance_invoices.InvoiceCreateRoute", nil),
					Children: []components.PageInterface{
						&components.FormComponent[Invoice]{
							Attr:          getters.FormBubbling(createName),
							Title:         "Create invoice",
							Subtitle:      "Customer, dates, payment term, lines, taxes, and status",
							ChildrenInput: invoiceCreateFormInputs(),
							ChildrenAction: []components.PageInterface{
								&components.ButtonSubmit{Label: "Save"},
							},
						},
					},
				},
			},
		}},
	}
}

// invoiceListPaymentTermSummaryGetter resolves a short summary for the invoice row's payment term (requires preload).
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
