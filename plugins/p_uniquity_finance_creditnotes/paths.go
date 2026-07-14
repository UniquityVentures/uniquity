package p_uniquity_finance_creditnotes

import (
	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/registry"
)

func pluginRoutes() lago.PluginFeatures[lago.Route] {
	cn := AppUrl + "cn/"
	return lago.PluginFeatures[lago.Route]{
		Entries: []registry.Pair[string, lago.Route]{
			{Key: "finance_credit_notes.DefaultRoute", Value: lago.Route{Path: AppUrl, Handler: lago.NewDynamicView("finance_credit_notes.CreditNoteListView")}},
			{Key: "finance_credit_notes.CreditNoteCreateRoute", Value: lago.Route{Path: AppUrl + "create/", Handler: lago.NewDynamicView("finance_credit_notes.CreditNoteCreateView")}},
			{Key: "finance_credit_notes.CreditNoteDetailRoute", Value: lago.Route{Path: cn + "{id}/", Handler: lago.NewDynamicView("finance_credit_notes.CreditNoteDetailView")}},
			{Key: "finance_credit_notes.JournalEntryFkSelectRoute", Value: lago.Route{Path: AppUrl + "pick-journal-entry/", Handler: lago.NewDynamicView("finance_credit_notes.JournalEntryFkSelectView")}},
		},
	}
}
