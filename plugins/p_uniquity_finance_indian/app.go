package p_uniquity_finance_indian

import (
	"log"
	"net/url"

	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

const AppUrl = "/finance-indian/"

// GetPlugin returns registry contributions for [lamu.BuildAllRegistries].
func GetPlugin() registry.Pair[string, lamu.Plugin] {
	u, err := url.Parse(AppUrl)
	if err != nil {
		log.Panic(err)
	}
	return registry.Pair[string, lamu.Plugin]{
		Key: "p_uniquity_finance_indian",
		Value: lamu.Plugin{
			Type:        lamu.PluginTypeAddon,
			Icon:        "map-pin",
			URL:         u,
			VerboseName: "Finance India (GST seed)",
			Roles:       []string{"superuser"},
			Migrations:  lamu.PluginStages(pluginMigrations),
		},
	}
}
