package p_uniquity_finance_accounts

import (
	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/registry"
)

func pluginRoutes() lago.PluginFeatures[lago.Route] {
	return lago.PluginFeatures[lago.Route]{
		Entries: []registry.Pair[string, lago.Route]{
			{Key: "finance_accounts.DefaultRoute", Value: lago.Route{Path: AppUrl, Handler: lago.NewDynamicView("finance_accounts.AccountListView")}},
			// Account CRUD lives under /accounts/… so /{id}/ patterns do not capture "journals", "currencies", etc.
			{Key: "finance_accounts.AccountCreateRoute", Value: lago.Route{Path: AppUrl + "accounts/create/", Handler: lago.NewDynamicView("finance_accounts.AccountCreateView")}},
			{Key: "finance_accounts.AccountSelectRoute", Value: lago.Route{Path: AppUrl + "accounts/select/", Handler: lago.NewDynamicView("finance_accounts.AccountSelectView")}},
			{Key: "finance_accounts.AccountDetailRoute", Value: lago.Route{Path: AppUrl + "accounts/{id}/", Handler: lago.NewDynamicView("finance_accounts.AccountDetailView")}},
			{Key: "finance_accounts.AccountUpdateRoute", Value: lago.Route{Path: AppUrl + "accounts/{id}/edit/", Handler: lago.NewDynamicView("finance_accounts.AccountUpdateView")}},
			{Key: "finance_accounts.AccountDeleteRoute", Value: lago.Route{Path: AppUrl + "accounts/{id}/delete/", Handler: lago.NewDynamicView("finance_accounts.AccountDeleteView")}},

			{Key: "finance_accounts.CurrencyListRoute", Value: lago.Route{Path: AppUrl + "currencies/", Handler: lago.NewDynamicView("finance_accounts.CurrencyListView")}},
			{Key: "finance_accounts.CurrencyCreateRoute", Value: lago.Route{Path: AppUrl + "currencies/create/", Handler: lago.NewDynamicView("finance_accounts.CurrencyCreateView")}},
			{Key: "finance_accounts.CurrencySelectRoute", Value: lago.Route{Path: AppUrl + "currencies/select/", Handler: lago.NewDynamicView("finance_accounts.CurrencySelectView")}},
			{Key: "finance_accounts.CurrencyDetailRoute", Value: lago.Route{Path: AppUrl + "currencies/{id}/", Handler: lago.NewDynamicView("finance_accounts.CurrencyDetailView")}},
			{Key: "finance_accounts.CurrencyUpdateRoute", Value: lago.Route{Path: AppUrl + "currencies/{id}/edit/", Handler: lago.NewDynamicView("finance_accounts.CurrencyUpdateView")}},
			{Key: "finance_accounts.CurrencyDeleteRoute", Value: lago.Route{Path: AppUrl + "currencies/{id}/delete/", Handler: lago.NewDynamicView("finance_accounts.CurrencyDeleteView")}},

			{Key: "finance_accounts.JournalListRoute", Value: lago.Route{Path: AppUrl + "journals/", Handler: lago.NewDynamicView("finance_accounts.JournalListView")}},
			{Key: "finance_accounts.JournalSelectRoute", Value: lago.Route{Path: AppUrl + "journals/select/", Handler: lago.NewDynamicView("finance_accounts.JournalSelectView")}},
			{Key: "finance_accounts.JournalCreateRoute", Value: lago.Route{Path: AppUrl + "journals/create/", Handler: lago.NewDynamicView("finance_accounts.JournalCreateView")}},
			{Key: "finance_accounts.JournalEntryCreateRoute", Value: lago.Route{Path: AppUrl + "journals/{journal_id}/entries/create/", Handler: lago.NewDynamicView("finance_accounts.JournalEntryCreateView")}},
			{Key: "finance_accounts.SourceDocSelectRoute", Value: lago.Route{Path: AppUrl + "source-docs/select/", Handler: lago.NewDynamicView("finance_accounts.SourceDocSelectView")}},
			{Key: "finance_accounts.JournalDetailRoute", Value: lago.Route{Path: AppUrl + "journals/{id}/", Handler: lago.NewDynamicView("finance_accounts.JournalDetailView")}},
			{Key: "finance_accounts.JournalEntryDetailRoute", Value: lago.Route{Path: AppUrl + "journal-entries/{id}/", Handler: lago.NewDynamicView("finance_accounts.JournalEntryDetailView")}},
			{Key: "finance_accounts.JournalUpdateRoute", Value: lago.Route{Path: AppUrl + "journals/{id}/edit/", Handler: lago.NewDynamicView("finance_accounts.JournalUpdateView")}},
			{Key: "finance_accounts.JournalDeleteRoute", Value: lago.Route{Path: AppUrl + "journals/{id}/delete/", Handler: lago.NewDynamicView("finance_accounts.JournalDeleteView")}},
			{Key: "finance_accounts.AccountingPreferencesRoute", Value: lago.Route{Path: AppUrl + "preferences/", Handler: lago.NewDynamicView("finance_accounts.AccountingPreferencesView")}},
		},
	}
}
