package p_uniquity_video

import (
	"log"
	"net/url"

	"github.com/UniquityVentures/lago/lago"
)

const AppUrl = "/video/"

func init() {
	u, err := url.Parse(AppUrl)
	if err != nil {
		log.Panic(err)
	}
	if err := lago.RegistryPlugin.Register("p_uniquity_video", lago.Plugin{
		Type:        lago.PluginTypeApp,
		Icon:        "film",
		URL:         u,
		VerboseName: "Video editors",
	}); err != nil {
		log.Panic(err)
	}
}
