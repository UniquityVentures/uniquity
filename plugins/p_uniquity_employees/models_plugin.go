package p_uniquity_employees

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
	"gorm.io/gorm"
)

func pluginModels() lamu.PluginFeatures[any] {
	return lamu.PluginFeatures[any]{
		Entries: []registry.Pair[string, any]{
			{Key: "p_uniquity_employees.Employee", Value: Employee{}},
			{Key: "p_uniquity_employees.PointsTransaction", Value: PointsTransaction{}},
		},
	}
}

func pluginDBInitHooks() lamu.PluginFeatures[lamu.DBInitHook] {
	return lamu.PluginFeatures[lamu.DBInitHook]{
		Entries: []registry.Pair[string, lamu.DBInitHook]{
			{Key: "p_uniquity_employees.trigger", Value: func(d *gorm.DB) *gorm.DB {
				installPointsTransactionSuperuserTrigger(d)
				return d
			}},
		},
	}
}
