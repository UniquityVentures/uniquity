package p_uniquity_finance_products

import (
	"slices"

	finance_accounts "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_accounts"
	"github.com/lariv-in/lago/components"
)

func accountingPreferencesFormHasProductPrefs(parent components.ParentInterface) bool {
	for _, input := range components.FindInputs(parent) {
		switch input.GetName() {
		case ProductPrefInventoryAccountIDField, ProductPrefCostOfSalesAcctIDField:
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
	if accountingPreferencesFormHasProductPrefs(&cloned) {
		return &cloned
	}
	appendProductPrefsToForm(&cloned)
	return &cloned
}

func appendProductPrefsToForm(parent components.MutableParentInterface) {
	children := parent.GetChildren()
	for i, child := range children {
		switch fc := child.(type) {
		case *components.FormComponent[finance_accounts.AccountingPreferences]:
			next := *fc
			next.ChildrenInput = append(slices.Clone(fc.ChildrenInput), productPreferencesFormInputs()...)
			children[i] = &next
		case components.FormComponent[finance_accounts.AccountingPreferences]:
			next := fc
			next.ChildrenInput = append(slices.Clone(fc.ChildrenInput), productPreferencesFormInputs()...)
			children[i] = next
		default:
			if mp, ok := child.(components.MutableParentInterface); ok {
				appendProductPrefsToForm(mp)
			}
		}
	}
	parent.SetChildren(children)
}
