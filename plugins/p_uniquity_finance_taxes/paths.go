package p_uniquity_finance_taxes

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pluginRoutes() lamu.PluginFeatures[lamu.Route] {
	t := AppUrl + "t/"
	return lamu.PluginFeatures[lamu.Route]{
		Entries: []registry.Pair[string, lamu.Route]{
			{Key: "finance_taxes.DefaultRoute", Value: lamu.Route{Path: AppUrl, Handler: lamu.NewDynamicView("finance_taxes.TaxListView")}},
			{Key: "finance_taxes.TaxCreateRoute", Value: lamu.Route{Path: AppUrl + "create/", Handler: lamu.NewDynamicView("finance_taxes.TaxCreateView")}},
			{Key: "finance_taxes.TaxDetailRoute", Value: lamu.Route{Path: t + "{id}/", Handler: lamu.NewDynamicView("finance_taxes.TaxDetailView")}},
			{Key: "finance_taxes.TaxUpdateRoute", Value: lamu.Route{Path: t + "{id}/edit/", Handler: lamu.NewDynamicView("finance_taxes.TaxUpdateView")}},
			{Key: "finance_taxes.TaxDeleteRoute", Value: lamu.Route{Path: t + "{id}/delete/", Handler: lamu.NewDynamicView("finance_taxes.TaxDeleteView")}},
			{Key: "finance_taxes.TaxMultiSelectRoute", Value: lamu.Route{Path: AppUrl + "multi-select/", Handler: lamu.NewDynamicView("finance_taxes.TaxMultiSelectView")}},
		},
	}
}
