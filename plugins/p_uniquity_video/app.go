package p_uniquity_video

import (
	"log"
	"net/url"

	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

const AppUrl = "/video/"

// GetPlugin returns registry contributions for [lamu.BuildAllRegistries].
func GetPlugin() registry.Pair[string, lamu.Plugin] {
	u, err := url.Parse(AppUrl)
	if err != nil {
		log.Panic(err)
	}
	return registry.Pair[string, lamu.Plugin]{
		Key: "p_uniquity_video",
		Value: lamu.Plugin{
			Type:        lamu.PluginTypeApp,
			Icon:        "film",
			URL:         u,
			VerboseName: "Video editors",
			Views:       lamu.PluginStages(pluginViews),
			Pages:       lamu.PluginStages(pluginPages),
			Routes:      lamu.PluginStages(pluginRoutes),
			Models:      lamu.PluginStages(pluginModels),
			Migrations:  lamu.PluginStages(pluginMigrations),
			Configs:     lamu.PluginStages(pluginConfigs),
		},
	}
}
