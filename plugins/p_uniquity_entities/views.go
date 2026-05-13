package p_uniquity_entities

import (
	"log/slog"
	"net/http"

	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/plugins/p_users"
	"github.com/UniquityVentures/lamu/registry"
	"github.com/UniquityVentures/lamu/views"
)

// SuperuserOnlyLayer returns 401 unless the authenticated user is a superuser.
type SuperuserOnlyLayer struct{}

func (SuperuserOnlyLayer) Next(_ views.View, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := p_users.UserFromContext(r.Context(), "entities.SuperuserOnlyLayer")
		if !user.IsSuperuser {
			slog.Error("entities.SuperuserOnlyLayer: forbidden", "user_id", user.ID)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func pluginViews() lamu.PluginFeatures[*views.View] {
	auth := p_users.AuthenticationLayer{}
	su := SuperuserOnlyLayer{}
	return lamu.PluginFeatures[*views.View]{
		Entries: []registry.Pair[string, *views.View]{
			{
				Key: "entities.EntityListView",
				Value: lamu.GetPageView("entities.EntityTable").
					WithLayer("p_users.auth", auth).
					WithLayer("entities.superuser", su).
					WithLayer("entities.entity_list", views.LayerList[Entity]{
						Key: getters.Static("entities"),
						QueryPatchers: views.QueryPatchers[Entity]{
							{Key: "entities.preload_currency", Value: views.QueryPatcherPreload[Entity]{Fields: []string{"Currency"}}},
						},
					}),
			},
			{
				Key: "entities.EntitySelectView",
				Value: lamu.GetPageView("entities.EntitySelectionTable").
					WithLayer("p_users.auth", auth).
					WithLayer("entities.superuser", su).
					WithLayer("entities.entity_select_list", views.LayerList[Entity]{
						Key: getters.Static("entities"),
						QueryPatchers: views.QueryPatchers[Entity]{
							{Key: "entities.preload_currency", Value: views.QueryPatcherPreload[Entity]{Fields: []string{"Currency"}}},
						},
					}),
			},
			{
				Key: "entities.EntityDetailView",
				Value: lamu.GetPageView("entities.EntityDetail").
					WithLayer("p_users.auth", auth).
					WithLayer("entities.superuser", su).
					WithLayer("entities.entity_detail", views.LayerDetail[Entity]{
						Key:          getters.Static("entity"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[Entity]{
							{Key: "entities.preload_currency", Value: views.QueryPatcherPreload[Entity]{Fields: []string{"Currency"}}},
						},
					}),
			},
			{
				Key: "entities.EntityCreateView",
				Value: lamu.GetPageView("entities.EntityCreateForm").
					WithLayer("p_users.auth", auth).
					WithLayer("entities.superuser", su).
					WithLayer("entities.entity_create", views.LayerCreate[Entity]{
						SuccessURL: lamu.RoutePath("entities.EntityDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("$id")),
						}),
					}),
			},
			{
				Key: "entities.EntityUpdateView",
				Value: lamu.GetPageView("entities.EntityUpdateForm").
					WithLayer("p_users.auth", auth).
					WithLayer("entities.superuser", su).
					WithLayer("entities.entity_update_detail", views.LayerDetail[Entity]{
						Key:          getters.Static("entity"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[Entity]{
							{Key: "entities.preload_currency", Value: views.QueryPatcherPreload[Entity]{Fields: []string{"Currency"}}},
						},
					}).
					WithLayer("entities.entity_update", views.LayerUpdate[Entity]{
						Key: getters.Static("entity"),
						SuccessURL: lamu.RoutePath("entities.EntityDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("entity.ID")),
						}),
					}),
			},
			{
				Key: "entities.EntityDeleteView",
				Value: lamu.GetPageView("entities.EntityDeleteForm").
					WithLayer("p_users.auth", auth).
					WithLayer("entities.superuser", su).
					WithLayer("entities.entity_delete_detail", views.LayerDetail[Entity]{
						Key:          getters.Static("entity"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[Entity]{
							{Key: "entities.preload_currency", Value: views.QueryPatcherPreload[Entity]{Fields: []string{"Currency"}}},
						},
					}).
					WithLayer("entities.entity_delete", views.LayerDelete[Entity]{
						Key:        getters.Static("entity"),
						SuccessURL: lamu.RoutePath("entities.EntityListRoute", nil),
					}),
			},
		},
	}
}
