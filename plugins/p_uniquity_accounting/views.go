package p_uniquity_accounting

import (
	"net/http"

	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/plugins/p_users"
	"github.com/UniquityVentures/lamu/registry"
	"github.com/UniquityVentures/lamu/views"
	"gorm.io/gorm"
)

type accountSelectPreload struct{}

func (accountSelectPreload) Patch(_ views.View, _ *http.Request, query gorm.ChainInterface[Account]) gorm.ChainInterface[Account] {
	return query.Preload("Entity", nil).Preload("Currency", nil)
}

func pluginViews() lamu.PluginFeatures[*views.View] {
	auth := p_users.AuthenticationLayer{}
	return lamu.PluginFeatures[*views.View]{
		Entries: []registry.Pair[string, *views.View]{
			{
				Key: "accounting.JournalEntryItemListView",
				Value: lamu.GetPageView("accounting.JournalEntryItemTable").
					WithLayer("p_users.auth", auth).
					WithLayer("accounting.journal_entry_item_list", views.LayerList[JournalEntryItem]{
						Key: getters.Static("journalEntryItems"),
						QueryPatchers: views.QueryPatchers[JournalEntryItem]{
							{Key: "accounting.journal_entry_item_preload", Value: views.QueryPatcherPreload[JournalEntryItem]{
								Fields: []string{"Account", "JournalEntry.Journal"},
							}},
						},
					}),
			},
			{
				Key: "accounting.JournalEntryItemDetailView",
				Value: lamu.GetPageView("accounting.JournalEntryItemDetail").
					WithLayer("p_users.auth", auth).
					WithLayer("accounting.journal_entry_item_detail", views.LayerDetail[JournalEntryItem]{
						Key:          getters.Static("journalEntryItem"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[JournalEntryItem]{
							{Key: "accounting.journal_entry_item_preload", Value: views.QueryPatcherPreload[JournalEntryItem]{
								Fields: []string{"Account", "JournalEntry.Journal"},
							}},
						},
					}),
			},
			{
				Key: "accounting.JournalEntryItemCreateView",
				Value: lamu.GetPageView("accounting.JournalEntryItemCreateForm").
					WithLayer("p_users.auth", auth).
					WithLayer("accounting.journal_entry_item_create", views.LayerCreate[JournalEntryItem]{
						SuccessURL: lamu.RoutePath("accounting.JournalEntryItemDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("$id")),
						}),
					}),
			},
			{
				Key: "accounting.JournalEntryItemAccountSelectView",
				Value: lamu.GetPageView("accounting.JournalEntryItemAccountSelectionTable").
					WithLayer("p_users.auth", auth).
					WithLayer("accounting.journal_entry_item_account_select", views.LayerList[Account]{
						Key: getters.Static("accounts"),
						QueryPatchers: views.QueryPatchers[Account]{
							{Key: "accounting.journal_entry_item_account_select_preload", Value: accountSelectPreload{}},
						},
					}),
			},
			{
				Key: "accounting.JournalEntryItemJournalEntrySelectView",
				Value: lamu.GetPageView("accounting.JournalEntryItemJournalEntrySelectionTable").
					WithLayer("p_users.auth", auth).
					WithLayer("accounting.journal_entry_item_journal_entry_select", views.LayerList[JournalEntry]{
						Key: getters.Static("journalEntries"),
						QueryPatchers: views.QueryPatchers[JournalEntry]{
							{Key: "accounting.journal_entry_item_journal_entry_select_preload", Value: views.QueryPatcherPreload[JournalEntry]{Fields: []string{"Journal"}}},
						},
					}),
			},
			{
				Key: "accounting.AccountListView",
				Value: lamu.GetPageView("accounting.AccountTable").
					WithLayer("p_users.auth", auth).
					WithLayer("accounting.account_list", views.LayerList[Account]{
						Key: getters.Static("accounts"),
						QueryPatchers: views.QueryPatchers[Account]{
							{Key: "accounting.account_list_preload", Value: views.QueryPatcherPreload[Account]{Fields: []string{"Entity", "Currency"}}},
						},
					}),
			},
			{
				Key: "accounting.AccountDetailView",
				Value: lamu.GetPageView("accounting.AccountDetail").
					WithLayer("p_users.auth", auth).
					WithLayer("accounting.account_detail", views.LayerDetail[Account]{
						Key: getters.Static("account"), PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[Account]{
							{Key: "accounting.account_detail_preload", Value: views.QueryPatcherPreload[Account]{Fields: []string{"Entity", "Currency"}}},
						},
					}),
			},
			{
				Key: "accounting.AccountUpdateView",
				Value: lamu.GetPageView("accounting.AccountUpdateForm").
					WithLayer("p_users.auth", auth).
					WithLayer("accounting.account_update_detail", views.LayerDetail[Account]{
						Key:          getters.Static("account"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[Account]{
							{Key: "accounting.account_update_preload", Value: views.QueryPatcherPreload[Account]{Fields: []string{"Entity", "Currency"}}},
						},
					}).
					WithLayer("accounting.account_update", views.LayerUpdate[Account]{
						Key: getters.Static("account"),
						SuccessURL: lamu.RoutePath("accounting.AccountDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("account.ID")),
						}),
					}),
			},
			{
				Key: "accounting.AccountTransferView",
				Value: lamu.GetPageView("accounting.AccountTransferForm").
					WithLayer("p_users.auth", auth).
					WithLayer("accounting.account_transfer_detail", views.LayerDetail[Account]{
						Key:          getters.Static("account"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[Account]{
							{Key: "accounting.account_transfer_preload", Value: views.QueryPatcherPreload[Account]{Fields: []string{"Entity", "Currency"}}},
						},
					}).
					WithLayer("accounting.account_transfer_post", views.MethodLayer{
						Method:  http.MethodPost,
						Handler: accountTransferPostHandler,
					}),
			},
			{
				Key: "accounting.AccountTransferToAccountSelectView",
				Value: lamu.GetPageView("accounting.AccountTransferToAccountSelectionTable").
					WithLayer("p_users.auth", auth).
					WithLayer("accounting.account_transfer_to_account_select", views.LayerList[Account]{
						Key: getters.Static("accounts"),
						QueryPatchers: views.QueryPatchers[Account]{
							{Key: "accounting.account_transfer_to_select_exclude_source", Value: accountTransferExcludeSourceQueryPatcher{}},
							{Key: "accounting.account_transfer_to_select_preload", Value: accountSelectPreload{}},
						},
					}),
			},
			{
				Key: "accounting.AccountCreateView",
				Value: lamu.GetPageView("accounting.AccountCreateForm").
					WithLayer("p_users.auth", auth).
					WithLayer("accounting.account_create", views.LayerCreate[Account]{
						SuccessURL: lamu.RoutePath("accounting.AccountDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("$id")),
						}),
					}),
			},
			{
				Key: "accounting.JournalListView",
				Value: lamu.GetPageView("accounting.JournalTable").
					WithLayer("p_users.auth", auth).
					WithLayer("accounting.journal_list", views.LayerList[Journal]{Key: getters.Static("journals")}),
			},
			{
				Key: "accounting.JournalDetailView",
				Value: lamu.GetPageView("accounting.JournalDetail").
					WithLayer("p_users.auth", auth).
					WithLayer("accounting.journal_detail", views.LayerDetail[Journal]{
						Key: getters.Static("journal"), PathParamKey: getters.Static("id"),
					}),
			},
			{
				Key: "accounting.JournalAccountTransferView",
				Value: lamu.GetPageView("accounting.JournalAccountTransferForm").
					WithLayer("p_users.auth", auth).
					WithLayer("accounting.journal_account_transfer_detail", views.LayerDetail[Journal]{
						Key:          getters.Static("journal"),
						PathParamKey: getters.Static("id"),
					}).
					WithLayer("accounting.journal_account_transfer_post", views.MethodLayer{
						Method:  http.MethodPost,
						Handler: journalAccountTransferPostHandler,
					}),
			},
			{
				Key: "accounting.JournalAccountTransferSelectFromView",
				Value: lamu.GetPageView("accounting.JournalAccountTransferFromAccountSelectionTable").
					WithLayer("p_users.auth", auth).
					WithLayer("accounting.journal_account_transfer_select_from", views.LayerList[Account]{
						Key: getters.Static("accounts"),
						QueryPatchers: views.QueryPatchers[Account]{
							{Key: "accounting.journal_account_transfer_exclude", Value: journalAccountTransferExcludeByQueryParam{Param: "ToAccountID"}},
							{Key: "accounting.journal_account_transfer_select_preload", Value: accountSelectPreload{}},
						},
					}),
			},
			{
				Key: "accounting.JournalAccountTransferSelectToView",
				Value: lamu.GetPageView("accounting.JournalAccountTransferToAccountSelectionTable").
					WithLayer("p_users.auth", auth).
					WithLayer("accounting.journal_account_transfer_select_to", views.LayerList[Account]{
						Key: getters.Static("accounts"),
						QueryPatchers: views.QueryPatchers[Account]{
							{Key: "accounting.journal_account_transfer_exclude", Value: journalAccountTransferExcludeByQueryParam{Param: "FromAccountID"}},
							{Key: "accounting.journal_account_transfer_select_preload", Value: accountSelectPreload{}},
						},
					}),
			},
			{
				Key: "accounting.JournalCreateView",
				Value: lamu.GetPageView("accounting.JournalCreateForm").
					WithLayer("p_users.auth", auth).
					WithLayer("accounting.journal_create", views.LayerCreate[Journal]{
						SuccessURL: lamu.RoutePath("accounting.JournalDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("$id")),
						}),
					}),
			},
			{
				Key: "accounting.JournalEntryListView",
				Value: lamu.GetPageView("accounting.JournalEntryTable").
					WithLayer("p_users.auth", auth).
					WithLayer("accounting.journal_entry_list", views.LayerList[JournalEntry]{
						Key: getters.Static("journalEntries"),
						QueryPatchers: views.QueryPatchers[JournalEntry]{
							{Key: "accounting.journal_entry_preload", Value: views.QueryPatcherPreload[JournalEntry]{Fields: []string{"Journal"}}},
						},
					}),
			},
			{
				Key: "accounting.JournalEntryDetailView",
				Value: lamu.GetPageView("accounting.JournalEntryDetail").
					WithLayer("p_users.auth", auth).
					WithLayer("accounting.journal_entry_detail", views.LayerDetail[JournalEntry]{
						Key:          getters.Static("journalEntry"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[JournalEntry]{
							{Key: "accounting.journal_entry_preload", Value: views.QueryPatcherPreload[JournalEntry]{Fields: []string{"Journal"}}},
						},
					}),
			},
			{
				Key: "accounting.JournalEntryCreateView",
				Value: lamu.GetPageView("accounting.JournalEntryCreateForm").
					WithLayer("p_users.auth", auth).
					WithLayer("accounting.journal_entry_create", views.LayerCreate[JournalEntry]{
						SuccessURL: lamu.RoutePath("accounting.JournalEntryDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("$id")),
						}),
					}),
			},
			{
				Key: "accounting.JournalEntryJournalSelectView",
				Value: lamu.GetPageView("accounting.JournalEntryJournalSelectionTable").
					WithLayer("p_users.auth", auth).
					WithLayer("accounting.journal_select", views.LayerList[Journal]{Key: getters.Static("journals")}),
			},
		},
	}
}
