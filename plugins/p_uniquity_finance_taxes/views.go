package p_uniquity_finance_taxes

import (
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/plugins/p_users"
	"github.com/UniquityVentures/lamu/registry"
	"github.com/UniquityVentures/lamu/views"
)


func pluginViews() lamu.PluginFeatures[*views.View] {
	qp := taxQueryPatchers()
	return lamu.PluginFeatures[*views.View]{
		Entries: []registry.Pair[string, *views.View]{
			{
				Key: "finance_taxes.TaxListView",
				Value: lamu.GetPageView("finance_taxes.TaxTable").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_taxes.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("finance_taxes.tax_list", views.LayerList[Tax]{
						Key:           getters.Static("taxes"),
						QueryPatchers: qp,
					}),
			},
			{
				Key: "finance_taxes.TaxDetailView",
				Value: lamu.GetPageView("finance_taxes.TaxDetail").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_taxes.superuser", p_users.SuperuserOnlyLayer{}).
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
					WithLayer("finance_taxes.superuser", p_users.SuperuserOnlyLayer{}).
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
					WithLayer("finance_taxes.superuser", p_users.SuperuserOnlyLayer{}).
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
					WithLayer("finance_taxes.superuser", p_users.SuperuserOnlyLayer{}).
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
					WithLayer("finance_taxes.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("finance_taxes.tax_multiselect_list", views.LayerList[Tax]{
						Key:           getters.Static("taxes"),
						QueryPatchers: qp,
					}),
			},
		},
	}
}
