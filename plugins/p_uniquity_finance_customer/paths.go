package p_uniquity_finance_customer

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pluginRoutes() lamu.PluginFeatures[lamu.Route] {
	cust := AppUrl + "c/"
	return lamu.PluginFeatures[lamu.Route]{
		Entries: []registry.Pair[string, lamu.Route]{
			{Key: "finance_customers.DefaultRoute", Value: lamu.Route{Path: AppUrl, Handler: lamu.NewDynamicView("finance_customers.CustomerListView")}},
			{Key: "finance_customers.CustomerCreateRoute", Value: lamu.Route{Path: AppUrl + "create/", Handler: lamu.NewDynamicView("finance_customers.CustomerCreateView")}},
			{Key: "finance_customers.CustomerDetailRoute", Value: lamu.Route{Path: cust + "{id}/", Handler: lamu.NewDynamicView("finance_customers.CustomerDetailView")}},
			{Key: "finance_customers.CustomerUpdateRoute", Value: lamu.Route{Path: cust + "{id}/edit/", Handler: lamu.NewDynamicView("finance_customers.CustomerUpdateView")}},
			{Key: "finance_customers.CustomerDeleteRoute", Value: lamu.Route{Path: cust + "{id}/delete/", Handler: lamu.NewDynamicView("finance_customers.CustomerDeleteView")}},
		},
	}
}
