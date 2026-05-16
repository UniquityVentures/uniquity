package p_uniquity_finance_creditnotes

import (
	"log/slog"
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

// SuperuserOnlyLayer returns 401 unless the authenticated user is a superuser.
type SuperuserOnlyLayer struct{}

func (SuperuserOnlyLayer) Next(_ views.View, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := p_users.UserFromContext(r.Context(), "finance_credit_notes.SuperuserOnlyLayer")
		if !user.IsSuperuser {
			slog.Error("finance_credit_notes.SuperuserOnlyLayer: forbidden", "user_id", user.ID)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

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
					WithLayer("finance_credit_notes.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_credit_notes.credit_note_list", views.LayerList[CreditNote]{
						Key: getters.Static("credit_notes"),
					}),
			},
			{
				Key: "finance_credit_notes.CreditNoteCreateView",
				Value: lamu.GetPageView("finance_credit_notes.CreditNoteCreateForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_credit_notes.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_credit_notes.credit_note_create", views.LayerCreate[CreditNote]{
						SuccessURL: lamu.RoutePath("finance_credit_notes.DefaultRoute", nil),
						FormPatchers: views.FormPatchers{
							{Key: "finance_credit_notes.credit_note_create_defaults", Value: creditNoteCreateFormDefaults{}},
						},
					}),
			},
			{
				Key: "finance_credit_notes.JournalEntryFkSelectView",
				Value: lamu.GetPageView("finance_credit_notes.JournalEntryFkSelectionTable").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_credit_notes.superuser", SuperuserOnlyLayer{}).
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
