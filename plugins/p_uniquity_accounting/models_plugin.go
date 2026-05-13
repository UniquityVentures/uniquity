package p_uniquity_accounting

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pluginModels() lamu.PluginFeatures[any] {
	return lamu.PluginFeatures[any]{
		Entries: []registry.Pair[string, any]{
			{Key: "p_uniquity_accounting.Account", Value: Account{}},
			{Key: "p_uniquity_accounting.Posting", Value: Posting{}},
		},
	}
}
