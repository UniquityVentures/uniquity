package p_uniquity_tax_rates

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pluginRoutes() lamu.PluginFeatures[lamu.Route] {
	return lamu.PluginFeatures[lamu.Route]{
		Entries: []registry.Pair[string, lamu.Route]{
			{Key: "tax_rates.DefaultRoute", Value: lamu.Route{Path: AppUrl, Handler: lamu.NewDynamicView("tax_rates.HubView")}},
			{Key: "tax_rates.TaxRateSelectRoute", Value: lamu.Route{Path: AppUrl + "select/", Handler: lamu.NewDynamicView("tax_rates.TaxRateSelectView")}},
		},
	}
}
