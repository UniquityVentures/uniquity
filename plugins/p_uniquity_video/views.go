package p_uniquity_video

import (
	"errors"
	"log/slog"
	"net/http"

	uniqempl "github.com/UniquityVentures/uniquity/plugins/p_uniquity_employees"
	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/getters"
	"github.com/lariv-in/lago/plugins/p_users"
	"github.com/lariv-in/lago/registry"
	"github.com/lariv-in/lago/views"
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

func pluginViews() lago.PluginFeatures[*views.View] {
	auth := p_users.AuthenticationLayer{}

	return lago.PluginFeatures[*views.View]{
		Entries: []registry.Pair[string, *views.View]{
			{
				Key: "video.HubView",
				Value: lago.GetPageView("video.HubPage").
					WithLayer("p_users.auth", auth),
			},
			{
				Key: "video.RawListView",
				Value: lago.GetPageView("video.RawFootageTable").
					WithLayer("p_users.auth", auth).
					WithLayer("video.raw_list", views.LayerList[RawFootage]{
						Key: getters.Static("rawFootages"),
						QueryPatchers: views.QueryPatchers[RawFootage]{
							registry.Pair[string, views.QueryPatcher[RawFootage]]{Key: "video.raw_preload", Value: rawPreload{}},
						},
					}),
			},
			{
				Key: "video.RawDetailView",
				Value: lago.GetPageView("video.RawFootageDetail").
					WithLayer("p_users.auth", auth).
					WithLayer("video.raw_detail", views.LayerDetail[RawFootage]{
						Key:          getters.Static("rawFootage"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[RawFootage]{
							registry.Pair[string, views.QueryPatcher[RawFootage]]{Key: "video.raw_preload", Value: rawPreload{}},
						},
					}),
			},
			{
				Key: "video.RawCreateView",
				Value: lago.GetPageView("video.RawFootageCreateForm").
					WithLayer("p_users.auth", auth).
					WithLayer("video.raw_create", views.LayerCreate[RawFootage]{
						SuccessURL: lago.RoutePath("video.RawDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("$id")),
						}),
					}),
			},
			{
				Key: "video.RawUpdateView",
				Value: lago.GetPageView("video.RawFootageUpdateForm").
					WithLayer("p_users.auth", auth).
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
					}),
			},
			{
				Key: "video.RawDeleteView",
				Value: lago.GetPageView("video.RawFootageDeleteForm").
					WithLayer("p_users.auth", auth).
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
					}),
			},
			{
				Key: "video.RawSelectView",
				Value: lago.GetPageView("video.RawFootageSelectionTable").
					WithLayer("p_users.auth", auth).
					WithLayer("video.raw_select", views.LayerList[RawFootage]{
						Key: getters.Static("rawFootages"),
						QueryPatchers: views.QueryPatchers[RawFootage]{
							registry.Pair[string, views.QueryPatcher[RawFootage]]{Key: "video.raw_select_assigned", Value: rawSelectAssignedFilter{}},
							registry.Pair[string, views.QueryPatcher[RawFootage]]{Key: "video.raw_preload", Value: rawPreload{}},
						},
					}),
			},
			{
				Key: "video.EmployeeSelectView",
				Value: lago.GetPageView("video.EmployeeSelectionTable").
					WithLayer("p_users.auth", auth).
					WithLayer("video.employee_select", views.LayerList[uniqempl.Employee]{
						Key: getters.Static("employees"),
						QueryPatchers: views.QueryPatchers[uniqempl.Employee]{
							registry.Pair[string, views.QueryPatcher[uniqempl.Employee]]{Key: "video.employee_select_preload", Value: employeeSelectPreload{}},
						},
					}),
			},
			{
				Key: "video.EditedListView",
				Value: lago.GetPageView("video.EditedVideoTable").
					WithLayer("p_users.auth", auth).
					WithLayer("video.edited_list", views.LayerList[EditedVideo]{
						Key: getters.Static("editedVideos"),
						QueryPatchers: views.QueryPatchers[EditedVideo]{
							registry.Pair[string, views.QueryPatcher[EditedVideo]]{Key: "video.edited_preload", Value: editedPreload{}},
						},
					}),
			},
			{
				Key: "video.EditedDetailView",
				Value: lago.GetPageView("video.EditedVideoDetail").
					WithLayer("p_users.auth", auth).
					WithLayer("video.edited_detail", views.LayerDetail[EditedVideo]{
						Key:          getters.Static("editedVideo"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[EditedVideo]{
							registry.Pair[string, views.QueryPatcher[EditedVideo]]{Key: "video.edited_preload", Value: editedPreload{}},
						},
					}),
			},
			{
				Key: "video.EditedCreateView",
				Value: lago.GetPageView("video.EditedVideoCreateForm").
					WithLayer("p_users.auth", auth).
					WithLayer("video.edited_create", views.LayerCreate[EditedVideo]{
						SuccessURL: lago.RoutePath("video.EditedDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("$id")),
						}),
					}),
			},
			{
				Key: "video.EditedUpdateView",
				Value: lago.GetPageView("video.EditedVideoUpdateForm").
					WithLayer("p_users.auth", auth).
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
					}),
			},
			{
				Key: "video.EditedDeleteView",
				Value: lago.GetPageView("video.EditedVideoDeleteForm").
					WithLayer("p_users.auth", auth).
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
					}),
			},
			{
				Key: "video.EditedSelectView",
				Value: lago.GetPageView("video.EditedVideoSelectionTable").
					WithLayer("p_users.auth", auth).
					WithLayer("video.edited_select", views.LayerList[EditedVideo]{
						Key: getters.Static("editedVideos"),
						QueryPatchers: views.QueryPatchers[EditedVideo]{
							registry.Pair[string, views.QueryPatcher[EditedVideo]]{Key: "video.edited_preload", Value: editedPreload{}},
						},
					}),
			},
			{
				Key: "video.PublishedListView",
				Value: lago.GetPageView("video.PublishedVideoTable").
					WithLayer("p_users.auth", auth).
					WithLayer("video.published_list", views.LayerList[PublishedVideo]{
						Key: getters.Static("publishedVideos"),
						QueryPatchers: views.QueryPatchers[PublishedVideo]{
							registry.Pair[string, views.QueryPatcher[PublishedVideo]]{Key: "video.published_preload", Value: publishedPreload{}},
						},
					}),
			},
			{
				Key: "video.PublishedDetailView",
				Value: lago.GetPageView("video.PublishedVideoDetail").
					WithLayer("p_users.auth", auth).
					WithLayer("video.published_detail", views.LayerDetail[PublishedVideo]{
						Key:          getters.Static("publishedVideo"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[PublishedVideo]{
							registry.Pair[string, views.QueryPatcher[PublishedVideo]]{Key: "video.published_preload", Value: publishedPreload{}},
						},
					}).
					WithLayer("video.published_youtube_meta", youtubePublishedMetaLayer{}),
			},
			{
				Key: "video.PublishedEditorPointsCreateView",
				Value: lago.GetPageView("video.PublishedEditorPointsForm").
					WithLayer("p_users.auth", auth).
					WithLayer("employees.superuser", p_users.SuperuserOnlyLayer{}).
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
					}),
			},
			{
				Key: "video.PublishedCreateView",
				Value: lago.GetPageView("video.PublishedVideoCreateForm").
					WithLayer("p_users.auth", auth).
					WithLayer("video.published_create", views.LayerCreate[PublishedVideo]{
						SuccessURL: lago.RoutePath("video.PublishedDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("$id")),
						}),
					}),
			},
			{
				Key: "video.PublishedUpdateView",
				Value: lago.GetPageView("video.PublishedVideoUpdateForm").
					WithLayer("p_users.auth", auth).
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
					}),
			},
			{
				Key: "video.PublishedDeleteView",
				Value: lago.GetPageView("video.PublishedVideoDeleteForm").
					WithLayer("p_users.auth", auth).
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
					}),
			},
			{
				Key: "video.PublishedSelectView",
				Value: lago.GetPageView("video.PublishedVideoSelectionTable").
					WithLayer("p_users.auth", auth).
					WithLayer("video.published_select", views.LayerList[PublishedVideo]{
						Key: getters.Static("publishedVideos"),
						QueryPatchers: views.QueryPatchers[PublishedVideo]{
							registry.Pair[string, views.QueryPatcher[PublishedVideo]]{Key: "video.published_preload", Value: publishedPreload{}},
						},
					}),
			},
		},
	}
}
