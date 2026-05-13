package p_uniquity_products

import (
	"log"
	"net/url"

	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

const AppUrl = "/products/"

// GetPlugin returns registry contributions for [lamu.BuildAllRegistries].
func GetPlugin() registry.Pair[string, lamu.Plugin] {
	u, err := url.Parse(AppUrl)
	if err != nil {
		log.Panic(err)
	}
	return registry.Pair[string, lamu.Plugin]{
		Key: "p_uniquity_products",
		Value: lamu.Plugin{
			Type:        lamu.PluginTypeApp,
			Icon:        "cube",
			URL:         u,
			VerboseName: "Products",
			Roles:       []string{"superuser"},
			Views:       pluginViews,
			Pages:       pluginPages,
			Routes:      pluginRoutes,
			Models:      pluginModels,
			Migrations:  pluginMigrations,
		},
	}
}
