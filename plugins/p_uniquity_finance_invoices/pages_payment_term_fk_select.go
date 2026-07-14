package p_uniquity_finance_invoices

import (
	"github.com/lariv-in/lago/components"
	"github.com/lariv-in/lago/getters"
	"github.com/lariv-in/lago/registry"
)

func pageEntriesPaymentTermFkSelectPages() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "finance_invoices.PaymentTermFkSelectionTable", Value: &components.Modal{
			UID: "finance-invoice-payment-term-fk-modal",
			Children: []components.PageInterface{
				&components.DataTable[PaymentTerm]{
					UID:   "finance-payment-term-fk-select-table",
					Title: "Select payment term",
					Data:  getters.Key[components.ObjectList[PaymentTerm]]("payment_terms"),
					RowAttr: getters.RowAttrSelect(
						"PaymentTermID",
						getters.Key[uint]("$row.ID"),
						paymentTermRowSummaryGetter(),
					),
					Columns: []components.TableColumn{
						{Label: "Kind", Name: "Type", Children: []components.PageInterface{
							&components.FieldText{Getter: registry.PairValueFromKey(getters.Key[string]("$row.Type"), paymentTermKindChoiceList)},
						}},
						{Label: "Summary", Name: "Summary", Children: []components.PageInterface{
							&components.FieldText{Getter: paymentTermRowSummaryGetter()},
						}},
					},
				},
			},
		}},
	}
}
