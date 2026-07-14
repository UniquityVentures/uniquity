package p_uniquity_finance_creditnotes

import (
	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/registry"
)

func pluginModels() lago.PluginFeatures[any] {
	return lago.PluginFeatures[any]{
		Entries: []registry.Pair[string, any]{
			{Key: "p_uniquity_finance_creditnotes.CreditNote", Value: CreditNote{}},
		},
	}
}
