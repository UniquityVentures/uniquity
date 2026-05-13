package p_uniquity_video

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pluginRoutes() lamu.PluginFeatures[lamu.Route] {
	raw := AppUrl + "raw/"
	rawR := raw + "r/"
	ed := AppUrl + "edited/"
	edR := ed + "r/"
	pub := AppUrl + "published/"
	pubR := pub + "r/"
	return lamu.PluginFeatures[lamu.Route]{
		Entries: []registry.Pair[string, lamu.Route]{
			{Key: "video.DefaultRoute", Value: lamu.Route{Path: AppUrl, Handler: lamu.NewDynamicView("video.HubView")}},
			{Key: "video.RawListRoute", Value: lamu.Route{Path: raw, Handler: lamu.NewDynamicView("video.RawListView")}},
			{Key: "video.RawCreateRoute", Value: lamu.Route{Path: raw + "create/", Handler: lamu.NewDynamicView("video.RawCreateView")}},
			{Key: "video.RawSelectRoute", Value: lamu.Route{Path: raw + "select/", Handler: lamu.NewDynamicView("video.RawSelectView")}},
			{Key: "video.EmployeeSelectRoute", Value: lamu.Route{Path: raw + "select-employee/", Handler: lamu.NewDynamicView("video.EmployeeSelectView")}},
			{Key: "video.RawDetailRoute", Value: lamu.Route{Path: rawR + "{id}/", Handler: lamu.NewDynamicView("video.RawDetailView")}},
			{Key: "video.RawUpdateRoute", Value: lamu.Route{Path: rawR + "{id}/edit/", Handler: lamu.NewDynamicView("video.RawUpdateView")}},
			{Key: "video.RawDeleteRoute", Value: lamu.Route{Path: rawR + "{id}/delete/", Handler: lamu.NewDynamicView("video.RawDeleteView")}},
			{Key: "video.EditedListRoute", Value: lamu.Route{Path: ed, Handler: lamu.NewDynamicView("video.EditedListView")}},
			{Key: "video.EditedCreateRoute", Value: lamu.Route{Path: ed + "create/", Handler: lamu.NewDynamicView("video.EditedCreateView")}},
			{Key: "video.EditedSelectRoute", Value: lamu.Route{Path: ed + "select/", Handler: lamu.NewDynamicView("video.EditedSelectView")}},
			{Key: "video.EditedDetailRoute", Value: lamu.Route{Path: edR + "{id}/", Handler: lamu.NewDynamicView("video.EditedDetailView")}},
			{Key: "video.EditedUpdateRoute", Value: lamu.Route{Path: edR + "{id}/edit/", Handler: lamu.NewDynamicView("video.EditedUpdateView")}},
			{Key: "video.EditedDeleteRoute", Value: lamu.Route{Path: edR + "{id}/delete/", Handler: lamu.NewDynamicView("video.EditedDeleteView")}},
			{Key: "video.PublishedListRoute", Value: lamu.Route{Path: pub, Handler: lamu.NewDynamicView("video.PublishedListView")}},
			{Key: "video.PublishedCreateRoute", Value: lamu.Route{Path: pub + "create/", Handler: lamu.NewDynamicView("video.PublishedCreateView")}},
			{Key: "video.PublishedSelectRoute", Value: lamu.Route{Path: pub + "select/", Handler: lamu.NewDynamicView("video.PublishedSelectView")}},
			{Key: "video.PublishedDetailRoute", Value: lamu.Route{Path: pubR + "{id}/", Handler: lamu.NewDynamicView("video.PublishedDetailView")}},
			{Key: "video.PublishedEditorPointsCreateRoute", Value: lamu.Route{Path: pubR + "{id}/editor-points/", Handler: lamu.NewDynamicView("video.PublishedEditorPointsCreateView")}},
			{Key: "video.PublishedUpdateRoute", Value: lamu.Route{Path: pubR + "{id}/edit/", Handler: lamu.NewDynamicView("video.PublishedUpdateView")}},
			{Key: "video.PublishedDeleteRoute", Value: lamu.Route{Path: pubR + "{id}/delete/", Handler: lamu.NewDynamicView("video.PublishedDeleteView")}},
		},
	}
}
