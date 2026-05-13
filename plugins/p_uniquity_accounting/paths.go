package p_uniquity_accounting

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pluginRoutes() lamu.PluginFeatures[lamu.Route] {
	return lamu.PluginFeatures[lamu.Route]{
		Entries: []registry.Pair[string, lamu.Route]{
			{Key: "accounting.TransactionListRoute", Value: lamu.Route{Path: AppUrl, Handler: lamu.NewDynamicView("accounting.TransactionListView")}},
			{Key: "accounting.TransactionCreateRoute", Value: lamu.Route{Path: AppUrl + "transactions/create/", Handler: lamu.NewDynamicView("accounting.TransactionCreateView")}},
			{Key: "accounting.TransactionDetailRoute", Value: lamu.Route{Path: AppUrl + "transactions/t/{id}/", Handler: lamu.NewDynamicView("accounting.TransactionDetailView")}},
			{Key: "accounting.TransactionAccountSelectRoute", Value: lamu.Route{Path: AppUrl + "transactions/select-account/", Handler: lamu.NewDynamicView("accounting.TransactionAccountSelectView")}},
			{Key: "accounting.AccountListRoute", Value: lamu.Route{Path: AppUrl + "accounts/", Handler: lamu.NewDynamicView("accounting.AccountListView")}},
			{Key: "accounting.AccountCreateRoute", Value: lamu.Route{Path: AppUrl + "accounts/create", Handler: lamu.NewDynamicView("accounting.AccountCreateView")}},
			{Key: "accounting.AccountDetailRoute", Value: lamu.Route{Path: AppUrl + "accounts/a/{id}/", Handler: lamu.NewDynamicView("accounting.AccountDetailView")}},
		},
	}
}
