package p_uniquity_finance_invoices

import (
	finance_accounts "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_accounts"
	finance_products "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_products"
	"github.com/lariv-in/lago/components"
	"github.com/lariv-in/lago/getters"
)

func paymentPreferencesFormInputs() []components.PageInterface {
	return []components.PageInterface{
		&components.ContainerError{
			Error: getters.Key[error]("$error." + PaymentPrefAccountIDField),
			Children: []components.PageInterface{
				&components.InputForeignKey[finance_accounts.Account]{
					Label:       "Payment account (receipts)",
					Name:        PaymentPrefAccountIDField,
					Url:         finance_accounts.AccountSelectRouteURL(finance_accounts.BalanceTypeDebit),
					Display:     getters.Format("%s (#%d)", getters.Any(getters.Key[string]("$in.Name")), getters.Any(getters.Key[uint]("$in.ID"))),
					Placeholder: "Bank or cash account…",
					Getter:      getters.Association[finance_accounts.Account, uint](finance_products.OptionalPrefUintGetter(PaymentPrefAccountIDField)),
				},
			},
		},
	}
}
