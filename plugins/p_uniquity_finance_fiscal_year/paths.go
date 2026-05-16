package p_uniquity_finance_fiscal_year

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pluginRoutes() lamu.PluginFeatures[lamu.Route] {
	fy := AppUrl + "fy/"
	return lamu.PluginFeatures[lamu.Route]{
		Entries: []registry.Pair[string, lamu.Route]{
			{Key: "finance_fiscal_years.DefaultRoute", Value: lamu.Route{Path: AppUrl, Handler: lamu.NewDynamicView("finance_fiscal_years.FiscalYearListView")}},
			{Key: "finance_fiscal_years.FiscalYearCreateRoute", Value: lamu.Route{Path: AppUrl + "create/", Handler: lamu.NewDynamicView("finance_fiscal_years.FiscalYearCreateView")}},
			{Key: "finance_fiscal_years.FiscalYearDetailRoute", Value: lamu.Route{Path: fy + "{id}/", Handler: lamu.NewDynamicView("finance_fiscal_years.FiscalYearDetailView")}},
			{Key: "finance_fiscal_years.FiscalYearUpdateRoute", Value: lamu.Route{Path: fy + "{id}/edit/", Handler: lamu.NewDynamicView("finance_fiscal_years.FiscalYearUpdateView")}},
			{Key: "finance_fiscal_years.FiscalYearDeleteRoute", Value: lamu.Route{Path: fy + "{id}/delete/", Handler: lamu.NewDynamicView("finance_fiscal_years.FiscalYearDeleteView")}},
			{Key: "finance_fiscal_years.FiscalYearSelectRoute", Value: lamu.Route{Path: AppUrl + "select/", Handler: lamu.NewDynamicView("finance_fiscal_years.FiscalYearSelectView")}},
		},
	}
}
