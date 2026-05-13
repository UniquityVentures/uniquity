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
	return query
}

func pluginViews() lamu.PluginFeatures[*views.View] {
	auth := p_users.AuthenticationLayer{}
	return lamu.PluginFeatures[*views.View]{
		Entries: []registry.Pair[string, *views.View]{
			{
				Key: "accounting.TransactionListView",
				Value: lamu.GetPageView("accounting.TransactionTable").
					WithLayer("p_users.auth", auth).
					WithLayer("accounting.transaction_list", views.LayerList[Posting]{
						Key: getters.Static("transactions"),
						QueryPatchers: views.QueryPatchers[Posting]{
							{Key: "accounting.transaction_preload", Value: views.QueryPatcherPreload[Posting]{Fields: []string{"Account"}}},
						},
					}),
			},
			{
				Key: "accounting.TransactionDetailView",
				Value: lamu.GetPageView("accounting.TransactionDetail").
					WithLayer("p_users.auth", auth).
					WithLayer("accounting.transaction_detail", views.LayerDetail[Posting]{
						Key:          getters.Static("transaction"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[Posting]{
							{Key: "accounting.transaction_preload", Value: views.QueryPatcherPreload[Posting]{}},
						},
					}),
			},
			{
				Key: "accounting.TransactionCreateView",
				Value: lamu.GetPageView("accounting.TransactionCreateForm").
					WithLayer("p_users.auth", auth).
					WithLayer("accounting.transaction_create", views.LayerCreate[Posting]{
						SuccessURL: lamu.RoutePath("accounting.TransactionDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("$id")),
						}),
					}),
			},
			{
				Key: "accounting.TransactionAccountSelectView",
				Value: lamu.GetPageView("accounting.TransactionAccountSelectionTable").
					WithLayer("p_users.auth", auth).
					WithLayer("accounting.account_select", views.LayerList[Account]{
						Key: getters.Static("accounts"),
						QueryPatchers: views.QueryPatchers[Account]{
							{Key: "accounting.account_select_preload", Value: accountSelectPreload{}},
						},
					}),
			},
			{
				Key: "accounting.AccountListView",
				Value: lamu.GetPageView("accounting.AccountTable").
					WithLayer("p_users.auth", auth).
					WithLayer("accounting.account_list", views.LayerList[Account]{Key: getters.Static("accounts")}),
			},
			{
				Key: "accounting.AccountDetailView",
				Value: lamu.GetPageView("accounting.AccountDetail").
					WithLayer("p_users.auth", auth).
					WithLayer("accounting.account_detail", views.LayerDetail[Account]{
						Key: getters.Static("account"), PathParamKey: getters.Static("id"),
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
		},
	}
}
