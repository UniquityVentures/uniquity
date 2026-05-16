package p_uniquity_finance_creditnotes

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pluginRoutes() lamu.PluginFeatures[lamu.Route] {
	return lamu.PluginFeatures[lamu.Route]{
		Entries: []registry.Pair[string, lamu.Route]{
			{Key: "finance_credit_notes.DefaultRoute", Value: lamu.Route{Path: AppUrl, Handler: lamu.NewDynamicView("finance_credit_notes.CreditNoteListView")}},
			{Key: "finance_credit_notes.CreditNoteCreateRoute", Value: lamu.Route{Path: AppUrl + "create/", Handler: lamu.NewDynamicView("finance_credit_notes.CreditNoteCreateView")}},
			{Key: "finance_credit_notes.JournalEntryFkSelectRoute", Value: lamu.Route{Path: AppUrl + "pick-journal-entry/", Handler: lamu.NewDynamicView("finance_credit_notes.JournalEntryFkSelectView")}},
		},
	}
}
