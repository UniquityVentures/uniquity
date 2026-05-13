package p_uniquity_products

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pluginRoutes() lamu.PluginFeatures[lamu.Route] {
	return lamu.PluginFeatures[lamu.Route]{
		Entries: []registry.Pair[string, lamu.Route]{
			{Key: "products.DefaultRoute", Value: lamu.Route{Path: AppUrl, Handler: lamu.NewDynamicView("products.HubView")}},
			{Key: "products.ProductSelectRoute", Value: lamu.Route{Path: AppUrl + "select/", Handler: lamu.NewDynamicView("products.ProductSelectView")}},
		},
	}
}
