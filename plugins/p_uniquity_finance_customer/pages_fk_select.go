package p_uniquity_finance_customer

import (
	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/components"
	"github.com/lariv-in/lago/getters"
	"github.com/lariv-in/lago/registry"
)

func pageEntriesCustomerFkSelectPages() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "finance_customers.CustomerFkSelectionFilter", Value: &components.FormComponent[Customer]{
			Attr: getters.FormBoostedGet(lago.RoutePath("finance_customers.CustomerFkSelectRoute", nil)),
			ChildrenInput: []components.PageInterface{
				&components.InputText{Label: "Name", Name: "Name", Getter: getters.Key[string]("$get.Name")},
			},
			ChildrenAction: []components.PageInterface{
				&components.ContainerRow{Classes: "flex gap-2", Children: []components.PageInterface{
					&components.ButtonSubmit{Label: "Apply"},
					&components.ButtonClear{Label: "Clear"},
				}},
			},
		}},
		{Key: "finance_customers.CustomerFkSelectionTable", Value: &components.Modal{
			UID: "finance-customer-fk-select-modal",
			Children: []components.PageInterface{
				&components.DataTable[Customer]{
					UID:   "finance-customer-fk-select-table",
					Title: "Select customer",
					Data:  getters.Key[components.ObjectList[Customer]]("customers"),
					RowAttr: getters.RowAttrSelect(
						"CustomerID",
						getters.Key[uint]("$row.ID"),
						getters.Key[string]("$row.Name"),
					),
					Actions: []components.PageInterface{
						&components.TableButtonFilter{Child: lago.DynamicPage{Name: "finance_customers.CustomerFkSelectionFilter"}},
					},
					Columns: []components.TableColumn{
						{Label: "Name", Name: "Name", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Name")},
						}},
						{Label: "Email", Name: "Email", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Email")},
						}},
						{Label: "Phone", Name: "Phone", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Phone")},
						}},
					},
				},
			},
		}},
	}
}
