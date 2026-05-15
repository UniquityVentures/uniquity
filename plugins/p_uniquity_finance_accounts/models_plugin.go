package p_uniquity_finance_accounts

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pluginModels() lamu.PluginFeatures[any] {
	return lamu.PluginFeatures[any]{
		Entries: []registry.Pair[string, any]{
			{Key: "p_uniquity_finance_accounts.Account", Value: Account{}},
			{Key: "p_uniquity_finance_accounts.Currency", Value: Currency{}},
			{Key: "p_uniquity_finance_accounts.Journal", Value: Journal{}},
			{Key: "p_uniquity_finance_accounts.SourceDoc", Value: SourceDoc{}},
			{Key: "p_uniquity_finance_accounts.JournalEntry", Value: JournalEntry{}},
			{Key: "p_uniquity_finance_accounts.JournalEntryItem", Value: JournalEntryItem{}},
		},
	}
}
