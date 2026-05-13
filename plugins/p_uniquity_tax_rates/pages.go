package p_uniquity_tax_rates

import (
	"context"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/fields"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type taxRatesHubBody struct {
	components.Page
}

func (e taxRatesHubBody) GetKey() string {
	return e.Key
}

func (e taxRatesHubBody) GetRoles() []string {
	return e.Roles
}

func (taxRatesHubBody) Build(ctx context.Context) Node {
	return Div(Class("p-4 prose max-w-none"),
		H1(Class("text-2xl font-bold"), Text("Tax rates")),
		P(Class("text-base-content/70"),
			Text("Sales and purchase taxes for invoice lines."),
		),
	)
}

func pluginPages() lamu.PluginFeatures[components.PageInterface] {
	return lamu.PluginFeatures[components.PageInterface]{
		Entries: []registry.Pair[string, components.PageInterface]{
			{Key: "tax_rates.MainMenu", Value: &components.SidebarMenu{
				Title: getters.Static("Tax rates"),
				Back: &components.SidebarMenuItem{
					Title: getters.Static("Back to Home"),
					Url:   lamu.RoutePath("dashboard.AppsPage", nil),
				},
				Children: []components.PageInterface{},
			}},
			{Key: "tax_rates.HubPage", Value: &components.ShellScaffold{
				Page:    components.Page{Roles: []string{"superuser"}},
				Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "tax_rates.MainMenu"}},
				Children: []components.PageInterface{
					taxRatesHubBody{Page: components.Page{Roles: []string{"superuser"}}},
				},
			}},
			{Key: "tax_rates.TaxRateSelectionTable", Value: &components.Modal{
				UID: "tax-rates-select-modal",
				Children: []components.PageInterface{
					&components.DataTable[TaxRate]{
						UID:   "tax-rates-select-table",
						Title: "Select tax rate",
						Data:  getters.Key[components.ObjectList[TaxRate]]("taxRates"),
						RowAttr: getters.RowAttrSelectMulti(
							getters.IfOrElse(getters.Key[string]("$get.target_input"), getters.Static("Taxes")),
							getters.Key[uint]("$row.ID"),
							getters.Key[string]("$row.Name"),
						),
						Columns: []components.TableColumn{
							{Label: "Name", Name: "Name", Children: []components.PageInterface{
								&components.FieldText{Getter: getters.Key[string]("$row.Name")},
							}},
							{Label: "Scope", Name: "Scope", Children: []components.PageInterface{
								&components.FieldText{Getter: getters.Key[string]("$row.Scope")},
							}},
							{Label: "Amount", Name: "Amount", Children: []components.PageInterface{
								&components.FieldText{Getter: getters.Format("%v", getters.Any(getters.Key[fields.DecimalSix]("$row.Amount")))},
							}},
						},
					},
				},
			}},
		},
	}
}
