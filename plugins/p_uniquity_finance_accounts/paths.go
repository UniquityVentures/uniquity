package p_uniquity_finance_accounts

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pluginRoutes() lamu.PluginFeatures[lamu.Route] {
	return lamu.PluginFeatures[lamu.Route]{
		Entries: []registry.Pair[string, lamu.Route]{
			{Key: "finance_accounts.DefaultRoute", Value: lamu.Route{Path: AppUrl, Handler: lamu.NewDynamicView("finance_accounts.AccountListView")}},
			// Account CRUD lives under /accounts/… so /{id}/ patterns do not capture "journals", "currencies", etc.
			{Key: "finance_accounts.AccountCreateRoute", Value: lamu.Route{Path: AppUrl + "accounts/create/", Handler: lamu.NewDynamicView("finance_accounts.AccountCreateView")}},
			{Key: "finance_accounts.AccountSelectRoute", Value: lamu.Route{Path: AppUrl + "accounts/select/", Handler: lamu.NewDynamicView("finance_accounts.AccountSelectView")}},
			{Key: "finance_accounts.AccountDetailRoute", Value: lamu.Route{Path: AppUrl + "accounts/{id}/", Handler: lamu.NewDynamicView("finance_accounts.AccountDetailView")}},
			{Key: "finance_accounts.AccountUpdateRoute", Value: lamu.Route{Path: AppUrl + "accounts/{id}/edit/", Handler: lamu.NewDynamicView("finance_accounts.AccountUpdateView")}},
			{Key: "finance_accounts.AccountDeleteRoute", Value: lamu.Route{Path: AppUrl + "accounts/{id}/delete/", Handler: lamu.NewDynamicView("finance_accounts.AccountDeleteView")}},

			{Key: "finance_accounts.CurrencyListRoute", Value: lamu.Route{Path: AppUrl + "currencies/", Handler: lamu.NewDynamicView("finance_accounts.CurrencyListView")}},
			{Key: "finance_accounts.CurrencyCreateRoute", Value: lamu.Route{Path: AppUrl + "currencies/create/", Handler: lamu.NewDynamicView("finance_accounts.CurrencyCreateView")}},
			{Key: "finance_accounts.CurrencySelectRoute", Value: lamu.Route{Path: AppUrl + "currencies/select/", Handler: lamu.NewDynamicView("finance_accounts.CurrencySelectView")}},
			{Key: "finance_accounts.CurrencyDetailRoute", Value: lamu.Route{Path: AppUrl + "currencies/{id}/", Handler: lamu.NewDynamicView("finance_accounts.CurrencyDetailView")}},
			{Key: "finance_accounts.CurrencyUpdateRoute", Value: lamu.Route{Path: AppUrl + "currencies/{id}/edit/", Handler: lamu.NewDynamicView("finance_accounts.CurrencyUpdateView")}},
			{Key: "finance_accounts.CurrencyDeleteRoute", Value: lamu.Route{Path: AppUrl + "currencies/{id}/delete/", Handler: lamu.NewDynamicView("finance_accounts.CurrencyDeleteView")}},

			{Key: "finance_accounts.JournalListRoute", Value: lamu.Route{Path: AppUrl + "journals/", Handler: lamu.NewDynamicView("finance_accounts.JournalListView")}},
			{Key: "finance_accounts.JournalCreateRoute", Value: lamu.Route{Path: AppUrl + "journals/create/", Handler: lamu.NewDynamicView("finance_accounts.JournalCreateView")}},
			{Key: "finance_accounts.JournalEntryCreateRoute", Value: lamu.Route{Path: AppUrl + "journals/{journal_id}/entries/create/", Handler: lamu.NewDynamicView("finance_accounts.JournalEntryCreateView")}},
			{Key: "finance_accounts.SourceDocSelectRoute", Value: lamu.Route{Path: AppUrl + "source-docs/select/", Handler: lamu.NewDynamicView("finance_accounts.SourceDocSelectView")}},
			{Key: "finance_accounts.JournalDetailRoute", Value: lamu.Route{Path: AppUrl + "journals/{id}/", Handler: lamu.NewDynamicView("finance_accounts.JournalDetailView")}},
			{Key: "finance_accounts.JournalEntryDetailRoute", Value: lamu.Route{Path: AppUrl + "journal-entries/{id}/", Handler: lamu.NewDynamicView("finance_accounts.JournalEntryDetailView")}},
			{Key: "finance_accounts.JournalUpdateRoute", Value: lamu.Route{Path: AppUrl + "journals/{id}/edit/", Handler: lamu.NewDynamicView("finance_accounts.JournalUpdateView")}},
			{Key: "finance_accounts.JournalDeleteRoute", Value: lamu.Route{Path: AppUrl + "journals/{id}/delete/", Handler: lamu.NewDynamicView("finance_accounts.JournalDeleteView")}},
		},
	}
}
