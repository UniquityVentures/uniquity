package p_uniquity_finance_products

import (
	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/registry"
)

func pluginRoutes() lago.PluginFeatures[lago.Route] {
	p := AppUrl + "p/"
	return lago.PluginFeatures[lago.Route]{
		Entries: []registry.Pair[string, lago.Route]{
			{Key: "finance_products.DefaultRoute", Value: lago.Route{Path: AppUrl, Handler: lago.NewDynamicView("finance_products.ProductListView")}},
			{Key: "finance_products.ProductCreateRoute", Value: lago.Route{Path: AppUrl + "create/", Handler: lago.NewDynamicView("finance_products.ProductCreateView")}},
			{Key: "finance_products.ProductDetailRoute", Value: lago.Route{Path: p + "{id}/", Handler: lago.NewDynamicView("finance_products.ProductDetailView")}},
			{Key: "finance_products.ProductUpdateRoute", Value: lago.Route{Path: p + "{id}/edit/", Handler: lago.NewDynamicView("finance_products.ProductUpdateView")}},
			{Key: "finance_products.ProductDeleteRoute", Value: lago.Route{Path: p + "{id}/delete/", Handler: lago.NewDynamicView("finance_products.ProductDeleteView")}},
			{Key: "finance_products.ProductFkSelectRoute", Value: lago.Route{Path: AppUrl + "pick-product/", Handler: lago.NewDynamicView("finance_products.ProductFkSelectView")}},
		},
	}
}
