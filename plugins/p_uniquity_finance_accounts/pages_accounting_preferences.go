package p_uniquity_finance_accounts

import (
	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pageAccountingPreferencesPages() []registry.Pair[string, components.PageInterface] {
	formName := getters.Static("finance_accounts.AccountingPreferencesForm")
	return []registry.Pair[string, components.PageInterface]{
		{Key: "finance_accounts.AccountingPreferencesForm", Value: &components.ShellScaffold{
			Page:    components.Page{Roles: []string{"superuser"}},
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_accounts.MainMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name:      formName,
					ActionURL: lamu.RoutePath("finance_accounts.AccountingPreferencesRoute", nil),
					Children: []components.PageInterface{
						&components.FormComponent[AccountingPreferences]{
							Attr:          getters.FormBubbling(formName),
							Title:         "Accounting preferences",
							Subtitle:      "Go text/template for posted invoice numbers when a draft has no number. Variables: FISCAL_CODE, YY, YYYY, POSTED_SEQ (next posted row id). Example: INV-{{.YYYY}}-{{.POSTED_SEQ}}. Default journal prefills new draft invoices. Invoice PDF template: Go text/template → Typst; root fields match each detail page’s $in; funcs num2words, num2wordsAnd, num2wordsRupees, invoiceGrandTotalWords.",
							ChildrenInput: accountingPreferencesFormInputs(),
							ChildrenAction: []components.PageInterface{
								&components.ButtonSubmit{Label: "Save preferences"},
							},
						},
					},
				},
			},
		}},
	}
}

func accountingPreferencesFormInputs() []components.PageInterface {
	return []components.PageInterface{
		&components.ContainerError{
			Error: getters.Key[error]("$error.InvoiceNumberFormat"),
			Children: []components.PageInterface{
				&components.InputText{
					Name:   "InvoiceNumberFormat",
					Label:  "Invoice number format",
					Getter: getters.Key[string]("$in.InvoiceNumberFormat"),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.InvoicePDFTemplate"),
			Children: []components.PageInterface{
				&components.InputTextarea{
					Name:    "InvoicePDFTemplate",
					Label:   "Invoice PDF template (Typst)",
					Getter:  getters.Key[string]("$in.InvoicePDFTemplate"),
					Rows:    16,
					Classes: "font-mono text-sm min-h-48",
				},
			},
		},
	}
}
