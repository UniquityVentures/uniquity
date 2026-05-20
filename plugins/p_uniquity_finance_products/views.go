package p_uniquity_finance_products

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
		user := p_users.UserFromContext(r.Context(), "finance_products.SuperuserOnlyLayer")
		if !user.IsSuperuser {
			slog.Error("finance_products.SuperuserOnlyLayer: forbidden", "user_id", user.ID)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func pluginViews() lamu.PluginFeatures[*views.View] {
	qp := views.QueryPatchers[Product]{
		registry.Pair[string, views.QueryPatcher[Product]]{Key: "finance_products.preload_taxes", Value: productPreloadTaxes},
	}
	return lamu.PluginFeatures[*views.View]{
		Entries: []registry.Pair[string, *views.View]{
			{
				Key: "finance_products.ProductListView",
				Value: lamu.GetPageView("finance_products.ProductTable").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_products.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_products.product_list", views.LayerList[Product]{
						Key:           getters.Static("products"),
						QueryPatchers: qp,
					}),
			},
			{
				Key: "finance_products.ProductDetailView",
				Value: lamu.GetPageView("finance_products.ProductDetail").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_products.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_products.product_detail", views.LayerDetail[Product]{
						Key:           getters.Static("product"),
						PathParamKey:  getters.Static("id"),
						QueryPatchers: qp,
					}),
			},
			{
				Key: "finance_products.ProductCreateView",
				Value: lamu.GetPageView("finance_products.ProductCreateForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_products.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_products.product_create", views.LayerCreate[Product]{
						SuccessURL: lamu.RoutePath("finance_products.ProductDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("$id")),
						}),
					}),
			},
			{
				Key: "finance_products.ProductUpdateView",
				Value: lamu.GetPageView("finance_products.ProductUpdateForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_products.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_products.product_detail", views.LayerDetail[Product]{
						Key:           getters.Static("product"),
						PathParamKey:  getters.Static("id"),
						QueryPatchers: qp,
					}).
					WithLayer("finance_products.product_update", views.LayerUpdate[Product]{
						Key:        getters.Static("product"),
						SuccessURL: lamu.RoutePath("finance_products.DefaultRoute", nil),
					}),
			},
			{
				Key: "finance_products.ProductDeleteView",
				Value: lamu.GetPageView("finance_products.ProductDeleteForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_products.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_products.product_detail", views.LayerDetail[Product]{
						Key:           getters.Static("product"),
						PathParamKey:  getters.Static("id"),
						QueryPatchers: qp,
					}).
					WithLayer("finance_products.product_delete", views.LayerDelete[Product]{
						Key:        getters.Static("product"),
						SuccessURL: lamu.RoutePath("finance_products.DefaultRoute", nil),
					}),
			},
			{
				Key: "finance_products.ProductFkSelectView",
				Value: lamu.GetPageView("finance_products.ProductFkSelectionTable").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_products.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_products.product_fk_list", views.LayerList[Product]{
						Key: getters.Static("products"),
					}),
			},
		},
		Patches: []registry.Pair[string, func(*views.View) *views.View]{
			{Key: "finance_accounts.AccountingPreferencesView", Value: patchAccountingPreferencesView},
		},
	}
}
