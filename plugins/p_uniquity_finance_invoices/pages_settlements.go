package p_uniquity_finance_invoices

import (
	"context"
	"fmt"
	"time"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
	finance_taxes "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_taxes"
)

var settlementPostedInvoiceDetailPreload = []string{
	"Payment",
	"PostedInvoice",
	"PostedInvoice.Customer",
	"PostedInvoice.PaymentTerm",
	"PostedInvoice.Taxes",
	"PostedInvoice.Lines",
	"PostedInvoice.Lines.Product",
	"PostedInvoice.Lines.Taxes",
	"PostedInvoice.JournalEntry",
	"PriorPartial",
}

func settlementPriorPartialSummaryGetter(inOrRowPrefix string) getters.Getter[string] {
	key := inOrRowPrefix + ".PriorPartial"
	return func(ctx context.Context) (string, error) {
		pp, err := getters.Key[*PartiallyPaidInvoice](key)(ctx)
		if err != nil || pp == nil || pp.ID == 0 {
			return "—", nil
		}
		return fmt.Sprintf("#%d", pp.ID), nil
	}
}

func settlementPaymentRowSummaryGetter(rowPrefix string) getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		pid, err := getters.Key[uint](rowPrefix + ".Payment.ID")(ctx)
		if err != nil {
			return "", err
		}
		amt, err := paymentDecimalStringGetter(rowPrefix + ".Payment.Amount")(ctx)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("#%d · %s", pid, amt), nil
	}
}

func settlementPaymentDetailLabelGetter() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		pid, err := getters.Key[uint]("$in.Payment.ID")(ctx)
		if err != nil {
			return "", err
		}
		amt, err := paymentDecimalStringGetter("$in.Payment.Amount")(ctx)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("#%d · %s", pid, amt), nil
	}
}

func settlementPostedCustomerLinkGetter() getters.Getter[string] {
	return lamu.RoutePath("finance_customers.CustomerDetailRoute", map[string]getters.Getter[any]{
		"id": getters.Any(getters.Key[uint]("$in.PostedInvoice.CustomerID")),
	})
}

func settlementPostedInvoicePaymentTermSummaryGetter() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		typ, err := getters.Key[string]("$in.PostedInvoice.PaymentTerm.Type")(ctx)
		if err != nil {
			return "", err
		}
		bid, err := getters.Key[uint]("$in.PostedInvoice.PaymentTerm.BackingID")(ctx)
		if err != nil {
			return "", err
		}
		id, err := getters.Key[uint]("$in.PostedInvoice.PaymentTerm.ID")(ctx)
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

func settlementPostedInvoiceLinesDisplayGetter() getters.Getter[[]InvoiceLineDisplay] {
	return func(ctx context.Context) ([]InvoiceLineDisplay, error) {
		lines, err := getters.Key[[]PostedInvoiceLine]("$in.PostedInvoice.Lines")(ctx)
		if err != nil || len(lines) == 0 {
			return nil, err
		}
		out := make([]InvoiceLineDisplay, 0, len(lines))
		for _, ln := range lines {
			name := ln.Product.Name
			if name == "" {
				name = fmt.Sprintf("#%d", ln.ProductID)
			}
			u, lev, wh, net := invoiceLineAmountBreakdown(ln.Quantity, ln.Rate, ln.Taxes)
			out = append(out, InvoiceLineDisplay{
				Product:           name,
				Quantity:          ln.Quantity.String(),
				Rate:              ln.Rate.String(),
				LineTaxes:         invoiceLineTaxesLabel(ln.Taxes),
				UntaxedAmount:     decimalSixDisplay(u),
				LeviedTaxAmount:   decimalSixDisplay(lev),
				WithholdingAmount: decimalSixDisplayWithholding(wh),
				LineTotal:         decimalSixDisplay(net),
			})
		}
		return out, nil
	}
}

func settlementPostedInvoiceLinesSummaryGetter() getters.Getter[InvoiceLinesSummary] {
	return func(ctx context.Context) (InvoiceLinesSummary, error) {
		lines, _ := getters.Key[[]PostedInvoiceLine]("$in.PostedInvoice.Lines")(ctx)
		taxes, _ := getters.Key[[]finance_taxes.Tax]("$in.PostedInvoice.Taxes")(ctx)
		return invoiceLinesSummaryFromPostedLines(lines, taxes), nil
	}
}

func settlementInvoiceActionRow(showPayAction bool, pdfRoute, recordIDKey string) components.PageInterface {
	children := []components.PageInterface{invoicePdfDownloadButton(pdfRoute, recordIDKey)}
	if showPayAction {
		children = append([]components.PageInterface{
			&components.ButtonLink{
				Page:    components.Page{Roles: []string{"superuser"}},
				Label:   getters.Static("Pay"),
				Link:    paymentCreateURLForPostedInvoiceID("$in.PostedInvoiceID"),
				Classes: "btn-primary",
			},
		}, children...)
	}
	return &components.ContainerRow{
		Classes:  "mb-4 flex flex-wrap gap-2 items-center",
		Children: children,
	}
}

func settlementPostedInvoiceDetailFields() []components.PageInterface {
	return []components.PageInterface{
		&components.LabelInline{Title: "Number", Children: []components.PageInterface{
			&components.FieldText{Getter: getters.Key[string]("$in.PostedInvoice.Number")},
		}},
		&components.LabelInline{Title: "Posted date", Children: []components.PageInterface{
			&optionalFieldDate{Getter: invoiceOptionalDateGetter("$in.PostedInvoice.PostedAt")},
		}},
		&components.LabelInline{Title: "Invoice date", Children: []components.PageInterface{
			&components.FieldDate{Getter: getters.Key[time.Time]("$in.PostedInvoice.Datetime")},
		}},
		&components.LabelInline{Title: "Customer", Children: []components.PageInterface{
			&components.FieldLink{
				Href:    settlementPostedCustomerLinkGetter(),
				Label:   getters.Key[string]("$in.PostedInvoice.Customer.Name"),
				Classes: "link link-hover",
			},
		}},
		&components.LabelInline{Title: "Payment term", Children: []components.PageInterface{
			&components.FieldText{Getter: settlementPostedInvoicePaymentTermSummaryGetter()},
		}},
		&components.LabelInline{Title: "Journal entry", Children: []components.PageInterface{
			&components.FieldLink{
				Href:  journalEntryLinkGetter("$in.PostedInvoice.JournalEntryID"),
				Label: getters.Format("Entry #%d", getters.Any(getters.Key[uint]("$in.PostedInvoice.JournalEntryID"))),
			},
		}},
		&components.LabelInline{Title: "Taxes", Children: []components.PageInterface{
			&components.FieldManyToMany[finance_taxes.Tax]{
				Getter:    getters.Key[[]finance_taxes.Tax]("$in.PostedInvoice.Taxes"),
				Display:   getters.Key[string]("$in.Name"),
				Link:      invoiceTaxDetailLinkHrefGetter(),
				Classes:   "w-full max-w-xl",
				EmptyText: "—",
			},
		}},
		&components.LabelNewline{Title: "Lines", Children: []components.PageInterface{
			&FieldInvoiceLines{Getter: settlementPostedInvoiceLinesDisplayGetter()},
			&FieldInvoiceLinesSummary{Getter: settlementPostedInvoiceLinesSummaryGetter()},
		}},
	}
}

func settlementPaymentDetailFields() []components.PageInterface {
	return []components.PageInterface{
		&components.LabelInline{Title: "Payment", Children: []components.PageInterface{
			&components.FieldLink{
				Href: lamu.RoutePath("finance_invoices.PaymentDetailRoute", map[string]getters.Getter[any]{
					"id": getters.Any(getters.Key[uint]("$in.PaymentID")),
				}),
				Label:   settlementPaymentDetailLabelGetter(),
				Classes: "link link-hover",
			},
		}},
		&components.LabelInline{Title: "Payment date", Children: []components.PageInterface{
			&components.FieldDate{Getter: getters.Key[time.Time]("$in.Payment.Datetime")},
		}},
		&components.LabelInline{Title: "Prior partial record", Children: []components.PageInterface{
			&components.FieldText{Getter: settlementPriorPartialSummaryGetter("$in")},
		}},
	}
}

func settlementInvoiceDetailColumnChildren(showPayAction bool, pdfRoute, recordIDKey string) []components.PageInterface {
	children := make([]components.PageInterface, 0, 16)
	children = append(children, settlementInvoiceActionRow(showPayAction, pdfRoute, recordIDKey))
	children = append(children, settlementPostedInvoiceDetailFields()...)
	children = append(children, settlementPaymentDetailFields()...)
	return children
}

func paidInvoiceHubTable() *components.DataTable[PaidInvoice] {
	return &components.DataTable[PaidInvoice]{
		UID:     "finance-paid-invoice-table",
		Classes: "w-full",
		Data:    getters.Key[components.ObjectList[PaidInvoice]]("paid_invoices"),
		RowAttr: getters.RowAttrNavigate(lamu.RoutePath("finance_invoices.PaidInvoiceDetailRoute", map[string]getters.Getter[any]{
			"id": getters.Any(getters.Key[uint]("$row.ID")),
		})),
		Columns: []components.TableColumn{
			{Label: "ID", Name: "ID", Children: []components.PageInterface{
				&components.FieldText{Getter: getters.Format("%d", getters.Any(getters.Key[uint]("$row.ID")))},
			}},
			{Label: "Invoice", Name: "Invoice", Children: []components.PageInterface{
				&components.FieldText{Getter: getters.Format("%s (#%d)", getters.Any(getters.Key[string]("$row.PostedInvoice.Number")), getters.Any(getters.Key[uint]("$row.PostedInvoice.ID")))},
			}},
			{Label: "Payment", Name: "Payment", Children: []components.PageInterface{
				&components.FieldText{Getter: settlementPaymentRowSummaryGetter("$row")},
			}},
			{Label: "Payment date", Name: "PaymentDatetime", Children: []components.PageInterface{
				&components.FieldDate{Getter: getters.Key[time.Time]("$row.Payment.Datetime")},
			}},
			{Label: "Prior partial", Name: "PriorPartial", Children: []components.PageInterface{
				&components.FieldText{Getter: settlementPriorPartialSummaryGetter("$row")},
			}},
		},
	}
}

func partiallyPaidInvoiceHubTable() *components.DataTable[PartiallyPaidInvoice] {
	return &components.DataTable[PartiallyPaidInvoice]{
		UID:     "finance-partially-paid-invoice-table",
		Classes: "w-full",
		Data:    getters.Key[components.ObjectList[PartiallyPaidInvoice]]("partially_paid_invoices"),
		RowAttr: getters.RowAttrNavigate(lamu.RoutePath("finance_invoices.PartiallyPaidInvoiceDetailRoute", map[string]getters.Getter[any]{
			"id": getters.Any(getters.Key[uint]("$row.ID")),
		})),
		Columns: []components.TableColumn{
			{Label: "ID", Name: "ID", Children: []components.PageInterface{
				&components.FieldText{Getter: getters.Format("%d", getters.Any(getters.Key[uint]("$row.ID")))},
			}},
			{Label: "Invoice", Name: "Invoice", Children: []components.PageInterface{
				&components.FieldText{Getter: getters.Format("%s (#%d)", getters.Any(getters.Key[string]("$row.PostedInvoice.Number")), getters.Any(getters.Key[uint]("$row.PostedInvoice.ID")))},
			}},
			{Label: "Payment", Name: "Payment", Children: []components.PageInterface{
				&components.FieldText{Getter: settlementPaymentRowSummaryGetter("$row")},
			}},
			{Label: "Payment date", Name: "PaymentDatetime", Children: []components.PageInterface{
				&components.FieldDate{Getter: getters.Key[time.Time]("$row.Payment.Datetime")},
			}},
			{Label: "Prior partial", Name: "PriorPartial", Children: []components.PageInterface{
				&components.FieldText{Getter: settlementPriorPartialSummaryGetter("$row")},
			}},
		},
	}
}

func pageEntriesSettlementPages() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "finance_invoices.PaidInvoiceDetail", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_invoices.PaidInvoiceDetailMenu"}},
			Children: []components.PageInterface{
				&components.ContainerError{
					Error: getters.Key[error]("$error._global"),
					Children: []components.PageInterface{
						&components.Detail[PaidInvoice]{
							Getter: getters.Key[PaidInvoice]("paid_invoice"),
							Children: []components.PageInterface{
								&components.ContainerColumn{
									Classes: "p-4",
									Children: settlementInvoiceDetailColumnChildren(false, "finance_invoices.PaidInvoicePdfRoute", "paid_invoice.ID"),
								},
							},
						},
					},
				},
			},
		}},
		{Key: "finance_invoices.PaidInvoiceDetailMenu", Value: &components.SidebarMenu{
			Title: getters.Format("Paid %s", getters.Any(getters.Key[string]("paid_invoice.PostedInvoice.Number"))),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("Invoices"),
				Url:   invoiceHubURLWithTabGetter("paid"),
			},
			Children: []components.PageInterface{
				&components.SidebarMenuItem{
					Title: getters.Static("Detail"),
					Url: lamu.RoutePath("finance_invoices.PaidInvoiceDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("paid_invoice.ID")),
					}),
				},
			},
		}},

		{Key: "finance_invoices.PartiallyPaidInvoiceDetail", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_invoices.PartiallyPaidInvoiceDetailMenu"}},
			Children: []components.PageInterface{
				&components.ContainerError{
					Error: getters.Key[error]("$error._global"),
					Children: []components.PageInterface{
						&components.Detail[PartiallyPaidInvoice]{
							Getter: getters.Key[PartiallyPaidInvoice]("partially_paid_invoice"),
							Children: []components.PageInterface{
								&components.ContainerColumn{
									Classes:  "p-4",
									Children: settlementInvoiceDetailColumnChildren(true, "finance_invoices.PartiallyPaidInvoicePdfRoute", "partially_paid_invoice.ID"),
								},
							},
						},
					},
				},
			},
		}},
		{Key: "finance_invoices.PartiallyPaidInvoiceDetailMenu", Value: &components.SidebarMenu{
			Title: getters.Format("Partial %s", getters.Any(getters.Key[string]("partially_paid_invoice.PostedInvoice.Number"))),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("Invoices"),
				Url:   invoiceHubURLWithTabGetter("partial"),
			},
			Children: []components.PageInterface{
				&components.SidebarMenuItem{
					Title: getters.Static("Detail"),
					Url: lamu.RoutePath("finance_invoices.PartiallyPaidInvoiceDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("partially_paid_invoice.ID")),
					}),
				},
			},
		}},
	}
}
