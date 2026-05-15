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
			Type:        lamu.PluginTypeApp,
			Icon:        "building-storefront",
			URL:         u,
			VerboseName: "Finance customers",
			Roles:       []string{"superuser"},
			Views:       pluginViews,
			Pages:       pluginPages,
			Routes:      pluginRoutes,
			Migrations:  pluginMigrations,
			Models:      pluginModels,
		},
	}
}
