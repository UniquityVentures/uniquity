package p_uniquity_finance_accounts

import (
	"log"
	"net/url"

	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/registry"
)

const AppUrl = "/finance/"

// GetPlugin returns registry contributions for [lago.BuildAllRegistries].
func GetPlugin() registry.Pair[string, lago.Plugin] {
	u, err := url.Parse(AppUrl)
	if err != nil {
		log.Panic(err)
	}
	return registry.Pair[string, lago.Plugin]{
		Key: "p_uniquity_finance_accounts",
		Value: lago.Plugin{
			Type:        lago.PluginTypeApp,
			Icon:        "building-library",
			URL:         u,
			VerboseName: "Finance",
			Roles:       []string{"superuser"},
			Views:       lago.PluginStages(pluginViews),
			Pages:       lago.PluginStages(pluginPages),
			Routes:      lago.PluginStages(pluginRoutes),
			Models:      lago.PluginStages(pluginModels),
			Migrations:  lago.PluginStages(pluginMigrations),
		},
	}
}
