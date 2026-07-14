package p_uniquity_finance_fiscal_year

import (
	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/registry"
)

func pluginModels() lago.PluginFeatures[any] {
	return lago.PluginFeatures[any]{
		Entries: []registry.Pair[string, any]{
			{Key: "p_uniquity_finance_fiscal_year.FiscalYear", Value: FiscalYear{}},
		},
	}
}
