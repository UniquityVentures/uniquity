package p_uniquity_video

import "github.com/UniquityVentures/lago/lago"

func init() {
	registerHubRoute()
	registerRawRoutes()
	registerEditedRoutes()
	registerPublishedRoutes()
}

func registerHubRoute() {
	_ = lago.RegistryRoute.Register("video.DefaultRoute", lago.Route{
		Path:    AppUrl,
		Handler: lago.NewDynamicView("video.HubView"),
	})
}

func registerRawRoutes() {
	raw := AppUrl + "raw/"
	rawR := raw + "r/"
	_ = lago.RegistryRoute.Register("video.RawListRoute", lago.Route{
		Path:    raw,
		Handler: lago.NewDynamicView("video.RawListView"),
	})
	_ = lago.RegistryRoute.Register("video.RawCreateRoute", lago.Route{
		Path:    raw + "create/",
		Handler: lago.NewDynamicView("video.RawCreateView"),
	})
	_ = lago.RegistryRoute.Register("video.RawSelectRoute", lago.Route{
		Path:    raw + "select/",
		Handler: lago.NewDynamicView("video.RawSelectView"),
	})
	_ = lago.RegistryRoute.Register("video.EmployeeSelectRoute", lago.Route{
		Path:    raw + "select-employee/",
		Handler: lago.NewDynamicView("video.EmployeeSelectView"),
	})
	_ = lago.RegistryRoute.Register("video.RawDetailRoute", lago.Route{
		Path:    rawR + "{id}/",
		Handler: lago.NewDynamicView("video.RawDetailView"),
	})
	_ = lago.RegistryRoute.Register("video.RawUpdateRoute", lago.Route{
		Path:    rawR + "{id}/edit/",
		Handler: lago.NewDynamicView("video.RawUpdateView"),
	})
	_ = lago.RegistryRoute.Register("video.RawDeleteRoute", lago.Route{
		Path:    rawR + "{id}/delete/",
		Handler: lago.NewDynamicView("video.RawDeleteView"),
	})
}

func registerEditedRoutes() {
	ed := AppUrl + "edited/"
	edR := ed + "r/"
	_ = lago.RegistryRoute.Register("video.EditedListRoute", lago.Route{
		Path:    ed,
		Handler: lago.NewDynamicView("video.EditedListView"),
	})
	_ = lago.RegistryRoute.Register("video.EditedCreateRoute", lago.Route{
		Path:    ed + "create/",
		Handler: lago.NewDynamicView("video.EditedCreateView"),
	})
	_ = lago.RegistryRoute.Register("video.EditedSelectRoute", lago.Route{
		Path:    ed + "select/",
		Handler: lago.NewDynamicView("video.EditedSelectView"),
	})
	_ = lago.RegistryRoute.Register("video.EditedDetailRoute", lago.Route{
		Path:    edR + "{id}/",
		Handler: lago.NewDynamicView("video.EditedDetailView"),
	})
	_ = lago.RegistryRoute.Register("video.EditedUpdateRoute", lago.Route{
		Path:    edR + "{id}/edit/",
		Handler: lago.NewDynamicView("video.EditedUpdateView"),
	})
	_ = lago.RegistryRoute.Register("video.EditedDeleteRoute", lago.Route{
		Path:    edR + "{id}/delete/",
		Handler: lago.NewDynamicView("video.EditedDeleteView"),
	})
}

func registerPublishedRoutes() {
	pub := AppUrl + "published/"
	pubR := pub + "r/"
	_ = lago.RegistryRoute.Register("video.PublishedListRoute", lago.Route{
		Path:    pub,
		Handler: lago.NewDynamicView("video.PublishedListView"),
	})
	_ = lago.RegistryRoute.Register("video.PublishedCreateRoute", lago.Route{
		Path:    pub + "create/",
		Handler: lago.NewDynamicView("video.PublishedCreateView"),
	})
	_ = lago.RegistryRoute.Register("video.PublishedSelectRoute", lago.Route{
		Path:    pub + "select/",
		Handler: lago.NewDynamicView("video.PublishedSelectView"),
	})
	_ = lago.RegistryRoute.Register("video.PublishedDetailRoute", lago.Route{
		Path:    pubR + "{id}/",
		Handler: lago.NewDynamicView("video.PublishedDetailView"),
	})
	_ = lago.RegistryRoute.Register("video.PublishedEditorPointsCreateRoute", lago.Route{
		Path:    pubR + "{id}/editor-points/",
		Handler: lago.NewDynamicView("video.PublishedEditorPointsCreateView"),
	})
	_ = lago.RegistryRoute.Register("video.PublishedUpdateRoute", lago.Route{
		Path:    pubR + "{id}/edit/",
		Handler: lago.NewDynamicView("video.PublishedUpdateView"),
	})
	_ = lago.RegistryRoute.Register("video.PublishedDeleteRoute", lago.Route{
		Path:    pubR + "{id}/delete/",
		Handler: lago.NewDynamicView("video.PublishedDeleteView"),
	})
}
