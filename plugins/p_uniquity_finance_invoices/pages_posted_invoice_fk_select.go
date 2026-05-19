package p_uniquity_finance_invoices

import (
	"context"
	"fmt"
	"time"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/registry"
)

func pageEntriesPostedInvoiceFkSelectPages() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "finance_invoices.PostedInvoiceFkSelectionTable", Value: &components.Modal{
			UID: "finance-posted-invoice-fk-modal",
			Children: []components.PageInterface{
				&components.DataTable[PostedInvoice]{
					UID:   "finance-posted-invoice-fk-select-table",
					Title: "Select posted invoice",
					Data:  getters.Key[components.ObjectList[PostedInvoice]]("posted_invoices"),
					RowAttr: getters.RowAttrSelect("PostedInvoiceID",
						getters.Key[uint]("$row.ID"),
						postedInvoiceFkRowSummaryGetter(),
					),
					Columns: []components.TableColumn{
						{Label: "Number", Name: "Number", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Number")},
						}},
						{Label: "Customer", Name: "Customer", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Customer.Name")},
						}},
						{Label: "Invoice date", Name: "Datetime", Children: []components.PageInterface{
							&components.FieldDate{Getter: getters.Key[time.Time]("$row.Datetime")},
						}},
					},
				},
			},
		}},
	}
}

func postedInvoiceFkRowSummaryGetter() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		id, err := getters.Key[uint]("$row.ID")(ctx)
		if err != nil {
			return "", err
		}
		num, err := getters.Key[string]("$row.Number")(ctx)
		if err != nil {
			return fmt.Sprintf("#%d", id), nil
		}
		return fmt.Sprintf("%s (#%d)", num, id), nil
	}
}
