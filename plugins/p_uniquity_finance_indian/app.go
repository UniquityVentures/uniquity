package p_uniquity_finance_indian

import (
	"log"
	"net/url"

	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/registry"
)

const AppUrl = "/finance-indian/"

// GetPlugin returns registry contributions for [lago.BuildAllRegistries].
func GetPlugin() registry.Pair[string, lago.Plugin] {
	u, err := url.Parse(AppUrl)
	if err != nil {
		log.Panic(err)
	}
	return registry.Pair[string, lago.Plugin]{
		Key: "p_uniquity_finance_indian",
		Value: lago.Plugin{
			Type:        lago.PluginTypeAddon,
			Icon:        "map-pin",
			URL:         u,
			VerboseName: "Finance India (GST seed)",
			Roles:       []string{"superuser"},
			Migrations:  lago.PluginStages(pluginMigrations),
		},
	}
}
