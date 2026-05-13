package p_uniquity_currencies

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pluginRoutes() lamu.PluginFeatures[lamu.Route] {
	base := AppUrl + "c/"
	return lamu.PluginFeatures[lamu.Route]{
		Entries: []registry.Pair[string, lamu.Route]{
			{Key: "currencies.CurrencyListRoute", Value: lamu.Route{Path: AppUrl, Handler: lamu.NewDynamicView("currencies.CurrencyListView")}},
			{Key: "currencies.CurrencySelectRoute", Value: lamu.Route{Path: AppUrl + "select/", Handler: lamu.NewDynamicView("currencies.CurrencySelectView")}},
			{Key: "currencies.CurrencyCreateRoute", Value: lamu.Route{Path: AppUrl + "create/", Handler: lamu.NewDynamicView("currencies.CurrencyCreateView")}},
			{Key: "currencies.CurrencyDetailRoute", Value: lamu.Route{Path: base + "{id}/", Handler: lamu.NewDynamicView("currencies.CurrencyDetailView")}},
			{Key: "currencies.CurrencyUpdateRoute", Value: lamu.Route{Path: base + "{id}/edit/", Handler: lamu.NewDynamicView("currencies.CurrencyUpdateView")}},
			{Key: "currencies.CurrencyDeleteRoute", Value: lamu.Route{Path: base + "{id}/delete/", Handler: lamu.NewDynamicView("currencies.CurrencyDeleteView")}},
		},
	}
}
