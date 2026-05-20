package p_uniquity_finance_products

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pluginModels() lamu.PluginFeatures[any] {
	return lamu.PluginFeatures[any]{
		Entries: []registry.Pair[string, any]{
			{Key: "p_uniquity_finance_products.Product", Value: Product{}},
			{Key: "p_uniquity_finance_products.ProductPreferences", Value: ProductPreferences{}},
		},
	}
}
