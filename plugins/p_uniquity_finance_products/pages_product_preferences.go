package p_uniquity_finance_products

import (
	finance_accounts "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_accounts"
	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/components"
	"github.com/lariv-in/lago/getters"
)

func productPreferencesFormInputs() []components.PageInterface {
	return []components.PageInterface{
		&components.ContainerError{
			Error: getters.Key[error]("$error." + ProductPrefInventoryAccountIDField),
			Children: []components.PageInterface{
				&components.InputForeignKey[finance_accounts.Account]{
					Label:       "Inventory account (products)",
					Name:        ProductPrefInventoryAccountIDField,
					Url:         lago.RoutePath("finance_accounts.AccountSelectRoute", nil),
					Display:     getters.Format("%s (#%d)", getters.Any(getters.Key[string]("$in.Name")), getters.Any(getters.Key[uint]("$in.ID"))),
					Placeholder: "Select…",
					Getter:      getters.Association[finance_accounts.Account, uint](OptionalPrefUintGetter(ProductPrefInventoryAccountIDField)),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error." + ProductPrefCostOfSalesAcctIDField),
			Children: []components.PageInterface{
				&components.InputForeignKey[finance_accounts.Account]{
					Label:       "Cost of sales account (products)",
					Name:        ProductPrefCostOfSalesAcctIDField,
					Url:         lago.RoutePath("finance_accounts.AccountSelectRoute", nil),
					Display:     getters.Format("%s (#%d)", getters.Any(getters.Key[string]("$in.Name")), getters.Any(getters.Key[uint]("$in.ID"))),
					Placeholder: "Select…",
					Getter:      getters.Association[finance_accounts.Account, uint](OptionalPrefUintGetter(ProductPrefCostOfSalesAcctIDField)),
				},
			},
		},
	}
}
