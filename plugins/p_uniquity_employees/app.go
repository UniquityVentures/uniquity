package p_uniquity_employees

import (
	"log"
	"net/url"

	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

const AppUrl = "/employees/"

// GetPlugin returns registry contributions for [lamu.BuildAllRegistries].
func GetPlugin() registry.Pair[string, lamu.Plugin] {
	u, err := url.Parse(AppUrl)
	if err != nil {
		log.Panic(err)
	}
	return registry.Pair[string, lamu.Plugin]{
		Key: "p_uniquity_employees",
		Value: lamu.Plugin{
			Type:        lamu.PluginTypeApp,
			Icon:        "users",
			URL:         u,
			VerboseName: "Employees & points",
			Roles:       []string{"superuser"},
			Views:       pluginViews,
			Pages:       pluginPages,
			Routes:      pluginRoutes,
			Models:      pluginModels,
			Migrations:  pluginMigrations,
			DBInitHooks: pluginDBInitHooks,
		},
	}
}
