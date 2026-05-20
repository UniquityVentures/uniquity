package p_uniquity_finance_invoices

import (
	"slices"

	"github.com/UniquityVentures/lamu/components"
	finance_accounts "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_accounts"
)

func accountingPreferencesFormHasInvoicePrefs(parent components.ParentInterface) bool {
	for _, input := range components.FindInputs(parent) {
		switch input.GetName() {
		case InvoicePrefAccountReceivableIDField, InvoicePrefAccountRevenueIDField, InvoicePrefAccountTaxPayableIDField, InvoicePrefJournalIDField,
			PaymentPrefAccountIDField:
			return true
		}
	}
	return false
}

func patchAccountingPreferencesForm(page components.PageInterface) components.PageInterface {
	shell, ok := page.(*components.ShellScaffold)
	if !ok {
		return page
	}
	cloned := *shell
	if accountingPreferencesFormHasInvoicePrefs(&cloned) {
		return &cloned
	}
	appendInvoicePrefsToForm(&cloned)
	return &cloned
}

func appendInvoicePrefsToForm(parent components.MutableParentInterface) {
	children := parent.GetChildren()
	for i, child := range children {
		switch fc := child.(type) {
		case *components.FormComponent[finance_accounts.AccountingPreferences]:
			next := *fc
			next.ChildrenInput = append(slices.Clone(fc.ChildrenInput), invoicePreferencesFormInputs()...)
			next.ChildrenInput = append(next.ChildrenInput, paymentPreferencesFormInputs()...)
			children[i] = &next
		case components.FormComponent[finance_accounts.AccountingPreferences]:
			next := fc
			next.ChildrenInput = append(slices.Clone(fc.ChildrenInput), invoicePreferencesFormInputs()...)
			next.ChildrenInput = append(next.ChildrenInput, paymentPreferencesFormInputs()...)
			children[i] = next
		default:
			if mp, ok := child.(components.MutableParentInterface); ok {
				appendInvoicePrefsToForm(mp)
			}
		}
	}
	parent.SetChildren(children)
}
