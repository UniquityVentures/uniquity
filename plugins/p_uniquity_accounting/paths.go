package p_uniquity_accounting

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pluginRoutes() lamu.PluginFeatures[lamu.Route] {
	return lamu.PluginFeatures[lamu.Route]{
		Entries: []registry.Pair[string, lamu.Route]{
			{Key: "accounting.JournalEntryItemListRoute", Value: lamu.Route{Path: AppUrl, Handler: lamu.NewDynamicView("accounting.JournalEntryItemListView")}},
			{Key: "accounting.JournalEntryItemCreateRoute", Value: lamu.Route{Path: AppUrl + "journal-entry-items/create/", Handler: lamu.NewDynamicView("accounting.JournalEntryItemCreateView")}},
			{Key: "accounting.JournalEntryItemDetailRoute", Value: lamu.Route{Path: AppUrl + "journal-entry-items/i/{id}/", Handler: lamu.NewDynamicView("accounting.JournalEntryItemDetailView")}},
			{Key: "accounting.JournalEntryItemAccountSelectRoute", Value: lamu.Route{Path: AppUrl + "journal-entry-items/select-account/", Handler: lamu.NewDynamicView("accounting.JournalEntryItemAccountSelectView")}},
			{Key: "accounting.JournalEntryItemJournalEntrySelectRoute", Value: lamu.Route{Path: AppUrl + "journal-entry-items/select-journal-entry/", Handler: lamu.NewDynamicView("accounting.JournalEntryItemJournalEntrySelectView")}},
			{Key: "accounting.AccountListRoute", Value: lamu.Route{Path: AppUrl + "accounts/", Handler: lamu.NewDynamicView("accounting.AccountListView")}},
			{Key: "accounting.AccountCreateRoute", Value: lamu.Route{Path: AppUrl + "accounts/create", Handler: lamu.NewDynamicView("accounting.AccountCreateView")}},
			{Key: "accounting.AccountDetailRoute", Value: lamu.Route{Path: AppUrl + "accounts/a/{id}/", Handler: lamu.NewDynamicView("accounting.AccountDetailView")}},
			{Key: "accounting.AccountUpdateRoute", Value: lamu.Route{Path: AppUrl + "accounts/a/{id}/edit/", Handler: lamu.NewDynamicView("accounting.AccountUpdateView")}},
			{Key: "accounting.AccountTransferRoute", Value: lamu.Route{Path: AppUrl + "accounts/a/{id}/transfer/", Handler: lamu.NewDynamicView("accounting.AccountTransferView")}},
			{Key: "accounting.AccountTransferToAccountSelectRoute", Value: lamu.Route{Path: AppUrl + "accounts/a/{id}/transfer/select-to-account/", Handler: lamu.NewDynamicView("accounting.AccountTransferToAccountSelectView")}},
			{Key: "accounting.JournalListRoute", Value: lamu.Route{Path: AppUrl + "journals/", Handler: lamu.NewDynamicView("accounting.JournalListView")}},
			{Key: "accounting.JournalCreateRoute", Value: lamu.Route{Path: AppUrl + "journals/create", Handler: lamu.NewDynamicView("accounting.JournalCreateView")}},
			{Key: "accounting.JournalDetailRoute", Value: lamu.Route{Path: AppUrl + "journals/j/{id}/", Handler: lamu.NewDynamicView("accounting.JournalDetailView")}},
			{Key: "accounting.JournalAccountTransferRoute", Value: lamu.Route{Path: AppUrl + "journals/j/{id}/account-transfer/", Handler: lamu.NewDynamicView("accounting.JournalAccountTransferView")}},
			{Key: "accounting.JournalAccountTransferSelectFromRoute", Value: lamu.Route{Path: AppUrl + "journals/j/{id}/account-transfer/select-from/", Handler: lamu.NewDynamicView("accounting.JournalAccountTransferSelectFromView")}},
			{Key: "accounting.JournalAccountTransferSelectToRoute", Value: lamu.Route{Path: AppUrl + "journals/j/{id}/account-transfer/select-to/", Handler: lamu.NewDynamicView("accounting.JournalAccountTransferSelectToView")}},
			{Key: "accounting.JournalEntryListRoute", Value: lamu.Route{Path: AppUrl + "journal-entries/", Handler: lamu.NewDynamicView("accounting.JournalEntryListView")}},
			{Key: "accounting.JournalEntryCreateRoute", Value: lamu.Route{Path: AppUrl + "journal-entries/create/", Handler: lamu.NewDynamicView("accounting.JournalEntryCreateView")}},
			{Key: "accounting.JournalEntryDetailRoute", Value: lamu.Route{Path: AppUrl + "journal-entries/e/{id}/", Handler: lamu.NewDynamicView("accounting.JournalEntryDetailView")}},
			{Key: "accounting.JournalEntryJournalSelectRoute", Value: lamu.Route{Path: AppUrl + "journal-entries/select-journal/", Handler: lamu.NewDynamicView("accounting.JournalEntryJournalSelectView")}},
		},
	}
}
