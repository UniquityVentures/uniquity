package p_uniquity_finance_taxes

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
		user := p_users.UserFromContext(r.Context(), "finance_taxes.SuperuserOnlyLayer")
		if !user.IsSuperuser {
			slog.Error("finance_taxes.SuperuserOnlyLayer: forbidden", "user_id", user.ID)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func pluginViews() lamu.PluginFeatures[*views.View] {
	qp := taxQueryPatchers()
	return lamu.PluginFeatures[*views.View]{
		Entries: []registry.Pair[string, *views.View]{
			{
				Key: "finance_taxes.TaxListView",
				Value: lamu.GetPageView("finance_taxes.TaxTable").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_taxes.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_taxes.tax_list", views.LayerList[Tax]{
						Key:           getters.Static("taxes"),
						QueryPatchers: qp,
					}),
			},
			{
				Key: "finance_taxes.TaxDetailView",
				Value: lamu.GetPageView("finance_taxes.TaxDetail").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_taxes.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_taxes.tax_detail", views.LayerDetail[Tax]{
						Key:           getters.Static("tax"),
						PathParamKey:  getters.Static("id"),
						QueryPatchers: qp,
					}),
			},
			{
				Key: "finance_taxes.TaxCreateView",
				Value: lamu.GetPageView("finance_taxes.TaxCreateForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_taxes.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_taxes.tax_create", views.LayerCreate[Tax]{
						SuccessURL: lamu.RoutePath("finance_taxes.TaxDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("$id")),
						}),
					}),
			},
			{
				Key: "finance_taxes.TaxUpdateView",
				Value: lamu.GetPageView("finance_taxes.TaxUpdateForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_taxes.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_taxes.tax_detail", views.LayerDetail[Tax]{
						Key:           getters.Static("tax"),
						PathParamKey:  getters.Static("id"),
						QueryPatchers: qp,
					}).
					WithLayer("finance_taxes.tax_update", views.LayerUpdate[Tax]{
						Key:        getters.Static("tax"),
						SuccessURL: lamu.RoutePath("finance_taxes.DefaultRoute", nil),
					}),
			},
			{
				Key: "finance_taxes.TaxDeleteView",
				Value: lamu.GetPageView("finance_taxes.TaxDeleteForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_taxes.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_taxes.tax_detail", views.LayerDetail[Tax]{
						Key:           getters.Static("tax"),
						PathParamKey:  getters.Static("id"),
						QueryPatchers: qp,
					}).
					WithLayer("finance_taxes.tax_delete", views.LayerDelete[Tax]{
						Key:        getters.Static("tax"),
						SuccessURL: lamu.RoutePath("finance_taxes.DefaultRoute", nil),
					}),
			},
			{
				Key: "finance_taxes.TaxMultiSelectView",
				Value: lamu.GetPageView("finance_taxes.TaxMultiSelectionTable").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_taxes.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_taxes.tax_multiselect_list", views.LayerList[Tax]{
						Key:           getters.Static("taxes"),
						QueryPatchers: qp,
					}),
			},
		},
	}
}
