package p_uniquity_video

import (
	"log"
	"net/url"

	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/registry"
)

const AppUrl = "/video/"

// GetPlugin returns registry contributions for [lago.BuildAllRegistries].
func GetPlugin() registry.Pair[string, lago.Plugin] {
	u, err := url.Parse(AppUrl)
	if err != nil {
		log.Panic(err)
	}
	return registry.Pair[string, lago.Plugin]{
		Key: "p_uniquity_video",
		Value: lago.Plugin{
			Type:        lago.PluginTypeApp,
			Icon:        "film",
			URL:         u,
			VerboseName: "Video editors",
			Views:       lago.PluginStages(pluginViews),
			Pages:       lago.PluginStages(pluginPages),
			Routes:      lago.PluginStages(pluginRoutes),
			Models:      lago.PluginStages(pluginModels),
			Migrations:  lago.PluginStages(pluginMigrations),
			Configs:     lago.PluginStages(pluginConfigs),
		},
	}
}
