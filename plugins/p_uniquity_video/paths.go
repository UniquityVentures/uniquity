package p_uniquity_video

import (
	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/registry"
)

func pluginRoutes() lago.PluginFeatures[lago.Route] {
	raw := AppUrl + "raw/"
	rawR := raw + "r/"
	ed := AppUrl + "edited/"
	edR := ed + "r/"
	pub := AppUrl + "published/"
	pubR := pub + "r/"
	return lago.PluginFeatures[lago.Route]{
		Entries: []registry.Pair[string, lago.Route]{
			{Key: "video.DefaultRoute", Value: lago.Route{Path: AppUrl, Handler: lago.NewDynamicView("video.HubView")}},
			{Key: "video.RawListRoute", Value: lago.Route{Path: raw, Handler: lago.NewDynamicView("video.RawListView")}},
			{Key: "video.RawCreateRoute", Value: lago.Route{Path: raw + "create/", Handler: lago.NewDynamicView("video.RawCreateView")}},
			{Key: "video.RawSelectRoute", Value: lago.Route{Path: raw + "select/", Handler: lago.NewDynamicView("video.RawSelectView")}},
			{Key: "video.EmployeeSelectRoute", Value: lago.Route{Path: raw + "select-employee/", Handler: lago.NewDynamicView("video.EmployeeSelectView")}},
			{Key: "video.RawDetailRoute", Value: lago.Route{Path: rawR + "{id}/", Handler: lago.NewDynamicView("video.RawDetailView")}},
			{Key: "video.RawUpdateRoute", Value: lago.Route{Path: rawR + "{id}/edit/", Handler: lago.NewDynamicView("video.RawUpdateView")}},
			{Key: "video.RawDeleteRoute", Value: lago.Route{Path: rawR + "{id}/delete/", Handler: lago.NewDynamicView("video.RawDeleteView")}},
			{Key: "video.EditedListRoute", Value: lago.Route{Path: ed, Handler: lago.NewDynamicView("video.EditedListView")}},
			{Key: "video.EditedCreateRoute", Value: lago.Route{Path: ed + "create/", Handler: lago.NewDynamicView("video.EditedCreateView")}},
			{Key: "video.EditedSelectRoute", Value: lago.Route{Path: ed + "select/", Handler: lago.NewDynamicView("video.EditedSelectView")}},
			{Key: "video.EditedDetailRoute", Value: lago.Route{Path: edR + "{id}/", Handler: lago.NewDynamicView("video.EditedDetailView")}},
			{Key: "video.EditedUpdateRoute", Value: lago.Route{Path: edR + "{id}/edit/", Handler: lago.NewDynamicView("video.EditedUpdateView")}},
			{Key: "video.EditedDeleteRoute", Value: lago.Route{Path: edR + "{id}/delete/", Handler: lago.NewDynamicView("video.EditedDeleteView")}},
			{Key: "video.PublishedListRoute", Value: lago.Route{Path: pub, Handler: lago.NewDynamicView("video.PublishedListView")}},
			{Key: "video.PublishedCreateRoute", Value: lago.Route{Path: pub + "create/", Handler: lago.NewDynamicView("video.PublishedCreateView")}},
			{Key: "video.PublishedSelectRoute", Value: lago.Route{Path: pub + "select/", Handler: lago.NewDynamicView("video.PublishedSelectView")}},
			{Key: "video.PublishedDetailRoute", Value: lago.Route{Path: pubR + "{id}/", Handler: lago.NewDynamicView("video.PublishedDetailView")}},
			{Key: "video.PublishedEditorPointsCreateRoute", Value: lago.Route{Path: pubR + "{id}/editor-points/", Handler: lago.NewDynamicView("video.PublishedEditorPointsCreateView")}},
			{Key: "video.PublishedUpdateRoute", Value: lago.Route{Path: pubR + "{id}/edit/", Handler: lago.NewDynamicView("video.PublishedUpdateView")}},
			{Key: "video.PublishedDeleteRoute", Value: lago.Route{Path: pubR + "{id}/delete/", Handler: lago.NewDynamicView("video.PublishedDeleteView")}},
		},
	}
}
