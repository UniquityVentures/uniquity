package p_uniquity_finance_customer

import (
	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/registry"
)

func pluginRoutes() lago.PluginFeatures[lago.Route] {
	cust := AppUrl + "c/"
	return lago.PluginFeatures[lago.Route]{
		Entries: []registry.Pair[string, lago.Route]{
			{Key: "finance_customers.DefaultRoute", Value: lago.Route{Path: AppUrl, Handler: lago.NewDynamicView("finance_customers.CustomerListView")}},
			{Key: "finance_customers.CustomerCreateRoute", Value: lago.Route{Path: AppUrl + "create/", Handler: lago.NewDynamicView("finance_customers.CustomerCreateView")}},
			{Key: "finance_customers.CustomerDetailRoute", Value: lago.Route{Path: cust + "{id}/", Handler: lago.NewDynamicView("finance_customers.CustomerDetailView")}},
			{Key: "finance_customers.CustomerUpdateRoute", Value: lago.Route{Path: cust + "{id}/edit/", Handler: lago.NewDynamicView("finance_customers.CustomerUpdateView")}},
			{Key: "finance_customers.CustomerDeleteRoute", Value: lago.Route{Path: cust + "{id}/delete/", Handler: lago.NewDynamicView("finance_customers.CustomerDeleteView")}},
			{Key: "finance_customers.CustomerFkSelectRoute", Value: lago.Route{Path: AppUrl + "pick-customer/", Handler: lago.NewDynamicView("finance_customers.CustomerFkSelectView")}},
		},
	}
}
