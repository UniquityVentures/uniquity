package p_uniquity_products

import (
	"context"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type productsHubBody struct {
	components.Page
}

func (e productsHubBody) GetKey() string {
	return e.Key
}

func (e productsHubBody) GetRoles() []string {
	return e.Roles
}

func (productsHubBody) Build(ctx context.Context) Node {
	return Div(Class("p-4 prose max-w-none"),
		H1(Class("text-2xl font-bold"), Text("Products")),
		P(Class("text-base-content/70"),
			Text("Catalog items for invoice lines."),
		),
	)
}

func pluginPages() lamu.PluginFeatures[components.PageInterface] {
	return lamu.PluginFeatures[components.PageInterface]{
		Entries: []registry.Pair[string, components.PageInterface]{
			{Key: "products.MainMenu", Value: &components.SidebarMenu{
				Title: getters.Static("Products"),
				Back: &components.SidebarMenuItem{
					Title: getters.Static("Back to Home"),
					Url:   lamu.RoutePath("dashboard.AppsPage", nil),
				},
				Children: []components.PageInterface{},
			}},
			{Key: "products.HubPage", Value: &components.ShellScaffold{
				Page:    components.Page{Roles: []string{"superuser"}},
				Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "products.MainMenu"}},
				Children: []components.PageInterface{
					productsHubBody{Page: components.Page{Roles: []string{"superuser"}}},
				},
			}},
			{Key: "products.ProductSelectionTable", Value: &components.Modal{
				UID: "products-product-select-modal",
				Children: []components.PageInterface{
					&components.DataTable[Product]{
						UID:   "products-product-select-table",
						Title: "Select product",
						Data:  getters.Key[components.ObjectList[Product]]("products"),
						RowAttr: getters.RowAttrSelect("ProductID",
							getters.Key[uint]("$row.ID"),
							getters.Format("%s — %s", getters.Any(getters.Key[string]("$row.Code")), getters.Any(getters.Key[string]("$row.Name"))),
						),
						Columns: []components.TableColumn{
							{Label: "Code", Name: "Code", Children: []components.PageInterface{
								&components.FieldText{Getter: getters.Key[string]("$row.Code")},
							}},
							{Label: "Name", Name: "Name", Children: []components.PageInterface{
								&components.FieldText{Getter: getters.Key[string]("$row.Name")},
							}},
							{Label: "Entity", Name: "Entity", Children: []components.PageInterface{
								&components.FieldText{Getter: getters.Key[string]("$row.Entity.Name")},
							}},
						},
					},
				},
			}},
		},
	}
}
