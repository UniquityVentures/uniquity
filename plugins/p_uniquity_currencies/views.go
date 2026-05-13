package p_uniquity_currencies

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
		user := p_users.UserFromContext(r.Context(), "currencies.SuperuserOnlyLayer")
		if !user.IsSuperuser {
			slog.Error("currencies.SuperuserOnlyLayer: forbidden", "user_id", user.ID)
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
				Key: "currencies.CurrencyListView",
				Value: lamu.GetPageView("currencies.CurrencyTable").
					WithLayer("p_users.auth", auth).
					WithLayer("currencies.superuser", su).
					WithLayer("currencies.currency_list", views.LayerList[Currency]{Key: getters.Static("currencies")}),
			},
			{
				Key: "currencies.CurrencySelectView",
				Value: lamu.GetPageView("currencies.CurrencySelectionTable").
					WithLayer("p_users.auth", auth).
					WithLayer("currencies.superuser", su).
					WithLayer("currencies.currency_select_list", views.LayerList[Currency]{Key: getters.Static("currencies")}),
			},
			{
				Key: "currencies.CurrencyDetailView",
				Value: lamu.GetPageView("currencies.CurrencyDetail").
					WithLayer("p_users.auth", auth).
					WithLayer("currencies.superuser", su).
					WithLayer("currencies.currency_detail", views.LayerDetail[Currency]{
						Key:          getters.Static("currency"),
						PathParamKey: getters.Static("id"),
					}),
			},
			{
				Key: "currencies.CurrencyCreateView",
				Value: lamu.GetPageView("currencies.CurrencyCreateForm").
					WithLayer("p_users.auth", auth).
					WithLayer("currencies.superuser", su).
					WithLayer("currencies.currency_create", views.LayerCreate[Currency]{
						SuccessURL: lamu.RoutePath("currencies.CurrencyDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("$id")),
						}),
					}),
			},
			{
				Key: "currencies.CurrencyUpdateView",
				Value: lamu.GetPageView("currencies.CurrencyUpdateForm").
					WithLayer("p_users.auth", auth).
					WithLayer("currencies.superuser", su).
					WithLayer("currencies.currency_update_detail", views.LayerDetail[Currency]{
						Key:          getters.Static("currency"),
						PathParamKey: getters.Static("id"),
					}).
					WithLayer("currencies.currency_update", views.LayerUpdate[Currency]{
						Key: getters.Static("currency"),
						SuccessURL: lamu.RoutePath("currencies.CurrencyDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("currency.ID")),
						}),
					}),
			},
			{
				Key: "currencies.CurrencyDeleteView",
				Value: lamu.GetPageView("currencies.CurrencyDeleteForm").
					WithLayer("p_users.auth", auth).
					WithLayer("currencies.superuser", su).
					WithLayer("currencies.currency_delete_detail", views.LayerDetail[Currency]{
						Key:          getters.Static("currency"),
						PathParamKey: getters.Static("id"),
					}).
					WithLayer("currencies.currency_delete", views.LayerDelete[Currency]{
						Key:        getters.Static("currency"),
						SuccessURL: lamu.RoutePath("currencies.CurrencyListRoute", nil),
					}),
			},
		},
	}
}
