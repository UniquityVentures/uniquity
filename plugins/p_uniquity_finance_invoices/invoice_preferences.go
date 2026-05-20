package p_uniquity_finance_invoices

import (
	"fmt"
	"log/slog"

	"github.com/UniquityVentures/lamu/lamu"
	finance_accounts "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_accounts"
	finance_products "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_products"
	"gorm.io/gorm"
)

// Invoice preference form field names shared with the accounting preferences page patch.
const (
	InvoicePrefAccountReceivableIDField = "AccountReceivableID"
	InvoicePrefAccountRevenueIDField    = "AccountRevenueID"
	InvoicePrefAccountTaxPayableIDField = "AccountTaxPayableID"
	InvoicePrefJournalIDField           = "JournalID"
)

// InvoicePreferences is the singleton row for invoice-wide GL and journal settings (one row, typically id = 1).
type InvoicePreferences struct {
	gorm.Model

	AccountReceivableID *uint                     `gorm:"column:account_receivable_id"`
	AccountReceivable   *finance_accounts.Account `gorm:"foreignKey:AccountReceivableID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AccountRevenueID    *uint                     `gorm:"column:account_revenue_id"`
	AccountRevenue      *finance_accounts.Account `gorm:"foreignKey:AccountRevenueID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AccountTaxPayableID *uint                     `gorm:"column:account_tax_payable_id"`
	AccountTaxPayable   *finance_accounts.Account `gorm:"foreignKey:AccountTaxPayableID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	JournalID           *uint                     `gorm:"column:journal_id"`
	Journal             *finance_accounts.Journal `gorm:"foreignKey:JournalID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// InvoicePreferenceFormFields returns form field names owned by [InvoicePreferences].
func InvoicePreferenceFormFields() map[string]struct{} {
	return map[string]struct{}{
		InvoicePrefAccountReceivableIDField: {},
		InvoicePrefAccountRevenueIDField:    {},
		InvoicePrefAccountTaxPayableIDField: {},
		InvoicePrefJournalIDField:           {},
	}
}

// LoadInvoicePreferences returns the singleton preferences row, creating id 1 if missing.
func LoadInvoicePreferences(db *gorm.DB) InvoicePreferences {
	var prefs InvoicePreferences
	if err := db.FirstOrCreate(&prefs, InvoicePreferences{Model: gorm.Model{ID: 1}}).Error; err != nil {
		slog.Warn("LoadInvoicePreferences", "error", err)
	}
	return prefs
}

// ValidateInvoicePreferencesForPosting ensures AR, revenue, and journal are configured for posting.
func ValidateInvoicePreferencesForPosting(tx *gorm.DB, prefs *InvoicePreferences) error {
	if prefs == nil {
		return fmt.Errorf("invoice preferences required")
	}
	if err := finance_accounts.ValidateLeafAccountBalanceType(tx, finance_products.OptionalUintValue(prefs.AccountReceivableID), finance_accounts.BalanceTypeDebit, "accounts receivable"); err != nil {
		return err
	}
	if err := finance_accounts.ValidateLeafAccountBalanceType(tx, finance_products.OptionalUintValue(prefs.AccountRevenueID), finance_accounts.BalanceTypeCredit, "revenue account"); err != nil {
		return err
	}
	if err := finance_accounts.ValidateLeafAccountBalanceType(tx, finance_products.OptionalUintValue(prefs.AccountTaxPayableID), finance_accounts.BalanceTypeCredit, "tax payable account"); err != nil {
		return err
	}
	if finance_products.OptionalUintValue(prefs.JournalID) == 0 {
		return fmt.Errorf("journal is required in invoice preferences")
	}
	return nil
}

func init() {
	lamu.RegistryAdmin.Register("p_uniquity_finance_invoices.InvoicePreferences", lamu.AdminPanel[InvoicePreferences]{
		SearchField: "",
	})
}
