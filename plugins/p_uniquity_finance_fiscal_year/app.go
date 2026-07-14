package p_uniquity_finance_fiscal_year

import (
	"log"
	"net/url"

	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/registry"
)

const AppUrl = "/finance-fiscal-years/"

// GetPlugin returns registry contributions for [lago.BuildAllRegistries].
func GetPlugin() registry.Pair[string, lago.Plugin] {
	u, err := url.Parse(AppUrl)
	if err != nil {
		log.Panic(err)
	}
	return registry.Pair[string, lago.Plugin]{
		Key: "p_uniquity_finance_fiscal_year",
		Value: lago.Plugin{
			Type:        lago.PluginTypeAddon,
			Icon:        "calendar-days",
			URL:         u,
			VerboseName: "Finance fiscal years",
			Roles:       []string{"superuser"},
			Views:       lago.PluginStages(pluginViews),
			Pages:       lago.PluginStages(pluginPages),
			Routes:      lago.PluginStages(pluginRoutes),
			Migrations:  lago.PluginStages(pluginMigrations),
			Models:      lago.PluginStages(pluginModels),
		},
	}
}
