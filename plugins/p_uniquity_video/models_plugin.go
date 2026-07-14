package p_uniquity_video

import (
	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/registry"
)

func pluginModels() lago.PluginFeatures[any] {
	return lago.PluginFeatures[any]{
		Entries: []registry.Pair[string, any]{
			{Key: "p_uniquity_video.RawFootage", Value: RawFootage{}},
			{Key: "p_uniquity_video.EditedVideo", Value: EditedVideo{}},
			{Key: "p_uniquity_video.PublishedVideo", Value: PublishedVideo{}},
		},
	}
}
