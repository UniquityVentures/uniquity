package p_uniquity_video

import (
	"net/http"

	"github.com/UniquityVentures/lago/getters"
	"github.com/UniquityVentures/lago/lago"
	"github.com/UniquityVentures/lago/plugins/p_users"
	"github.com/UniquityVentures/lago/registry"
	"github.com/UniquityVentures/lago/views"
	"gorm.io/gorm"
)

type rawPreload struct{}

func (rawPreload) Patch(_ views.View, _ *http.Request, query gorm.ChainInterface[RawFootage]) gorm.ChainInterface[RawFootage] {
	return query.Preload("Files", nil).Preload("AssignedTo", nil).Preload("AssignedTo.User", nil)
}

type editedPreload struct{}

func (editedPreload) Patch(_ views.View, _ *http.Request, query gorm.ChainInterface[EditedVideo]) gorm.ChainInterface[EditedVideo] {
	return query.Preload("RawFootage", nil).Preload("EditedVNode", nil)
}

type publishedPreload struct{}

func (publishedPreload) Patch(_ views.View, _ *http.Request, query gorm.ChainInterface[PublishedVideo]) gorm.ChainInterface[PublishedVideo] {
	return query.Preload("EditedVideo", nil).Preload("EditedVideo.RawFootage", nil)
}

func init() {
	auth := p_users.AuthenticationLayer{}

	lago.RegistryView.Register("video.HubView",
		lago.GetPageView("video.HubPage").
			WithLayer("users.auth", auth))

	// --- Raw footage ---
	lago.RegistryView.Register("video.RawListView",
		lago.GetPageView("video.RawFootageTable").
			WithLayer("users.auth", auth).
			WithLayer("video.raw_list", views.LayerList[RawFootage]{
				Key: getters.Static("rawFootages"),
				QueryPatchers: views.QueryPatchers[RawFootage]{
					registry.Pair[string, views.QueryPatcher[RawFootage]]{Key: "video.raw_preload", Value: rawPreload{}},
				},
			}))

	lago.RegistryView.Register("video.RawDetailView",
		lago.GetPageView("video.RawFootageDetail").
			WithLayer("users.auth", auth).
			WithLayer("video.raw_detail", views.LayerDetail[RawFootage]{
				Key:          getters.Static("rawFootage"),
				PathParamKey: getters.Static("id"),
				QueryPatchers: views.QueryPatchers[RawFootage]{
					registry.Pair[string, views.QueryPatcher[RawFootage]]{Key: "video.raw_preload", Value: rawPreload{}},
				},
			}))

	lago.RegistryView.Register("video.RawCreateView",
		lago.GetPageView("video.RawFootageCreateForm").
			WithLayer("users.auth", auth).
			WithLayer("video.raw_create", views.LayerCreate[RawFootage]{
				SuccessURL: lago.RoutePath("video.RawDetailRoute", map[string]getters.Getter[any]{
					"id": getters.Any(getters.Key[uint]("$id")),
				}),
			}))

	lago.RegistryView.Register("video.RawUpdateView",
		lago.GetPageView("video.RawFootageUpdateForm").
			WithLayer("users.auth", auth).
			WithLayer("video.raw_detail", views.LayerDetail[RawFootage]{
				Key:          getters.Static("rawFootage"),
				PathParamKey: getters.Static("id"),
				QueryPatchers: views.QueryPatchers[RawFootage]{
					registry.Pair[string, views.QueryPatcher[RawFootage]]{Key: "video.raw_preload", Value: rawPreload{}},
				},
			}).
			WithLayer("video.raw_update", views.LayerUpdate[RawFootage]{
				Key:        getters.Static("rawFootage"),
				SuccessURL: lago.RoutePath("video.RawListRoute", nil),
			}))

	lago.RegistryView.Register("video.RawDeleteView",
		lago.GetPageView("video.RawFootageDeleteForm").
			WithLayer("users.auth", auth).
			WithLayer("video.raw_detail", views.LayerDetail[RawFootage]{
				Key:          getters.Static("rawFootage"),
				PathParamKey: getters.Static("id"),
				QueryPatchers: views.QueryPatchers[RawFootage]{
					registry.Pair[string, views.QueryPatcher[RawFootage]]{Key: "video.raw_preload", Value: rawPreload{}},
				},
			}).
			WithLayer("video.raw_delete", views.LayerDelete[RawFootage]{
				Key:        getters.Static("rawFootage"),
				SuccessURL: lago.RoutePath("video.RawListRoute", nil),
			}))

	lago.RegistryView.Register("video.RawSelectView",
		lago.GetPageView("video.RawFootageSelectionTable").
			WithLayer("users.auth", auth).
			WithLayer("video.raw_select", views.LayerList[RawFootage]{
				Key: getters.Static("rawFootages"),
				QueryPatchers: views.QueryPatchers[RawFootage]{
					registry.Pair[string, views.QueryPatcher[RawFootage]]{Key: "video.raw_preload", Value: rawPreload{}},
				},
			}))

	// --- Edited videos ---
	lago.RegistryView.Register("video.EditedListView",
		lago.GetPageView("video.EditedVideoTable").
			WithLayer("users.auth", auth).
			WithLayer("video.edited_list", views.LayerList[EditedVideo]{
				Key: getters.Static("editedVideos"),
				QueryPatchers: views.QueryPatchers[EditedVideo]{
					registry.Pair[string, views.QueryPatcher[EditedVideo]]{Key: "video.edited_preload", Value: editedPreload{}},
				},
			}))

	lago.RegistryView.Register("video.EditedDetailView",
		lago.GetPageView("video.EditedVideoDetail").
			WithLayer("users.auth", auth).
			WithLayer("video.edited_detail", views.LayerDetail[EditedVideo]{
				Key:          getters.Static("editedVideo"),
				PathParamKey: getters.Static("id"),
				QueryPatchers: views.QueryPatchers[EditedVideo]{
					registry.Pair[string, views.QueryPatcher[EditedVideo]]{Key: "video.edited_preload", Value: editedPreload{}},
				},
			}))

	lago.RegistryView.Register("video.EditedCreateView",
		lago.GetPageView("video.EditedVideoCreateForm").
			WithLayer("users.auth", auth).
			WithLayer("video.edited_create", views.LayerCreate[EditedVideo]{
				SuccessURL: lago.RoutePath("video.EditedDetailRoute", map[string]getters.Getter[any]{
					"id": getters.Any(getters.Key[uint]("$id")),
				}),
			}))

	lago.RegistryView.Register("video.EditedUpdateView",
		lago.GetPageView("video.EditedVideoUpdateForm").
			WithLayer("users.auth", auth).
			WithLayer("video.edited_detail", views.LayerDetail[EditedVideo]{
				Key:          getters.Static("editedVideo"),
				PathParamKey: getters.Static("id"),
				QueryPatchers: views.QueryPatchers[EditedVideo]{
					registry.Pair[string, views.QueryPatcher[EditedVideo]]{Key: "video.edited_preload", Value: editedPreload{}},
				},
			}).
			WithLayer("video.edited_update", views.LayerUpdate[EditedVideo]{
				Key:        getters.Static("editedVideo"),
				SuccessURL: lago.RoutePath("video.EditedListRoute", nil),
			}))

	lago.RegistryView.Register("video.EditedDeleteView",
		lago.GetPageView("video.EditedVideoDeleteForm").
			WithLayer("users.auth", auth).
			WithLayer("video.edited_detail", views.LayerDetail[EditedVideo]{
				Key:          getters.Static("editedVideo"),
				PathParamKey: getters.Static("id"),
				QueryPatchers: views.QueryPatchers[EditedVideo]{
					registry.Pair[string, views.QueryPatcher[EditedVideo]]{Key: "video.edited_preload", Value: editedPreload{}},
				},
			}).
			WithLayer("video.edited_delete", views.LayerDelete[EditedVideo]{
				Key:        getters.Static("editedVideo"),
				SuccessURL: lago.RoutePath("video.EditedListRoute", nil),
			}))

	lago.RegistryView.Register("video.EditedSelectView",
		lago.GetPageView("video.EditedVideoSelectionTable").
			WithLayer("users.auth", auth).
			WithLayer("video.edited_select", views.LayerList[EditedVideo]{
				Key: getters.Static("editedVideos"),
				QueryPatchers: views.QueryPatchers[EditedVideo]{
					registry.Pair[string, views.QueryPatcher[EditedVideo]]{Key: "video.edited_preload", Value: editedPreload{}},
				},
			}))

	// --- Published videos ---
	lago.RegistryView.Register("video.PublishedListView",
		lago.GetPageView("video.PublishedVideoTable").
			WithLayer("users.auth", auth).
			WithLayer("video.published_list", views.LayerList[PublishedVideo]{
				Key: getters.Static("publishedVideos"),
				QueryPatchers: views.QueryPatchers[PublishedVideo]{
					registry.Pair[string, views.QueryPatcher[PublishedVideo]]{Key: "video.published_preload", Value: publishedPreload{}},
				},
			}))

	lago.RegistryView.Register("video.PublishedDetailView",
		lago.GetPageView("video.PublishedVideoDetail").
			WithLayer("users.auth", auth).
			WithLayer("video.published_detail", views.LayerDetail[PublishedVideo]{
				Key:          getters.Static("publishedVideo"),
				PathParamKey: getters.Static("id"),
				QueryPatchers: views.QueryPatchers[PublishedVideo]{
					registry.Pair[string, views.QueryPatcher[PublishedVideo]]{Key: "video.published_preload", Value: publishedPreload{}},
				},
			}))

	lago.RegistryView.Register("video.PublishedCreateView",
		lago.GetPageView("video.PublishedVideoCreateForm").
			WithLayer("users.auth", auth).
			WithLayer("video.published_create", views.LayerCreate[PublishedVideo]{
				SuccessURL: lago.RoutePath("video.PublishedDetailRoute", map[string]getters.Getter[any]{
					"id": getters.Any(getters.Key[uint]("$id")),
				}),
			}))

	lago.RegistryView.Register("video.PublishedUpdateView",
		lago.GetPageView("video.PublishedVideoUpdateForm").
			WithLayer("users.auth", auth).
			WithLayer("video.published_detail", views.LayerDetail[PublishedVideo]{
				Key:          getters.Static("publishedVideo"),
				PathParamKey: getters.Static("id"),
				QueryPatchers: views.QueryPatchers[PublishedVideo]{
					registry.Pair[string, views.QueryPatcher[PublishedVideo]]{Key: "video.published_preload", Value: publishedPreload{}},
				},
			}).
			WithLayer("video.published_update", views.LayerUpdate[PublishedVideo]{
				Key:        getters.Static("publishedVideo"),
				SuccessURL: lago.RoutePath("video.PublishedListRoute", nil),
			}))

	lago.RegistryView.Register("video.PublishedDeleteView",
		lago.GetPageView("video.PublishedVideoDeleteForm").
			WithLayer("users.auth", auth).
			WithLayer("video.published_detail", views.LayerDetail[PublishedVideo]{
				Key:          getters.Static("publishedVideo"),
				PathParamKey: getters.Static("id"),
				QueryPatchers: views.QueryPatchers[PublishedVideo]{
					registry.Pair[string, views.QueryPatcher[PublishedVideo]]{Key: "video.published_preload", Value: publishedPreload{}},
				},
			}).
			WithLayer("video.published_delete", views.LayerDelete[PublishedVideo]{
				Key:        getters.Static("publishedVideo"),
				SuccessURL: lago.RoutePath("video.PublishedListRoute", nil),
			}))

	lago.RegistryView.Register("video.PublishedSelectView",
		lago.GetPageView("video.PublishedVideoSelectionTable").
			WithLayer("users.auth", auth).
			WithLayer("video.published_select", views.LayerList[PublishedVideo]{
				Key: getters.Static("publishedVideos"),
				QueryPatchers: views.QueryPatchers[PublishedVideo]{
					registry.Pair[string, views.QueryPatcher[PublishedVideo]]{Key: "video.published_preload", Value: publishedPreload{}},
				},
			}))
}
