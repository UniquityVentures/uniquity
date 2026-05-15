package p_uniquity_finance_customer

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
		user := p_users.UserFromContext(r.Context(), "finance_customers.SuperuserOnlyLayer")
		if !user.IsSuperuser {
			slog.Error("finance_customers.SuperuserOnlyLayer: forbidden", "user_id", user.ID)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func pluginViews() lamu.PluginFeatures[*views.View] {
	return lamu.PluginFeatures[*views.View]{
		Entries: []registry.Pair[string, *views.View]{
			{
				Key: "finance_customers.CustomerListView",
				Value: lamu.GetPageView("finance_customers.CustomerTable").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_customers.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_customers.customer_list", views.LayerList[Customer]{
						Key: getters.Static("customers"),
					}),
			},
			{
				Key: "finance_customers.CustomerDetailView",
				Value: lamu.GetPageView("finance_customers.CustomerDetail").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_customers.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_customers.customer_detail", views.LayerDetail[Customer]{
						Key:          getters.Static("customer"),
						PathParamKey: getters.Static("id"),
					}),
			},
			{
				Key: "finance_customers.CustomerCreateView",
				Value: lamu.GetPageView("finance_customers.CustomerCreateForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_customers.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_customers.customer_create", views.LayerCreate[Customer]{
						SuccessURL: lamu.RoutePath("finance_customers.CustomerDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("$id")),
						}),
					}),
			},
			{
				Key: "finance_customers.CustomerUpdateView",
				Value: lamu.GetPageView("finance_customers.CustomerUpdateForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_customers.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_customers.customer_detail", views.LayerDetail[Customer]{
						Key:          getters.Static("customer"),
						PathParamKey: getters.Static("id"),
					}).
					WithLayer("finance_customers.customer_update", views.LayerUpdate[Customer]{
						Key:        getters.Static("customer"),
						SuccessURL: lamu.RoutePath("finance_customers.DefaultRoute", nil),
					}),
			},
			{
				Key: "finance_customers.CustomerDeleteView",
				Value: lamu.GetPageView("finance_customers.CustomerDeleteForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_customers.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_customers.customer_detail", views.LayerDetail[Customer]{
						Key:          getters.Static("customer"),
						PathParamKey: getters.Static("id"),
					}).
					WithLayer("finance_customers.customer_delete", views.LayerDelete[Customer]{
						Key:        getters.Static("customer"),
						SuccessURL: lamu.RoutePath("finance_customers.DefaultRoute", nil),
					}),
			},
			{
				Key: "finance_customers.CustomerFkSelectView",
				Value: lamu.GetPageView("finance_customers.CustomerFkSelectionTable").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_customers.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_customers.customer_fk_list", views.LayerList[Customer]{
						Key: getters.Static("customers"),
					}),
			},
		},
	}
}
