package p_uniquity_video

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/UniquityVentures/lago/getters"
	"github.com/UniquityVentures/lago/lago"
	"github.com/UniquityVentures/lago/plugins/p_users"
	"github.com/UniquityVentures/lago/registry"
	"github.com/UniquityVentures/lago/views"
	uniqempl "github.com/UniquityVentures/uniquity/plugins/p_uniquity_employees"
	"gorm.io/gorm"
)

type rawSelectAssignedFilter struct{}

// Patch limits raw footage choices to rows assigned to the current user's employee,
// unless the user is a superuser (IsSuperuser), in which case all rows are shown.
func (rawSelectAssignedFilter) Patch(_ views.View, r *http.Request, query gorm.ChainInterface[RawFootage]) gorm.ChainInterface[RawFootage] {
	ctx := r.Context()
	user := p_users.UserFromContext(ctx, "video.raw_select_assigned")
	if user.IsSuperuser {
		return query
	}
	db, err := getters.DBFromContext(ctx)
	if err != nil {
		slog.Error("video.raw_select_assigned: db from context", "error", err)
		return query.Where("1 = 0")
	}
	emp, err := gorm.G[uniqempl.Employee](db).Where("user_id = ?", user.ID).Take(ctx)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Error("video.raw_select_assigned: load employee", "error", err, "userID", user.ID)
		}
		return query.Where("1 = 0")
	}
	return query.Where("assigned_to_id = ?", emp.ID)
}

type rawPreload struct{}

func (rawPreload) Patch(_ views.View, _ *http.Request, query gorm.ChainInterface[RawFootage]) gorm.ChainInterface[RawFootage] {
	return query.Preload("Files", nil).Preload("AssignedTo", nil).Preload("AssignedTo.User", nil)
}

type employeeSelectPreload struct{}

func (employeeSelectPreload) Patch(_ views.View, _ *http.Request, query gorm.ChainInterface[uniqempl.Employee]) gorm.ChainInterface[uniqempl.Employee] {
	return query.Preload("User", nil)
}

type editedPreload struct{}

func (editedPreload) Patch(_ views.View, _ *http.Request, query gorm.ChainInterface[EditedVideo]) gorm.ChainInterface[EditedVideo] {
	return query.
		Preload("RawFootage.Files", nil).
		Preload("RawFootage.AssignedTo", nil).
		Preload("RawFootage.AssignedTo.User", nil).
		Preload("EditedVNode", nil)
}

type publishedPreload struct{}

func (publishedPreload) Patch(_ views.View, _ *http.Request, query gorm.ChainInterface[PublishedVideo]) gorm.ChainInterface[PublishedVideo] {
	return query.
		Preload("EditedVideo", nil).
		Preload("EditedVideo.RawFootage", nil).
		Preload("EditedVideo.RawFootage.AssignedTo", nil).
		Preload("EditedVideo.RawFootage.AssignedTo.User", nil)
}

// publishedVideoEditorPointsToEmployeePatcher sets ToEmployeeID from the loaded published video’s
// raw-footage assignee on POST (authoritative; ignores tampering with hidden fields).
type publishedVideoEditorPointsToEmployeePatcher struct{}

func (publishedVideoEditorPointsToEmployeePatcher) Patch(_ views.View, r *http.Request, formData map[string]any, formErrors map[string]error) (map[string]any, map[string]error) {
	pv, err := getters.Key[PublishedVideo]("publishedVideo")(r.Context())
	if err != nil {
		formErrors["_form"] = errors.New("published video not loaded")
		return formData, formErrors
	}
	id := pv.EditedVideo.RawFootage.AssignedToID
	if id == 0 {
		formErrors["_form"] = errors.New("this publication has no responsible editor (raw footage assignee)")
		return formData, formErrors
	}
	formData["ToEmployeeID"] = id
	return formData, formErrors
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
					registry.Pair[string, views.QueryPatcher[RawFootage]]{Key: "video.raw_select_assigned", Value: rawSelectAssignedFilter{}},
					registry.Pair[string, views.QueryPatcher[RawFootage]]{Key: "video.raw_preload", Value: rawPreload{}},
				},
			}))

	lago.RegistryView.Register("video.EmployeeSelectView",
		lago.GetPageView("video.EmployeeSelectionTable").
			WithLayer("users.auth", auth).
			WithLayer("video.employee_select", views.LayerList[uniqempl.Employee]{
				Key: getters.Static("employees"),
				QueryPatchers: views.QueryPatchers[uniqempl.Employee]{
					registry.Pair[string, views.QueryPatcher[uniqempl.Employee]]{Key: "video.employee_select_preload", Value: employeeSelectPreload{}},
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
			}).
			WithLayer("video.published_youtube_meta", youtubePublishedMetaLayer{}))

	lago.RegistryView.Register("video.PublishedEditorPointsCreateView",
		lago.GetPageView("video.PublishedEditorPointsForm").
			WithLayer("users.auth", auth).
			WithLayer("employees.superuser", uniqempl.SuperuserOnlyLayer{}).
			WithLayer("video.published_detail", views.LayerDetail[PublishedVideo]{
				Key:          getters.Static("publishedVideo"),
				PathParamKey: getters.Static("id"),
				QueryPatchers: views.QueryPatchers[PublishedVideo]{
					registry.Pair[string, views.QueryPatcher[PublishedVideo]]{Key: "video.published_preload", Value: publishedPreload{}},
				},
			}).
			WithLayer("video.published_editor_points_create", views.LayerCreate[uniqempl.PointsTransaction]{
				SuccessURL: lago.RoutePath("employees.PointsDetailRoute", map[string]getters.Getter[any]{
					"id": getters.Any(getters.Key[uint]("$id")),
				}),
				FormPatchers: views.FormPatchers{
					registry.Pair[string, views.FormPatcher]{Key: "employees.points_from_user", Value: uniqempl.PointsFormFromUserPatcher{}},
					registry.Pair[string, views.FormPatcher]{Key: "video.published_editor_points_to", Value: publishedVideoEditorPointsToEmployeePatcher{}},
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
