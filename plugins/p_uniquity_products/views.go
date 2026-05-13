package p_uniquity_products

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
		user := p_users.UserFromContext(r.Context(), "products.SuperuserOnlyLayer")
		if !user.IsSuperuser {
			slog.Error("products.SuperuserOnlyLayer: forbidden", "user_id", user.ID)
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
				Key:   "products.HubView",
				Value: lamu.GetPageView("products.HubPage").WithLayer("p_users.auth", auth).WithLayer("products.superuser", su),
			},
			{
				Key: "products.ProductSelectView",
				Value: lamu.GetPageView("products.ProductSelectionTable").
					WithLayer("p_users.auth", auth).
					WithLayer("products.superuser", su).
					WithLayer("products.product_select_list", views.LayerList[Product]{
						Key: getters.Static("products"),
						QueryPatchers: views.QueryPatchers[Product]{
							{Key: "products.product_select_preload", Value: views.QueryPatcherPreload[Product]{Fields: []string{"Entity"}}},
						},
					}),
			},
		},
	}
}
