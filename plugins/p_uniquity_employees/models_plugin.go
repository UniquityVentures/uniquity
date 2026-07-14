package p_uniquity_employees

import (
	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/registry"
	"gorm.io/gorm"
)

func pluginModels() lago.PluginFeatures[any] {
	return lago.PluginFeatures[any]{
		Entries: []registry.Pair[string, any]{
			{Key: "p_uniquity_employees.Employee", Value: Employee{}},
			{Key: "p_uniquity_employees.PointsTransaction", Value: PointsTransaction{}},
		},
	}
}

func pluginDBInitHooks() lago.PluginFeatures[lago.DBInitHook] {
	return lago.PluginFeatures[lago.DBInitHook]{
		Entries: []registry.Pair[string, lago.DBInitHook]{
			{Key: "p_uniquity_employees.trigger", Value: func(d *gorm.DB) *gorm.DB {
				installPointsTransactionSuperuserTrigger(d)
				return d
			}},
		},
	}
}
