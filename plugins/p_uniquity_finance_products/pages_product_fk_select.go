package p_uniquity_finance_products

import (
	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/components"
	"github.com/lariv-in/lago/getters"
	"github.com/lariv-in/lago/registry"
)

func pageEntriesProductFkSelectPages() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "finance_products.ProductFkSelectionFilter", Value: &components.FormComponent[Product]{
			Attr: getters.FormBoostedGet(lago.RoutePath("finance_products.ProductFkSelectRoute", nil)),
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
		{Key: "finance_products.ProductFkSelectionTable", Value: &components.Modal{
			UID: "finance-product-fk-select-modal",
			Children: []components.PageInterface{
				&components.DataTable[Product]{
					UID:   "finance-product-fk-select-table",
					Title: "Select product",
					Data:  getters.Key[components.ObjectList[Product]]("products"),
					RowAttr: getters.RowAttrSelectNamed(
						getters.IfOrElse(getters.Key[string]("$get.target_input"), getters.Static("ProductID")),
						getters.Key[uint]("$row.ID"),
						getters.Key[string]("$row.Name"),
					),
					Actions: []components.PageInterface{
						&components.TableButtonFilter{Child: lago.DynamicPage{Name: "finance_products.ProductFkSelectionFilter"}},
					},
					Columns: []components.TableColumn{
						{Label: "Reference", Name: "Reference", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Reference")},
						}},
						{Label: "Name", Name: "Name", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Name")},
						}},
						{Label: "Sales price", Name: "SalesPrice", Children: []components.PageInterface{
							&components.FieldText{Getter: productDecimalStringGetter("$row.SalesPrice")},
						}},
					},
				},
			},
		}},
	}
}
