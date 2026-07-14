package p_uniquity_finance_invoices

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	finance_taxes "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_taxes"
	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/components"
	"github.com/lariv-in/lago/fields"
	"github.com/lariv-in/lago/getters"
	"github.com/lariv-in/lago/registry"
	"maragu.dev/gomponents"
)

func paymentDecimalStringGetter(ctxKey string) getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		pd, err := getters.Key[fields.DecimalSix](ctxKey)(ctx)
		if err != nil {
			return "", err
		}
		return pd.String(), nil
	}
}

func paymentDecimalGetter(ctxKey string) getters.Getter[fields.DecimalSix] {
	return func(ctx context.Context) (fields.DecimalSix, error) {
		return getters.Key[fields.DecimalSix](ctxKey)(ctx)
	}
}

func paymentCreateFormAttr(formName getters.Getter[string]) getters.Getter[gomponents.Node] {
	return func(ctx context.Context) (gomponents.Node, error) {
		bub, err := getters.FormBubbling(formName)(ctx)
		if err != nil {
			return nil, err
		}
		reload, err := paymentCreatePostedInvoiceReloadAttr()(ctx)
		if err != nil {
			return nil, err
		}
		return gomponents.Group{bub, reload}, nil
	}
}

// paymentCreatePostedInvoiceReloadAttr reloads the create form via boosted GET when a posted
// invoice is picked so the query-defaults layer can pre-fill Amount with the open balance.
func paymentCreatePostedInvoiceReloadAttr() getters.Getter[gomponents.Node] {
	route := lago.RoutePath("finance_invoices.PaymentCreateRoute", nil)
	return func(ctx context.Context) (gomponents.Node, error) {
		url, err := route(ctx)
		if err != nil {
			return nil, err
		}
		urlLit, err := json.Marshal(url)
		if err != nil {
			return nil, err
		}
		script := fmt.Sprintf(
			`(function(evt){var d=evt.detail;if(!d||d.name!=='PostedInvoiceID'||!d.value)return;var f=evt.target.closest('form');if(!f)return;var p=new URLSearchParams(htmx.values(f));p.set('PostedInvoiceID',d.value);var o={swap:'outerHTML',headers:{'HX-Boosted':'true'}};o.target='body';htmx.ajax('GET',%s+'?'+p.toString(),o)})($event)`,
			urlLit,
		)
		return gomponents.Attr("@fk-select.window", script), nil
	}
}

func paymentTaxNamesGetter(ctxKey string) getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		taxes, err := getters.Key[[]finance_taxes.Tax](ctxKey)(ctx)
		if err != nil || len(taxes) == 0 {
			return "—", err
		}
		names := make([]string, 0, len(taxes))
		for _, t := range taxes {
			n := t.Name
			if n == "" {
				n = fmt.Sprintf("#%d", t.ID)
			}
			names = append(names, n)
		}
		return strings.Join(names, ", "), nil
	}
}

func paymentComputedDecimalGetter(amountKey, taxesKey string, fn func(fields.DecimalSix, []finance_taxes.Tax) fields.DecimalSix) getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		amt, err := getters.Key[fields.DecimalSix](amountKey)(ctx)
		if err != nil {
			return "", err
		}
		taxes, _ := getters.Key[[]finance_taxes.Tax](taxesKey)(ctx)
		return fn(amt, taxes).String(), nil
	}
}

func paymentCreateDatetimeGetter() getters.Getter[time.Time] {
	return func(ctx context.Context) (time.Time, error) {
		t, err := getters.Key[time.Time]("$in.Datetime")(ctx)
		if err != nil || t.IsZero() {
			return time.Now(), nil
		}
		return t, nil
	}
}

func paymentCreateURLForPostedInvoiceID(postedInvoiceIDKey string) getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		base, err := lago.RoutePath("finance_invoices.PaymentCreateRoute", nil)(ctx)
		if err != nil {
			return "", err
		}
		postedID, err := getters.Key[uint](postedInvoiceIDKey)(ctx)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s?PostedInvoiceID=%d", base, postedID), nil
	}
}

func paymentSettlementKindAndHref(ctx context.Context) (label string, href string, err error) {
	pid, err := getters.Key[uint]("payment.ID")(ctx)
	if err != nil {
		return "", "", err
	}
	db, err := getters.DBFromContext(ctx)
	if err != nil {
		return "", "", err
	}
	var paid PaidInvoice
	if err := db.Where("payment_id = ? AND deleted_at IS NULL", pid).Take(&paid).Error; err == nil {
		href, err = lago.RoutePath("finance_invoices.PaidInvoiceDetailRoute", map[string]getters.Getter[any]{
			"id": getters.Any(getters.Static(paid.ID)),
		})(ctx)
		if err != nil {
			return "", "", err
		}
		return fmt.Sprintf("Paid in full #%d", paid.ID), href, nil
	}
	var partial PartiallyPaidInvoice
	if err := db.Where("payment_id = ? AND deleted_at IS NULL", pid).Take(&partial).Error; err == nil {
		href, err = lago.RoutePath("finance_invoices.PartiallyPaidInvoiceDetailRoute", map[string]getters.Getter[any]{
			"id": getters.Any(getters.Static(partial.ID)),
		})(ctx)
		if err != nil {
			return "", "", err
		}
		return fmt.Sprintf("Partial payment #%d", partial.ID), href, nil
	}
	return "—", "", nil
}

func paymentSettlementLabelGetter() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		label, _, err := paymentSettlementKindAndHref(ctx)
		return label, err
	}
}

func paymentSettlementHrefGetter() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		_, href, err := paymentSettlementKindAndHref(ctx)
		if err != nil {
			return "", err
		}
		if href != "" {
			return href, nil
		}
		return lago.RoutePath("finance_invoices.PaymentListRoute", nil)(ctx)
	}
}

func paymentCreateFormInputs() []components.PageInterface {
	return []components.PageInterface{
		paymentImmutableWarning{Page: components.Page{Key: "finance_invoices.PaymentCreateWarning"}},
		&components.ContainerError{
			Error: getters.Key[error]("$error.PostedInvoiceID"),
			Children: []components.PageInterface{
				&components.InputForeignKey[PostedInvoice]{
					Label:       "Posted invoice",
					Name:        "PostedInvoiceID",
					Required:    true,
					Url:         lago.RoutePath("finance_invoices.PostedInvoiceFkSelectRoute", nil),
					Display:     getters.Format("%s (#%d)", getters.Any(getters.Key[string]("$in.Number")), getters.Any(getters.Key[uint]("$in.ID"))),
					Placeholder: "Select posted invoice…",
					Getter:      getters.Association[PostedInvoice, uint](getters.Key[uint]("$in.PostedInvoiceID")),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Amount"),
			Children: []components.PageInterface{
				&components.InputPointsDecimal{
					Label:    "Settlement amount",
					Name:     "Amount",
					Required: true,
					Getter:   paymentDecimalGetter("$in.Amount"),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Taxes"),
			Children: []components.PageInterface{
				&components.InputManyToMany[finance_taxes.Tax]{
					Label:       "Withholding taxes",
					Name:        "Taxes",
					Display:     getters.Key[string]("$in.Name"),
					Getter:      getters.Key[[]finance_taxes.Tax]("$in.Taxes"),
					Url:         lago.RoutePath("finance_taxes.TaxMultiSelectRoute", nil),
					Placeholder: "Optional withholding at collection…",
					Classes:     "w-full",
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Datetime"),
			Children: []components.PageInterface{
				&components.InputDatetime{Label: "Payment date & time", Name: "Datetime", Required: true, Getter: paymentCreateDatetimeGetter()},
			},
		},
	}
}

func pageEntriesPaymentPages() []registry.Pair[string, components.PageInterface] {
	createName := getters.Static("finance_invoices.PaymentCreateForm")
	return []registry.Pair[string, components.PageInterface]{
		{Key: "finance_invoices.PaymentTable", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lago.DynamicPage{Name: "finance_accounts.MainMenu"}},
			Children: []components.PageInterface{
				&components.DataTable[Payment]{
					UID:     "finance-payment-table",
					Classes: "w-full",
					Data:    getters.Key[components.ObjectList[Payment]]("payments"),
					Actions: []components.PageInterface{
						&components.TableButtonCreate{
							Link: lago.RoutePath("finance_invoices.PaymentCreateRoute", nil),
							Page: components.Page{Roles: []string{"superuser"}},
						},
					},
					RowAttr: getters.RowAttrNavigate(lago.RoutePath("finance_invoices.PaymentDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("$row.ID")),
					})),
					Columns: []components.TableColumn{
						{Label: "ID", Name: "ID", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Format("%d", getters.Any(getters.Key[uint]("$row.ID")))},
						}},
						{Label: "Date", Name: "Datetime", Children: []components.PageInterface{
							&components.FieldDate{Getter: getters.Key[time.Time]("$row.Datetime")},
						}},
						{Label: "Settlement", Name: "Amount", Children: []components.PageInterface{
							&components.FieldText{Getter: paymentDecimalStringGetter("$row.Amount")},
						}},
						{Label: "Invoice", Name: "Invoice", Children: []components.PageInterface{
							&components.FieldText{Getter: paymentPostedInvoiceNumberRowGetter()},
						}},
						{Label: "Account", Name: "Account", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Format("%s (#%d)", getters.Any(getters.Key[string]("$row.Account.Name")), getters.Any(getters.Key[uint]("$row.Account.ID")))},
						}},
					},
				},
			},
		}},
		{Key: "finance_invoices.PaymentCreateForm", Value: &components.ShellScaffold{
			Page:    components.Page{Roles: []string{"superuser"}},
			Sidebar: []components.PageInterface{lago.DynamicPage{Name: "finance_accounts.MainMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name:      createName,
					ActionURL: lago.RoutePath("finance_invoices.PaymentCreateRoute", nil),
					Children: []components.PageInterface{
						&components.FormComponent[Payment]{
							Attr:          paymentCreateFormAttr(createName),
							Title:         "Record payment",
							Subtitle:      "Applies to a posted invoice and posts a journal entry",
							ChildrenInput: paymentCreateFormInputs(),
							ChildrenAction: []components.PageInterface{
								&components.ButtonSubmit{Label: "Save"},
							},
						},
					},
				},
			},
		}},
		{Key: "finance_invoices.PaymentDetail", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lago.DynamicPage{Name: "finance_invoices.PaymentDetailMenu"}},
			Children: []components.PageInterface{
				&components.ContainerError{
					Error: getters.Key[error]("$error._global"),
					Children: []components.PageInterface{
						&components.Detail[Payment]{
							Getter: getters.Key[Payment]("payment"),
							Children: []components.PageInterface{
								&components.ContainerColumn{
									Classes: "p-4",
									Children: []components.PageInterface{
										&components.LabelInline{Title: "Settlement amount", Children: []components.PageInterface{
											&components.FieldText{Getter: paymentDecimalStringGetter("$in.Amount")},
										}},
										&components.LabelInline{Title: "Withholding", Children: []components.PageInterface{
											&components.FieldText{Getter: paymentComputedDecimalGetter("$in.Amount", "$in.Taxes", paymentWithholdingTotal)},
										}},
										&components.LabelInline{Title: "Bank received", Children: []components.PageInterface{
											&components.FieldText{Getter: paymentComputedDecimalGetter("$in.Amount", "$in.Taxes", paymentBankAmount)},
										}},
										&components.LabelInline{Title: "Withholding taxes", Children: []components.PageInterface{
											&components.FieldText{Getter: paymentTaxNamesGetter("$in.Taxes")},
										}},
										&components.LabelInline{Title: "Date", Children: []components.PageInterface{
											&components.FieldDate{Getter: getters.Key[time.Time]("$in.Datetime")},
										}},
										&components.LabelInline{Title: "Posted invoice", Children: []components.PageInterface{
											&components.FieldLink{
												Href: lago.RoutePath("finance_invoices.PostedInvoiceDetailRoute", map[string]getters.Getter[any]{
													"id": getters.Any(getters.Key[uint]("$in.PostedInvoiceID")),
												}),
												Label:   getters.Format("%s (#%d)", getters.Any(getters.Key[string]("$in.PostedInvoice.Number")), getters.Any(getters.Key[uint]("$in.PostedInvoice.ID"))),
												Classes: "link link-hover",
											},
										}},
										&components.LabelInline{Title: "Account", Children: []components.PageInterface{
											&components.FieldLink{
												Href: lago.RoutePath("finance_accounts.AccountDetailRoute", map[string]getters.Getter[any]{
													"id": getters.Any(getters.Key[uint]("$in.AccountID")),
												}),
												Label:   getters.Format("%s (#%d)", getters.Any(getters.Key[string]("$in.Account.Name")), getters.Any(getters.Key[uint]("$in.Account.ID"))),
												Classes: "link link-hover",
											},
										}},
										&components.LabelInline{Title: "Journal entry", Children: []components.PageInterface{
											&components.FieldLink{
												Href:    journalEntryLinkGetter("$in.JournalEntryID"),
												Label:   getters.Format("Entry #%d", getters.Any(getters.Key[uint]("$in.JournalEntryID"))),
												Classes: "link link-hover",
											},
										}},
										&components.LabelInline{Title: "Settlement", Children: []components.PageInterface{
											&components.FieldLink{
												Href:    paymentSettlementHrefGetter(),
												Label:   paymentSettlementLabelGetter(),
												Classes: "link link-hover",
											},
										}},
									},
								},
							},
						},
					},
				},
			},
		}},
		{Key: "finance_invoices.PaymentDetailMenu", Value: &components.SidebarMenu{
			Title: getters.Format("Payment #%d", getters.Any(getters.Key[uint]("payment.ID"))),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("Payments"),
				Url:   lago.RoutePath("finance_invoices.PaymentListRoute", nil),
			},
			Children: []components.PageInterface{
				&components.SidebarMenuItem{
					Title: getters.Static("Detail"),
					Url: lago.RoutePath("finance_invoices.PaymentDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("payment.ID")),
					}),
				},
			},
		}},
	}
}

func paymentPostedInvoiceNumberRowGetter() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		n, err := getters.Key[string]("$row.PostedInvoice.Number")(ctx)
		if err != nil {
			return "", err
		}
		id, err := getters.Key[uint]("$row.PostedInvoice.ID")(ctx)
		if err != nil {
			return n, nil
		}
		return fmt.Sprintf("%s (#%d)", n, id), nil
	}
}
