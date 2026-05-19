package p_uniquity_finance_taxes

import (
	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pageEntriesTaxMultiSelectPages() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "finance_taxes.TaxMultiSelectionFilter", Value: &components.FormComponent[Tax]{
			Attr: getters.FormBoostedGet(lamu.RoutePath("finance_taxes.TaxMultiSelectRoute", nil)),
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
						&components.TableButtonFilter{Child: lamu.DynamicPage{Name: "finance_taxes.TaxMultiSelectionFilter"}},
					},
					Columns: []components.TableColumn{
						{Label: "Name", Name: "Name", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Name")},
						}},
						{Label: "Type", Name: "TaxType", Children: []components.PageInterface{
							&components.FieldText{Getter: taxKindLabel("$row.TaxType")},
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
