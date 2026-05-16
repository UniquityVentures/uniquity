package p_uniquity_finance_customer

import (
	"log"
	"net/url"

	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

const AppUrl = "/finance-customers/"

// GetPlugin returns registry contributions for [lamu.BuildAllRegistries].
func GetPlugin() registry.Pair[string, lamu.Plugin] {
	u, err := url.Parse(AppUrl)
	if err != nil {
		log.Panic(err)
	}
	return registry.Pair[string, lamu.Plugin]{
		Key: "p_uniquity_finance_customer",
		Value: lamu.Plugin{
			Type:        lamu.PluginTypeAddon,
			Icon:        "building-storefront",
			URL:         u,
			VerboseName: "Finance customers",
			Roles:       []string{"superuser"},
			Views:       lamu.PluginStages(pluginViews),
			Pages:       lamu.PluginStages(pluginPages),
			Routes:      lamu.PluginStages(pluginRoutes),
			Migrations:  lamu.PluginStages(pluginMigrations),
			Models:      lamu.PluginStages(pluginModels),
		},
	}
}
