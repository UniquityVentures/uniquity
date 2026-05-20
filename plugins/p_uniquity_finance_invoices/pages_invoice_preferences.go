package p_uniquity_finance_invoices

import (
	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	finance_accounts "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_accounts"
	finance_products "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_products"
)

func invoicePreferencesFormInputs() []components.PageInterface {
	return []components.PageInterface{
		&components.ContainerError{
			Error: getters.Key[error]("$error." + InvoicePrefAccountReceivableIDField),
			Children: []components.PageInterface{
				&components.InputForeignKey[finance_accounts.Account]{
					Label:       "Accounts receivable (invoices)",
					Name:        InvoicePrefAccountReceivableIDField,
					Url:         finance_accounts.AccountSelectRouteURL(finance_accounts.BalanceTypeDebit),
					Display:     getters.Format("%s (#%d)", getters.Any(getters.Key[string]("$in.Name")), getters.Any(getters.Key[uint]("$in.ID"))),
					Placeholder: "Select debit account…",
					Getter:      getters.Association[finance_accounts.Account, uint](finance_products.OptionalPrefUintGetter(InvoicePrefAccountReceivableIDField)),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error." + InvoicePrefAccountRevenueIDField),
			Children: []components.PageInterface{
				&components.InputForeignKey[finance_accounts.Account]{
					Label:       "Revenue account (invoices)",
					Name:        InvoicePrefAccountRevenueIDField,
					Url:         finance_accounts.AccountSelectRouteURL(finance_accounts.BalanceTypeCredit),
					Display:     getters.Format("%s (#%d)", getters.Any(getters.Key[string]("$in.Name")), getters.Any(getters.Key[uint]("$in.ID"))),
					Placeholder: "Select credit account…",
					Getter:      getters.Association[finance_accounts.Account, uint](finance_products.OptionalPrefUintGetter(InvoicePrefAccountRevenueIDField)),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error." + InvoicePrefAccountTaxPayableIDField),
			Children: []components.PageInterface{
				&components.InputForeignKey[finance_accounts.Account]{
					Label:       "Tax payable (invoices)",
					Name:        InvoicePrefAccountTaxPayableIDField,
					Url:         finance_accounts.AccountSelectRouteURL(finance_accounts.BalanceTypeCredit),
					Display:     getters.Format("%s (#%d)", getters.Any(getters.Key[string]("$in.Name")), getters.Any(getters.Key[uint]("$in.ID"))),
					Placeholder: "Select credit account…",
					Getter:      getters.Association[finance_accounts.Account, uint](finance_products.OptionalPrefUintGetter(InvoicePrefAccountTaxPayableIDField)),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error." + InvoicePrefJournalIDField),
			Children: []components.PageInterface{
				&components.InputForeignKey[finance_accounts.Journal]{
					Label:       "Journal (invoices)",
					Name:        InvoicePrefJournalIDField,
					Url:         lamu.RoutePath("finance_accounts.JournalSelectRoute", nil),
					Display:     getters.Key[string]("$in.Name"),
					Placeholder: "Select journal…",
					Getter:      getters.Association[finance_accounts.Journal, uint](finance_products.OptionalPrefUintGetter(InvoicePrefJournalIDField)),
				},
			},
		},
	}
}
