package p_uniquity_finance_products

import (
	"log"
	"net/url"

	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/registry"
)

const AppUrl = "/finance-products/"

// GetPlugin returns registry contributions for [lago.BuildAllRegistries].
func GetPlugin() registry.Pair[string, lago.Plugin] {
	u, err := url.Parse(AppUrl)
	if err != nil {
		log.Panic(err)
	}
	return registry.Pair[string, lago.Plugin]{
		Key: "p_uniquity_finance_products",
		Value: lago.Plugin{
			Type:        lago.PluginTypeAddon,
			Icon:        "cube",
			URL:         u,
			VerboseName: "Finance products",
			Roles:       []string{"superuser"},
			Views:       lago.PluginStages(pluginViews),
			Pages:       lago.PluginStages(pluginPages),
			Routes:      lago.PluginStages(pluginRoutes),
			Migrations:  lago.PluginStages(pluginMigrations),
			Models:      lago.PluginStages(pluginModels),
		},
	}
}
