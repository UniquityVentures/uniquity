package p_uniquity_finance_accounts

import (
	"context"

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
							Subtitle:      "Go text/template for posted invoice numbers when a draft has no number. Variables: FISCAL_CODE, YY, YYYY, POSTED_SEQ (next posted row id). Example: INV-{{.YYYY}}-{{.POSTED_SEQ}}. Default journal prefills new draft invoices.",
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
			Error: getters.Key[error]("$error.DefaultJournalID"),
			Children: []components.PageInterface{
				&components.InputForeignKey[Journal]{
					Name:        "DefaultJournalID",
					Label:       "Default journal (draft invoices)",
					Required:    false,
					Url:         lamu.RoutePath("finance_accounts.JournalSelectRoute", nil),
					Display:     getters.Key[string]("$in.Name"),
					Placeholder: "None — user must choose each time",
					Getter:      getters.Association[Journal, uint](accountingPrefsDefaultJournalIDUintGetter()),
				},
			},
		},
	}
}

func accountingPrefsDefaultJournalIDUintGetter() getters.Getter[uint] {
	return func(ctx context.Context) (uint, error) {
		ptr, err := getters.Key[*uint]("$in.DefaultJournalID")(ctx)
		if err != nil || ptr == nil {
			return 0, nil
		}
		return *ptr, nil
	}
}
