package p_uniquity_tax_rates

import (
	"log/slog"
	"net/http"

	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/plugins/p_users"
	"github.com/UniquityVentures/lamu/registry"
	"github.com/UniquityVentures/lamu/views"
)

type SuperuserOnlyLayer struct{}

func (SuperuserOnlyLayer) Next(_ views.View, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := p_users.UserFromContext(r.Context(), "tax_rates.SuperuserOnlyLayer")
		if !user.IsSuperuser {
			slog.Error("tax_rates.SuperuserOnlyLayer: forbidden", "user_id", user.ID)
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
				Key:   "tax_rates.HubView",
				Value: lamu.GetPageView("tax_rates.HubPage").WithLayer("p_users.auth", auth).WithLayer("tax_rates.superuser", su),
			},
			{
				Key: "tax_rates.TaxRateSelectView",
				Value: lamu.GetPageView("tax_rates.TaxRateSelectionTable").
					WithLayer("p_users.auth", auth).
					WithLayer("tax_rates.superuser", su).
					WithLayer("tax_rates.tax_rate_select_list", views.LayerList[TaxRate]{
						Key: getters.Static("taxRates"),
						QueryPatchers: views.QueryPatchers[TaxRate]{
							{Key: "tax_rates.tax_rate_select_preload", Value: views.QueryPatcherPreload[TaxRate]{Fields: []string{"Entity"}}},
						},
					}),
			},
		},
	}
}
