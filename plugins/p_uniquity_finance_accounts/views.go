package p_uniquity_finance_accounts

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/fields"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/plugins/p_users"
	"github.com/UniquityVentures/lamu/registry"
	"github.com/UniquityVentures/lamu/views"
	"gorm.io/gorm"
)


// accountChildrenContextKey holds [components.ObjectList[Account]] for direct children on the detail page.
const accountChildrenContextKey = "accountChildren"

// accountBalanceTotalContextKey holds a formatted sum of [JournalEntryItem] amounts for this account
// and all descendant accounts (subtree), for the account detail page.
const accountBalanceTotalContextKey = "accountBalanceTotal"

type accountListPreload struct{}

func (accountListPreload) Patch(_ views.View, _ *http.Request, query gorm.ChainInterface[Account]) gorm.ChainInterface[Account] {
	return query.Preload("Parent", nil)
}

// accountListRootOnly scopes the main accounts list to top-level rows (no parent).
type accountListRootOnly struct{}

func (accountListRootOnly) Patch(_ views.View, _ *http.Request, query gorm.ChainInterface[Account]) gorm.ChainInterface[Account] {
	return query.Where("parent_id IS NULL")
}

// accountDetailChildrenLayer loads direct child accounts after the parent row is in context.
type accountDetailChildrenLayer struct{}

func (accountDetailChildrenLayer) Next(_ views.View, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		acc, ok := ctx.Value("account").(Account)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}
		if !acc.IsGroup {
			ctx = context.WithValue(ctx, accountChildrenContextKey, components.ObjectList[Account]{})
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		db, err := getters.DBFromContext(ctx)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		items, err := gorm.G[Account](db).Where("parent_id = ?", acc.ID).Order("code ASC").Find(ctx)
		if err != nil {
			slog.Error("finance_accounts.account_detail_children: query", "error", err, "parent_id", acc.ID)
			next.ServeHTTP(w, r)
			return
		}
		n := uint64(len(items))
		ol := components.ObjectList[Account]{
			Items:    items,
			Number:   1,
			NumPages: 1,
			Total:    n,
		}
		ol = prependAccountParentUpRow(ol)
		ctx = context.WithValue(ctx, accountChildrenContextKey, ol)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// accountDescendantIDs returns rootID and every descendant account id (BFS), unique.
func accountDescendantIDs(db *gorm.DB, rootID uint) ([]uint, error) {
	var out []uint
	queue := []uint{rootID}
	seen := map[uint]struct{}{rootID: {}}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		out = append(out, cur)
		var kids []uint
		if err := db.Model(&Account{}).Where("parent_id = ?", cur).Pluck("id", &kids).Error; err != nil {
			return nil, err
		}
		for _, k := range kids {
			if _, ok := seen[k]; ok {
				continue
			}
			seen[k] = struct{}{}
			queue = append(queue, k)
		}
	}
	return out, nil
}

// accountDetailBalanceLayer sums journal line amounts for the account subtree into [accountBalanceTotalContextKey].
type accountDetailBalanceLayer struct{}

func (accountDetailBalanceLayer) Next(_ views.View, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		acc, ok := ctx.Value("account").(Account)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}
		db, err := getters.DBFromContext(ctx)
		if err != nil {
			slog.Error("finance_accounts.account_detail_balance: db", "error", err)
			ctx = context.WithValue(ctx, accountBalanceTotalContextKey, "—")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		ids, err := accountDescendantIDs(db, acc.ID)
		if err != nil {
			slog.Error("finance_accounts.account_detail_balance: descendants", "error", err, "account_id", acc.ID)
			ctx = context.WithValue(ctx, accountBalanceTotalContextKey, "—")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		var row struct {
			Total fields.DecimalSix `gorm:"column:total"`
		}
		q := db.Model(&JournalEntryItem{}).Where("account_id IN ?", ids)
		if err := q.Select("COALESCE(SUM(amount), 0) AS total").Scan(&row).Error; err != nil {
			slog.Error("finance_accounts.account_detail_balance: sum", "error", err, "account_id", acc.ID)
			ctx = context.WithValue(ctx, accountBalanceTotalContextKey, "—")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		ctx = context.WithValue(ctx, accountBalanceTotalContextKey, row.Total.String())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type accountDetailPreload struct{}

func (accountDetailPreload) Patch(_ views.View, _ *http.Request, query gorm.ChainInterface[Account]) gorm.ChainInterface[Account] {
	return query.Preload("Parent", nil)
}

type currencyListOrder struct{}

func (currencyListOrder) Patch(_ views.View, _ *http.Request, query gorm.ChainInterface[Currency]) gorm.ChainInterface[Currency] {
	return query.Order("code ASC")
}

type journalListPreload struct{}

func (journalListPreload) Patch(_ views.View, _ *http.Request, query gorm.ChainInterface[Journal]) gorm.ChainInterface[Journal] {
	return query.Preload("Currency", nil)
}

type journalDetailPreload struct{}

func (journalDetailPreload) Patch(_ views.View, _ *http.Request, query gorm.ChainInterface[Journal]) gorm.ChainInterface[Journal] {
	return query.Preload("Currency", nil)
}

// journalDetailEntriesContextKey holds [JournalEntry] rows for the journal detail page.
const journalDetailEntriesContextKey = "journalEntries"

// journalDetailEntriesLayer loads [JournalEntry] rows for the current journal.
type journalDetailEntriesLayer struct{}

func (journalDetailEntriesLayer) Next(_ views.View, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		journal, ok := ctx.Value("journal").(Journal)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}
		db, err := getters.DBFromContext(ctx)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		items, err := gorm.G[JournalEntry](db).
			Where("journal_id = ?", journal.ID).
			Preload("SourceDoc", nil).
			Order("journal_entries.datetime DESC").
			Order("journal_entries.id DESC").
			Find(ctx)
		if err != nil {
			slog.Error("finance_accounts.journal_detail_entries", "error", err, "journal_id", journal.ID)
			next.ServeHTTP(w, r)
			return
		}
		n := uint64(len(items))
		ol := components.ObjectList[JournalEntry]{
			Items:    items,
			Number:   1,
			NumPages: 1,
			Total:    n,
		}
		ctx = context.WithValue(ctx, journalDetailEntriesContextKey, ol)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// journalCreateFormDefaults applies defaults before insert (checkbox / enum omissions on POST).
type journalCreateFormDefaults struct{}

func (journalCreateFormDefaults) Patch(_ views.View, _ *http.Request, formData map[string]any, formErrors map[string]error) (map[string]any, map[string]error) {
	if _, ok := formData["IsActive"]; !ok {
		formData["IsActive"] = true
	}
	if v, ok := formData["Type"]; !ok || v == nil || v == "" {
		formData["Type"] = JournalTypeGeneral
	}
	return formData, formErrors
}

// journalEntryCreateFormDefaults sets JournalID from the parent journal in context and a missing Datetime.
type journalEntryCreateFormDefaults struct{}

func (journalEntryCreateFormDefaults) Patch(_ views.View, r *http.Request, formData map[string]any, formErrors map[string]error) (map[string]any, map[string]error) {
	j, ok := r.Context().Value("journal").(Journal)
	if ok {
		formData["JournalID"] = j.ID
	}
	if v, ok := formData["Datetime"]; !ok || v == nil {
		formData["Datetime"] = time.Now()
	}
	return formData, formErrors
}

// journalEntryDetailItemsContextKey holds line items for the journal entry detail page.
const journalEntryDetailItemsContextKey = "journalEntryItems"

type journalEntryDetailPreload struct{}

func (journalEntryDetailPreload) Patch(_ views.View, _ *http.Request, query gorm.ChainInterface[JournalEntry]) gorm.ChainInterface[JournalEntry] {
	return query.Preload("Journal", nil).Preload("SourceDoc", nil)
}

// journalEntryDetailItemsLayer loads [JournalEntryItem] rows for the current journal entry.
type journalEntryDetailItemsLayer struct{}

func (journalEntryDetailItemsLayer) Next(_ views.View, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		entry, ok := ctx.Value("journalEntry").(JournalEntry)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}
		db, err := getters.DBFromContext(ctx)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		items, err := gorm.G[JournalEntryItem](db).
			Where("journal_entry_id = ?", entry.ID).
			Preload("Account", nil).
			Order("journal_entry_items.id ASC").
			Find(ctx)
		if err != nil {
			slog.Error("finance_accounts.journal_entry_detail_items", "error", err, "journal_entry_id", entry.ID)
			next.ServeHTTP(w, r)
			return
		}
		n := uint64(len(items))
		ol := components.ObjectList[JournalEntryItem]{
			Items:    items,
			Number:   1,
			NumPages: 1,
			Total:    n,
		}
		ctx = context.WithValue(ctx, journalEntryDetailItemsContextKey, ol)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func pluginViews() lamu.PluginFeatures[*views.View] {
	return lamu.PluginFeatures[*views.View]{
		Entries: []registry.Pair[string, *views.View]{
			{
				Key: "finance_accounts.AccountListView",
				Value: lamu.GetPageView("finance_accounts.AccountTable").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_accounts.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("finance_accounts.account_list", views.LayerList[Account]{
						Key: getters.Static("accounts"),
						QueryPatchers: views.QueryPatchers[Account]{
							{Key: "finance_accounts.list_root_only", Value: accountListRootOnly{}},
							{Key: "finance_accounts.preload_parent", Value: accountListPreload{}},
						},
					}),
			},
			{
				Key: "finance_accounts.AccountDetailView",
				Value: lamu.GetPageView("finance_accounts.AccountDetail").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_accounts.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("finance_accounts.account_detail", views.LayerDetail[Account]{
						Key:          getters.Static("account"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[Account]{
							{Key: "finance_accounts.preload_parent", Value: accountDetailPreload{}},
						},
					}).
					WithLayer("finance_accounts.account_detail_children", accountDetailChildrenLayer{}).
					WithLayer("finance_accounts.account_detail_balance", accountDetailBalanceLayer{}),
			},
			{
				Key: "finance_accounts.AccountCreateView",
				Value: lamu.GetPageView("finance_accounts.AccountCreateForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_accounts.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("finance_accounts.account_create", views.LayerCreate[Account]{
						SuccessURL: lamu.RoutePath("finance_accounts.AccountDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("$id")),
						}),
					}),
			},
			{
				Key: "finance_accounts.AccountUpdateView",
				Value: lamu.GetPageView("finance_accounts.AccountUpdateForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_accounts.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("finance_accounts.account_detail", views.LayerDetail[Account]{
						Key:          getters.Static("account"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[Account]{
							{Key: "finance_accounts.preload_parent", Value: accountDetailPreload{}},
						},
					}).
					WithLayer("finance_accounts.account_update", views.LayerUpdate[Account]{
						Key:        getters.Static("account"),
						SuccessURL: lamu.RoutePath("finance_accounts.DefaultRoute", nil),
					}),
			},
			{
				Key: "finance_accounts.AccountDeleteView",
				Value: lamu.GetPageView("finance_accounts.AccountDeleteForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_accounts.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("finance_accounts.account_detail", views.LayerDetail[Account]{
						Key:          getters.Static("account"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[Account]{
							{Key: "finance_accounts.preload_parent", Value: accountDetailPreload{}},
						},
					}).
					WithLayer("finance_accounts.account_delete", views.LayerDelete[Account]{
						Key:        getters.Static("account"),
						SuccessURL: lamu.RoutePath("finance_accounts.DefaultRoute", nil),
					}),
			},
			{
				Key: "finance_accounts.AccountSelectView",
				Value: lamu.GetPageView("finance_accounts.AccountSelectionTable").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_accounts.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("finance_accounts.account_select_list", views.LayerList[Account]{
						Key: getters.Static("accounts"),
						QueryPatchers: views.QueryPatchers[Account]{
							{Key: "finance_accounts.preload_parent", Value: accountListPreload{}},
							{Key: "finance_accounts.account_select_balance_type_scope", Value: accountSelectBalanceTypeScope{}},
						},
					}).
					WithLayer("finance_accounts.account_select_parent_up", accountSelectParentUpLayer{}),
			},
			{
				Key: "finance_accounts.CurrencyListView",
				Value: lamu.GetPageView("finance_accounts.CurrencyTable").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_accounts.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("finance_accounts.currency_list", views.LayerList[Currency]{
						Key: getters.Static("currencies"),
						QueryPatchers: views.QueryPatchers[Currency]{
							{Key: "finance_accounts.currency_order", Value: currencyListOrder{}},
						},
					}),
			},
			{
				Key: "finance_accounts.CurrencyDetailView",
				Value: lamu.GetPageView("finance_accounts.CurrencyDetail").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_accounts.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("finance_accounts.currency_detail", views.LayerDetail[Currency]{
						Key:          getters.Static("currency"),
						PathParamKey: getters.Static("id"),
					}),
			},
			{
				Key: "finance_accounts.CurrencyCreateView",
				Value: lamu.GetPageView("finance_accounts.CurrencyCreateForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_accounts.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("finance_accounts.currency_create", views.LayerCreate[Currency]{
						SuccessURL: lamu.RoutePath("finance_accounts.CurrencyDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("$id")),
						}),
					}),
			},
			{
				Key: "finance_accounts.CurrencyUpdateView",
				Value: lamu.GetPageView("finance_accounts.CurrencyUpdateForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_accounts.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("finance_accounts.currency_detail", views.LayerDetail[Currency]{
						Key:          getters.Static("currency"),
						PathParamKey: getters.Static("id"),
					}).
					WithLayer("finance_accounts.currency_update", views.LayerUpdate[Currency]{
						Key:        getters.Static("currency"),
						SuccessURL: lamu.RoutePath("finance_accounts.CurrencyListRoute", nil),
					}),
			},
			{
				Key: "finance_accounts.CurrencyDeleteView",
				Value: lamu.GetPageView("finance_accounts.CurrencyDeleteForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_accounts.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("finance_accounts.currency_detail", views.LayerDetail[Currency]{
						Key:          getters.Static("currency"),
						PathParamKey: getters.Static("id"),
					}).
					WithLayer("finance_accounts.currency_delete", views.LayerDelete[Currency]{
						Key:        getters.Static("currency"),
						SuccessURL: lamu.RoutePath("finance_accounts.CurrencyListRoute", nil),
					}),
			},
			{
				Key: "finance_accounts.CurrencySelectView",
				Value: lamu.GetPageView("finance_accounts.CurrencySelectionTable").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_accounts.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("finance_accounts.currency_select_list", views.LayerList[Currency]{
						Key: getters.Static("currencies"),
						QueryPatchers: views.QueryPatchers[Currency]{
							{Key: "finance_accounts.currency_order", Value: currencyListOrder{}},
						},
					}),
			},
			{
				Key: "finance_accounts.JournalListView",
				Value: lamu.GetPageView("finance_accounts.JournalTable").
				 WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_accounts.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("finance_accounts.journal_list", views.LayerList[Journal]{
						Key: getters.Static("journals"),
						QueryPatchers: views.QueryPatchers[Journal]{
							{Key: "finance_accounts.journal_preload_currency", Value: journalListPreload{}},
						},
					}),
			},
			{
				Key: "finance_accounts.JournalSelectView",
				Value: lamu.GetPageView("finance_accounts.JournalFkSelectionTable").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_accounts.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("finance_accounts.journal_fk_list", views.LayerList[Journal]{
						Key: getters.Static("journals"),
						QueryPatchers: views.QueryPatchers[Journal]{
							{Key: "finance_accounts.journal_preload_currency", Value: journalListPreload{}},
						},
					}),
			},
			{
				Key: "finance_accounts.JournalDetailView",
				Value: lamu.GetPageView("finance_accounts.JournalDetail").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_accounts.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("finance_accounts.journal_detail", views.LayerDetail[Journal]{
						Key:          getters.Static("journal"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[Journal]{
							{Key: "finance_accounts.journal_preload_currency", Value: journalDetailPreload{}},
						},
					}).
					WithLayer("finance_accounts.journal_entries", journalDetailEntriesLayer{}),
			},
			{
				Key: "finance_accounts.JournalCreateView",
				Value: lamu.GetPageView("finance_accounts.JournalCreateForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_accounts.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("finance_accounts.journal_create", views.LayerCreate[Journal]{
						SuccessURL: lamu.RoutePath("finance_accounts.JournalDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("$id")),
						}),
						FormPatchers: views.FormPatchers{
							{Key: "finance_accounts.journal_create_defaults", Value: journalCreateFormDefaults{}},
						},
					}),
			},
			{
				Key: "finance_accounts.JournalUpdateView",
				Value: lamu.GetPageView("finance_accounts.JournalUpdateForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_accounts.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("finance_accounts.journal_detail", views.LayerDetail[Journal]{
						Key:          getters.Static("journal"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[Journal]{
							{Key: "finance_accounts.journal_preload_currency", Value: journalDetailPreload{}},
						},
					}).
					WithLayer("finance_accounts.journal_update", views.LayerUpdate[Journal]{
						Key:        getters.Static("journal"),
						SuccessURL: lamu.RoutePath("finance_accounts.JournalListRoute", nil),
					}),
			},
			{
				Key: "finance_accounts.JournalDeleteView",
				Value: lamu.GetPageView("finance_accounts.JournalDeleteForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_accounts.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("finance_accounts.journal_detail", views.LayerDetail[Journal]{
						Key:          getters.Static("journal"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[Journal]{
							{Key: "finance_accounts.journal_preload_currency", Value: journalDetailPreload{}},
						},
					}).
					WithLayer("finance_accounts.journal_delete", views.LayerDelete[Journal]{
						Key:        getters.Static("journal"),
						SuccessURL: lamu.RoutePath("finance_accounts.JournalListRoute", nil),
					}),
			},
			{
				Key: "finance_accounts.SourceDocSelectView",
				Value: lamu.GetPageView("finance_accounts.SourceDocSelectionTable").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_accounts.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("finance_accounts.sourcedoc_select_list", views.LayerList[SourceDoc]{
						Key: getters.Static("source_docs"),
					}),
			},
			{
				Key: "finance_accounts.JournalEntryCreateView",
				Value: lamu.GetPageView("finance_accounts.JournalEntryCreateForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_accounts.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("finance_accounts.journal_for_entry_create", views.LayerDetail[Journal]{
						Key:          getters.Static("journal"),
						PathParamKey: getters.Static("journal_id"),
						QueryPatchers: views.QueryPatchers[Journal]{
							{Key: "finance_accounts.journal_preload_currency", Value: journalDetailPreload{}},
						},
					}).
					WithLayer("finance_accounts.journal_entry_create", views.LayerCreate[JournalEntry]{
						SuccessURL: lamu.RoutePath("finance_accounts.JournalDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("journal.ID")),
						}),
						FormPatchers: views.FormPatchers{
							{Key: "finance_accounts.journal_entry_create_defaults", Value: journalEntryCreateFormDefaults{}},
						},
					}),
			},
			{
				Key: "finance_accounts.JournalEntryDetailView",
				Value: lamu.GetPageView("finance_accounts.JournalEntryDetail").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_accounts.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("finance_accounts.journal_entry_detail", views.LayerDetail[JournalEntry]{
						Key:          getters.Static("journalEntry"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[JournalEntry]{
							{Key: "finance_accounts.journal_entry_detail_preload", Value: journalEntryDetailPreload{}},
						},
					}).
					WithLayer("finance_accounts.journal_entry_items", journalEntryDetailItemsLayer{}),
			},
			{
				Key: "finance_accounts.AccountingPreferencesView",
				Value: lamu.GetPageView("finance_accounts.AccountingPreferencesForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_accounts.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("finance_accounts.accounting_preferences", views.LayerSingleton[AccountingPreferences]{
						SuccessURL: lamu.RoutePath("finance_accounts.AccountingPreferencesRoute", nil),
					}),
			},
		},
	}
}
