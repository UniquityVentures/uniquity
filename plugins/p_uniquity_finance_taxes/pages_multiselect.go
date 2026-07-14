package p_uniquity_finance_taxes

import (
	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/components"
	"github.com/lariv-in/lago/getters"
	"github.com/lariv-in/lago/registry"
)

func pageEntriesTaxMultiSelectPages() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "finance_taxes.TaxMultiSelectionFilter", Value: &components.FormComponent[Tax]{
			Attr: getters.FormBoostedGet(lago.RoutePath("finance_taxes.TaxMultiSelectRoute", nil)),
			ChildrenInput: []components.PageInterface{
				&components.InputText{Hidden: true, Name: "target_input", Getter: getters.Key[string]("$get.target_input")},
				&components.InputText{Label: "Name", Name: "Name", Getter: getters.Key[string]("$get.Name")},
			},
			ChildrenAction: []components.PageInterface{
				&components.ContainerRow{Classes: "flex gap-2", Children: []components.PageInterface{
					&components.ButtonSubmit{Label: "Apply"},
					&components.ButtonClear{Label: "Clear"},
				}},
			},
		}},
		{Key: "finance_taxes.TaxMultiSelectionTable", Value: &components.Modal{
			UID: "finance-tax-multi-selection-modal",
			Children: []components.PageInterface{
				&components.DataTable[Tax]{
					UID:   "finance-tax-multi-selection-table",
					Title: "Select taxes",
					Data:  getters.Key[components.ObjectList[Tax]]("taxes"),
					RowAttr: getters.RowAttrSelectMulti(
						getters.IfOrElse(
							getters.Key[string]("$get.target_input"),
							getters.Static("Taxes"),
						),
						getters.Key[uint]("$row.ID"),
						getters.Key[string]("$row.Name"),
					),
					Actions: []components.PageInterface{
						&components.TableButtonFilter{Child: lago.DynamicPage{Name: "finance_taxes.TaxMultiSelectionFilter"}},
					},
					Columns: []components.TableColumn{
						{Label: "Name", Name: "Name", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Name")},
						}},
						{Label: "Type", Name: "TaxType", Children: []components.PageInterface{
							&components.FieldText{Getter: registry.PairValueFromKey(getters.Key[TaxKind]("$row.TaxType"), taxKindChoiceList)},
						}},
						{Label: "Percentage", Name: "Percentage", Children: []components.PageInterface{
							&components.FieldText{Getter: taxDecimalStringGetter("$row.Percentage")},
						}},
					},
				},
			},
		}},
	}
}
