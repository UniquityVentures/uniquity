package p_uniquity_finance_fiscal_year

import (
	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/registry"
)

func pluginRoutes() lago.PluginFeatures[lago.Route] {
	fy := AppUrl + "fy/"
	return lago.PluginFeatures[lago.Route]{
		Entries: []registry.Pair[string, lago.Route]{
			{Key: "finance_fiscal_years.DefaultRoute", Value: lago.Route{Path: AppUrl, Handler: lago.NewDynamicView("finance_fiscal_years.FiscalYearListView")}},
			{Key: "finance_fiscal_years.FiscalYearCreateRoute", Value: lago.Route{Path: AppUrl + "create/", Handler: lago.NewDynamicView("finance_fiscal_years.FiscalYearCreateView")}},
			{Key: "finance_fiscal_years.FiscalYearDetailRoute", Value: lago.Route{Path: fy + "{id}/", Handler: lago.NewDynamicView("finance_fiscal_years.FiscalYearDetailView")}},
			{Key: "finance_fiscal_years.FiscalYearUpdateRoute", Value: lago.Route{Path: fy + "{id}/edit/", Handler: lago.NewDynamicView("finance_fiscal_years.FiscalYearUpdateView")}},
			{Key: "finance_fiscal_years.FiscalYearDeleteRoute", Value: lago.Route{Path: fy + "{id}/delete/", Handler: lago.NewDynamicView("finance_fiscal_years.FiscalYearDeleteView")}},
			{Key: "finance_fiscal_years.FiscalYearSelectRoute", Value: lago.Route{Path: AppUrl + "select/", Handler: lago.NewDynamicView("finance_fiscal_years.FiscalYearSelectView")}},
		},
	}
}
