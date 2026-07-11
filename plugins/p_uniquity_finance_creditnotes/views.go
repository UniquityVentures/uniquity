package p_uniquity_finance_creditnotes

import (
	"net/http"
	"time"

	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/plugins/p_users"
	"github.com/UniquityVentures/lamu/registry"
	"github.com/UniquityVentures/lamu/views"
	finance_accounts "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_accounts"
	"gorm.io/gorm"
)


type creditNoteCreateFormDefaults struct{}

func (creditNoteCreateFormDefaults) Patch(_ views.View, _ *http.Request, formData map[string]any, formErrors map[string]error) (map[string]any, map[string]error) {
	if v, ok := formData["Datetime"]; !ok || v == nil {
		formData["Datetime"] = time.Now()
	}
	return formData, formErrors
}

type journalEntryFkListPreload struct{}

func (journalEntryFkListPreload) Patch(_ views.View, _ *http.Request, query gorm.ChainInterface[finance_accounts.JournalEntry]) gorm.ChainInterface[finance_accounts.JournalEntry] {
	return query.
		Preload("Journal", nil).
		Order("journal_entries.datetime DESC").
		Order("journal_entries.id DESC")
}

func pluginViews() lamu.PluginFeatures[*views.View] {
	return lamu.PluginFeatures[*views.View]{
		Entries: []registry.Pair[string, *views.View]{
			{
				Key: "finance_credit_notes.CreditNoteListView",
				Value: lamu.GetPageView("finance_credit_notes.CreditNoteTable").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_credit_notes.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("finance_credit_notes.credit_note_list", views.LayerList[CreditNote]{
						Key: getters.Static("credit_notes"),
					}),
			},
			{
				Key: "finance_credit_notes.CreditNoteCreateView",
				Value: lamu.GetPageView("finance_credit_notes.CreditNoteCreateForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_credit_notes.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("finance_credit_notes.credit_note_create", views.LayerCreate[CreditNote]{
						SuccessURL: lamu.RoutePath("finance_credit_notes.DefaultRoute", nil),
						FormPatchers: views.FormPatchers{
							{Key: "finance_credit_notes.credit_note_create_defaults", Value: creditNoteCreateFormDefaults{}},
						},
					}),
			},
			{
				Key: "finance_credit_notes.CreditNoteDetailView",
				Value: lamu.GetPageView("finance_credit_notes.CreditNoteDetail").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_credit_notes.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("finance_credit_notes.credit_note_detail", views.LayerDetail[CreditNote]{
						Key:          getters.Static("credit_note"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[CreditNote]{
							{Key: "finance_credit_notes.preload_credit_note_detail", Value: views.QueryPatcherPreload[CreditNote]{Fields: []string{"JournalEntry", "JournalEntry.Journal", "ReversedJournalEntry", "ReversedJournalEntry.Journal"}}},
						},
					}),
			},
			{
				Key: "finance_credit_notes.JournalEntryFkSelectView",
				Value: lamu.GetPageView("finance_credit_notes.JournalEntryFkSelectionTable").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_credit_notes.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("finance_credit_notes.journal_entry_fk_list", views.LayerList[finance_accounts.JournalEntry]{
						Key: getters.Static("journal_entries"),
						QueryPatchers: views.QueryPatchers[finance_accounts.JournalEntry]{
							{Key: "finance_credit_notes.journal_entry_fk_preload", Value: journalEntryFkListPreload{}},
						},
					}),
			},
		},
	}
}
