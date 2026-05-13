package p_uniquity_video

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pluginModels() lamu.PluginFeatures[any] {
	return lamu.PluginFeatures[any]{
		Entries: []registry.Pair[string, any]{
			{Key: "p_uniquity_video.RawFootage", Value: RawFootage{}},
			{Key: "p_uniquity_video.EditedVideo", Value: EditedVideo{}},
			{Key: "p_uniquity_video.PublishedVideo", Value: PublishedVideo{}},
		},
	}
}
