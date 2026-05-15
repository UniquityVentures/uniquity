package p_uniquity_finance_accounts

import (
	"github.com/UniquityVentures/lamu/lamu"
	"gorm.io/gorm"
)

// Account is a chart-of-accounts style row with an optional parent account.
type Account struct {
	gorm.Model

	Name          string      `gorm:"not null"`
	Code          int         `gorm:"not null"`
	IsGroup       bool        `gorm:"column:is_group;not null"`
	BalanceType   BalanceType `gorm:"type:\"BalanceType\";not null"`

	ParentID *uint
	Parent   *Account `gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

// Currency is an ISO 4217 currency (numeric code, name, symbol, minor units / exponent).
type Currency struct {
	gorm.Model

	Code      int    `gorm:"not null;uniqueIndex"` // ISO 4217 numeric code (e.g. 840 for USD)
	Name      string `gorm:"not null"`
	Symbol    string `gorm:"not null"`
	MinorUnit int    `gorm:"column:minor_unit;not null"` // decimal places for the minor currency unit
}

func init() {
	lamu.RegistryAdmin.Register("p_uniquity_finance_accounts.Account", lamu.AdminPanel[Account]{
		SearchField: "Name",
		ListFields:  []string{"Name", "Code", "IsGroup", "BalanceType", "Parent.Name", "UpdatedAt"},
		Preload:     []string{"Parent"},
	})
	lamu.RegistryAdmin.Register("p_uniquity_finance_accounts.Currency", lamu.AdminPanel[Currency]{
		SearchField: "Name",
		ListFields:  []string{"Code", "Name", "Symbol", "MinorUnit", "UpdatedAt"},
	})
	lamu.RegistryAdmin.Register("p_uniquity_finance_accounts.Journal", lamu.AdminPanel[Journal]{
		SearchField: "Name",
		ListFields:  []string{"Name", "IsActive", "Currency.Name", "Type", "UpdatedAt"},
		Preload:     []string{"Currency"},
	})
	lamu.RegistryAdmin.Register("p_uniquity_finance_accounts.JournalEntry", lamu.AdminPanel[JournalEntry]{
		SearchField: "Journal.Name",
		ListFields:  []string{"Datetime", "Journal.Name", "SourceDoc.ID", "UpdatedAt"},
		Preload:     []string{"Journal", "SourceDoc"},
	})
	lamu.RegistryAdmin.Register("p_uniquity_finance_accounts.JournalEntryItem", lamu.AdminPanel[JournalEntryItem]{
		SearchField: "Account.Name",
		ListFields:  []string{"Datetime", "Amount", "Account.Name", "JournalEntry.ID", "UpdatedAt"},
		Preload:     []string{"Account", "JournalEntry"},
	})
}
