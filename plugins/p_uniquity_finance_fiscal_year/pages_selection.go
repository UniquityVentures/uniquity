package p_uniquity_finance_fiscal_year

import (
	"time"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pageEntriesFiscalYearSelectionPages() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "finance_fiscal_years.FiscalYearSelectionFilter", Value: &components.FormComponent[FiscalYear]{
			Attr: getters.FormBoostedGet(lamu.RoutePath("finance_fiscal_years.FiscalYearSelectRoute", nil)),
			ChildrenInput: []components.PageInterface{
				&components.InputText{Label: "Code", Name: "Code", Getter: getters.Key[string]("$get.Code")},
				&components.InputText{Label: "Name", Name: "Name", Getter: getters.Key[string]("$get.Name")},
			},
			ChildrenAction: []components.PageInterface{
				&components.ContainerRow{Classes: "flex gap-2", Children: []components.PageInterface{
					&components.ButtonSubmit{Label: "Apply"},
					&components.ButtonClear{Label: "Clear"},
				}},
			},
		}},
		{Key: "finance_fiscal_years.FiscalYearSelectionTable", Value: &components.Modal{
			UID: "finance-fiscal-year-select-modal",
			Children: []components.PageInterface{
				&components.DataTable[FiscalYear]{
					UID:   "finance-fiscal-year-select-table",
					Title: "Select fiscal year",
					Data:  getters.Key[components.ObjectList[FiscalYear]]("fiscal_years"),
					Actions: []components.PageInterface{
						&components.TableButtonFilter{Child: lamu.DynamicPage{Name: "finance_fiscal_years.FiscalYearSelectionFilter"}},
					},
					RowAttr: getters.RowAttrSelect("FiscalYearID",
						getters.Key[uint]("$row.ID"),
						getters.Key[string]("$row.Name"),
					),
					Columns: []components.TableColumn{
						{Label: "Code", Name: "Code", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Code")},
						}},
						{Label: "Name", Name: "Name", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Name")},
						}},
						{Label: "Start", Name: "Start", Children: []components.PageInterface{
							&components.FieldDatetime{Getter: getters.Key[time.Time]("$row.Start")},
						}},
						{Label: "End", Name: "End", Children: []components.PageInterface{
							&components.FieldDatetime{Getter: getters.Key[time.Time]("$row.End")},
						}},
					},
				},
			},
		}},
	}
}
