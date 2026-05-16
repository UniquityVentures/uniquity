package p_uniquity_finance_products

import (
	"log"
	"net/url"

	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

const AppUrl = "/finance-products/"

// GetPlugin returns registry contributions for [lamu.BuildAllRegistries].
func GetPlugin() registry.Pair[string, lamu.Plugin] {
	u, err := url.Parse(AppUrl)
	if err != nil {
		log.Panic(err)
	}
	return registry.Pair[string, lamu.Plugin]{
		Key: "p_uniquity_finance_products",
		Value: lamu.Plugin{
			Type:        lamu.PluginTypeAddon,
			Icon:        "cube",
			URL:         u,
			VerboseName: "Finance products",
			Roles:       []string{"superuser"},
			Views:       lamu.PluginStages(pluginViews),
			Pages:       lamu.PluginStages(pluginPages),
			Routes:      lamu.PluginStages(pluginRoutes),
			Migrations:  lamu.PluginStages(pluginMigrations),
			Models:      lamu.PluginStages(pluginModels),
		},
	}
}
