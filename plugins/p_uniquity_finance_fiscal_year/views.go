package p_uniquity_finance_fiscal_year

import (
	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/getters"
	"github.com/lariv-in/lago/plugins/p_users"
	"github.com/lariv-in/lago/registry"
	"github.com/lariv-in/lago/views"
)

func pluginViews() lago.PluginFeatures[*views.View] {
	return lago.PluginFeatures[*views.View]{
		Entries: []registry.Pair[string, *views.View]{
			{
				Key: "finance_fiscal_years.FiscalYearListView",
				Value: lago.GetPageView("finance_fiscal_years.FiscalYearTable").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_fiscal_years.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("finance_fiscal_years.fiscal_year_list", views.LayerList[FiscalYear]{
						Key: getters.Static("fiscal_years"),
					}),
			},
			{
				Key: "finance_fiscal_years.FiscalYearDetailView",
				Value: lago.GetPageView("finance_fiscal_years.FiscalYearDetail").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_fiscal_years.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("finance_fiscal_years.fiscal_year_detail", views.LayerDetail[FiscalYear]{
						Key:          getters.Static("fiscal_year"),
						PathParamKey: getters.Static("id"),
					}),
			},
			{
				Key: "finance_fiscal_years.FiscalYearCreateView",
				Value: lago.GetPageView("finance_fiscal_years.FiscalYearCreateForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_fiscal_years.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("finance_fiscal_years.fiscal_year_create", views.LayerCreate[FiscalYear]{
						SuccessURL: lago.RoutePath("finance_fiscal_years.FiscalYearDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("$id")),
						}),
					}),
			},
			{
				Key: "finance_fiscal_years.FiscalYearUpdateView",
				Value: lago.GetPageView("finance_fiscal_years.FiscalYearUpdateForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_fiscal_years.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("finance_fiscal_years.fiscal_year_detail", views.LayerDetail[FiscalYear]{
						Key:          getters.Static("fiscal_year"),
						PathParamKey: getters.Static("id"),
					}).
					WithLayer("finance_fiscal_years.fiscal_year_update", views.LayerUpdate[FiscalYear]{
						Key:        getters.Static("fiscal_year"),
						SuccessURL: lago.RoutePath("finance_fiscal_years.DefaultRoute", nil),
					}),
			},
			{
				Key: "finance_fiscal_years.FiscalYearDeleteView",
				Value: lago.GetPageView("finance_fiscal_years.FiscalYearDeleteForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_fiscal_years.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("finance_fiscal_years.fiscal_year_detail", views.LayerDetail[FiscalYear]{
						Key:          getters.Static("fiscal_year"),
						PathParamKey: getters.Static("id"),
					}).
					WithLayer("finance_fiscal_years.fiscal_year_delete", views.LayerDelete[FiscalYear]{
						Key:        getters.Static("fiscal_year"),
						SuccessURL: lago.RoutePath("finance_fiscal_years.DefaultRoute", nil),
					}),
			},
			{
				Key: "finance_fiscal_years.FiscalYearSelectView",
				Value: lago.GetPageView("finance_fiscal_years.FiscalYearSelectionTable").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_fiscal_years.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("finance_fiscal_years.fiscal_year_select", views.LayerList[FiscalYear]{
						Key: getters.Static("fiscal_years"),
					}),
			},
		},
	}
}
