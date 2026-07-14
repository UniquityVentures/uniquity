package p_uniquity_video

import (
	"strings"

	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/registry"
)

// PluginConfig is loaded from TOML [Plugins.p_uniquity_video].
//
// youtubeApiKey is required: YouTube Data API v3 is used on published video detail
// (title, publish date, upload status). Create a key in Google Cloud Console and
// enable the YouTube Data API v3 for the project.
type PluginConfig struct {
	YouTubeAPIKey string `toml:"youtubeApiKey"`
}

// VideoPluginConfig is the decoded [Plugins.p_uniquity_video] block.
var VideoPluginConfig = &PluginConfig{}

func (c *PluginConfig) PostConfig() {
	if c == nil {
		return
	}
	c.YouTubeAPIKey = strings.TrimSpace(c.YouTubeAPIKey)
	requireYouTubeAPIKey(c.YouTubeAPIKey)
}

func requireYouTubeAPIKey(apiKey string) {
	if apiKey == "" {
		panic("p_uniquity_video: [Plugins.p_uniquity_video] youtubeApiKey is required but empty or missing")
	}
}

func pluginConfigs() lago.PluginFeatures[lago.Config] {
	return lago.PluginFeatures[lago.Config]{
		Entries: []registry.Pair[string, lago.Config]{
			{Key: "p_uniquity_video", Value: VideoPluginConfig},
		},
	}
}
