package p_uniquity_finance_products

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pluginRoutes() lamu.PluginFeatures[lamu.Route] {
	p := AppUrl + "p/"
	return lamu.PluginFeatures[lamu.Route]{
		Entries: []registry.Pair[string, lamu.Route]{
			{Key: "finance_products.DefaultRoute", Value: lamu.Route{Path: AppUrl, Handler: lamu.NewDynamicView("finance_products.ProductListView")}},
			{Key: "finance_products.ProductCreateRoute", Value: lamu.Route{Path: AppUrl + "create/", Handler: lamu.NewDynamicView("finance_products.ProductCreateView")}},
			{Key: "finance_products.ProductDetailRoute", Value: lamu.Route{Path: p + "{id}/", Handler: lamu.NewDynamicView("finance_products.ProductDetailView")}},
			{Key: "finance_products.ProductUpdateRoute", Value: lamu.Route{Path: p + "{id}/edit/", Handler: lamu.NewDynamicView("finance_products.ProductUpdateView")}},
			{Key: "finance_products.ProductDeleteRoute", Value: lamu.Route{Path: p + "{id}/delete/", Handler: lamu.NewDynamicView("finance_products.ProductDeleteView")}},
			{Key: "finance_products.ProductFkSelectRoute", Value: lamu.Route{Path: AppUrl + "pick-product/", Handler: lamu.NewDynamicView("finance_products.ProductFkSelectView")}},
		},
	}
}
