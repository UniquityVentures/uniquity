package p_uniquity_finance_taxes

import (
	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/registry"
)

func pluginRoutes() lago.PluginFeatures[lago.Route] {
	t := AppUrl + "t/"
	return lago.PluginFeatures[lago.Route]{
		Entries: []registry.Pair[string, lago.Route]{
			{Key: "finance_taxes.DefaultRoute", Value: lago.Route{Path: AppUrl, Handler: lago.NewDynamicView("finance_taxes.TaxListView")}},
			{Key: "finance_taxes.TaxCreateRoute", Value: lago.Route{Path: AppUrl + "create/", Handler: lago.NewDynamicView("finance_taxes.TaxCreateView")}},
			{Key: "finance_taxes.TaxDetailRoute", Value: lago.Route{Path: t + "{id}/", Handler: lago.NewDynamicView("finance_taxes.TaxDetailView")}},
			{Key: "finance_taxes.TaxUpdateRoute", Value: lago.Route{Path: t + "{id}/edit/", Handler: lago.NewDynamicView("finance_taxes.TaxUpdateView")}},
			{Key: "finance_taxes.TaxDeleteRoute", Value: lago.Route{Path: t + "{id}/delete/", Handler: lago.NewDynamicView("finance_taxes.TaxDeleteView")}},
			{Key: "finance_taxes.TaxMultiSelectRoute", Value: lago.Route{Path: AppUrl + "multi-select/", Handler: lago.NewDynamicView("finance_taxes.TaxMultiSelectView")}},
		},
	}
}
