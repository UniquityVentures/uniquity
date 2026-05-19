package p_uniquity_finance_invoices

import (
	"context"
	"fmt"
	"time"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

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
									Children: []components.PageInterface{
										&components.LabelInline{Title: "Posted invoice", Children: []components.PageInterface{
											&components.FieldLink{
												Href: lamu.RoutePath("finance_invoices.PostedInvoiceDetailRoute", map[string]getters.Getter[any]{
													"id": getters.Any(getters.Key[uint]("$in.PostedInvoiceID")),
												}),
												Label: getters.Format("%s (#%d)", getters.Any(getters.Key[string]("$in.PostedInvoice.Number")), getters.Any(getters.Key[uint]("$in.PostedInvoice.ID"))),
												Classes: "link link-hover",
											},
										}},
										&components.LabelInline{Title: "Payment", Children: []components.PageInterface{
											&components.FieldLink{
												Href: lamu.RoutePath("finance_invoices.PaymentDetailRoute", map[string]getters.Getter[any]{
													"id": getters.Any(getters.Key[uint]("$in.PaymentID")),
												}),
												Label: settlementPaymentDetailLabelGetter(),
												Classes: "link link-hover",
											},
										}},
										&components.LabelInline{Title: "Payment date", Children: []components.PageInterface{
											&components.FieldDate{Getter: getters.Key[time.Time]("$in.Payment.Datetime")},
										}},
										&components.LabelInline{Title: "Prior partial record", Children: []components.PageInterface{
											&components.FieldText{Getter: settlementPriorPartialSummaryGetter("$in")},
										}},
									},
								},
							},
						},
					},
				},
			},
		}},
		{Key: "finance_invoices.PaidInvoiceDetailMenu", Value: &components.SidebarMenu{
			Title: getters.Static("Paid invoice"),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("Paid invoices"),
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
									Classes: "p-4",
									Children: []components.PageInterface{
										&components.LabelInline{Title: "Posted invoice", Children: []components.PageInterface{
											&components.FieldLink{
												Href: lamu.RoutePath("finance_invoices.PostedInvoiceDetailRoute", map[string]getters.Getter[any]{
													"id": getters.Any(getters.Key[uint]("$in.PostedInvoiceID")),
												}),
												Label: getters.Format("%s (#%d)", getters.Any(getters.Key[string]("$in.PostedInvoice.Number")), getters.Any(getters.Key[uint]("$in.PostedInvoice.ID"))),
												Classes: "link link-hover",
											},
										}},
										&components.LabelInline{Title: "Payment", Children: []components.PageInterface{
											&components.FieldLink{
												Href: lamu.RoutePath("finance_invoices.PaymentDetailRoute", map[string]getters.Getter[any]{
													"id": getters.Any(getters.Key[uint]("$in.PaymentID")),
												}),
												Label: settlementPaymentDetailLabelGetter(),
												Classes: "link link-hover",
											},
										}},
										&components.LabelInline{Title: "Payment date", Children: []components.PageInterface{
											&components.FieldDate{Getter: getters.Key[time.Time]("$in.Payment.Datetime")},
										}},
										&components.LabelInline{Title: "Prior partial record", Children: []components.PageInterface{
											&components.FieldText{Getter: settlementPriorPartialSummaryGetter("$in")},
										}},
									},
								},
							},
						},
					},
				},
			},
		}},
		{Key: "finance_invoices.PartiallyPaidInvoiceDetailMenu", Value: &components.SidebarMenu{
			Title: getters.Static("Partial payment"),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("Partial payments"),
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
