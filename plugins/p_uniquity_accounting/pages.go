package p_uniquity_accounting

import (
	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pluginPages() lamu.PluginFeatures[components.PageInterface] {
	entries := []registry.Pair[string, components.PageInterface]{
		{Key: "accounting.MainMenu", Value: &components.SidebarMenu{
			Title: getters.Static("Accounting"),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("Back to Home"),
				Url:   lamu.RoutePath("dashboard.AppsPage", nil),
			},
			Children: []components.PageInterface{},
		}},
	}
	entries = append(entries, pageEntriesAccountPages()...)
	entries = append(entries, pageEntriesTransactionPages()...)
	return lamu.PluginFeatures[components.PageInterface]{
		Entries: entries,
		Patches: []registry.Pair[string, func(components.PageInterface) components.PageInterface]{
			{Key: "accounting.MainMenu", Value: patchAccountingMainMenuAccounts},
			{Key: "accounting.MainMenu", Value: patchAccountingMainMenuTransactions},
		},
	}
}
