package p_uniquity_finance_creditnotes

import (
	"log"
	"net/url"

	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/registry"
)

const AppUrl = "/finance-credit-notes/"

// GetPlugin returns registry contributions for [lago.BuildAllRegistries].
func GetPlugin() registry.Pair[string, lago.Plugin] {
	u, err := url.Parse(AppUrl)
	if err != nil {
		log.Panic(err)
	}
	return registry.Pair[string, lago.Plugin]{
		Key: "p_uniquity_finance_creditnotes",
		Value: lago.Plugin{
			Type:        lago.PluginTypeAddon,
			Icon:        "arrow-uturn-left",
			URL:         u,
			VerboseName: "Finance credit notes",
			Roles:       []string{"superuser"},
			Views:       lago.PluginStages(pluginViews),
			Pages:       lago.PluginStages(pluginPages),
			Routes:      lago.PluginStages(pluginRoutes),
			Migrations:  lago.PluginStages(pluginMigrations),
			Models:      lago.PluginStages(pluginModels),
		},
	}
}
