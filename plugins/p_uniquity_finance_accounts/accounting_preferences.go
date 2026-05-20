package p_uniquity_finance_accounts

import (
	"log/slog"

	"github.com/UniquityVentures/lamu/lamu"
	"gorm.io/gorm"
)

// AccountingPreferences is the singleton row for finance-wide accounting settings (one row, typically id = 1).
type AccountingPreferences struct {
	gorm.Model

	// InvoiceNumberFormat is a template for suggested/autofilled invoice numbers (interpretation is up to callers).
	InvoiceNumberFormat string `gorm:"column:invoice_number_format"`

	// InvoicePDFTemplate is Go text/template source whose output must be valid Typst markup; executed with getters.MapFromStruct on invoice rows (same shape as each detail page’s $in). Template funcs: num2words, num2wordsAnd, num2wordsRupees, invoiceGrandTotalWords (github.com/divan/num2words). Compiled to PDF via Typst (gotypst).
	InvoicePDFTemplate string `gorm:"column:invoice_pdf_template"`

	// DefaultJournalID optionally prefills the journal on new draft invoices (see finance invoices create form).
	DefaultJournalID *uint    `gorm:"column:default_journal_id"`
	DefaultJournal   *Journal `gorm:"foreignKey:DefaultJournalID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

// LoadAccountingPreferences returns the singleton preferences row, creating id 1 if missing (same idea as p_otp OTPPreferences).
func LoadAccountingPreferences(db *gorm.DB) AccountingPreferences {
	var prefs AccountingPreferences
	if err := db.FirstOrCreate(&prefs, AccountingPreferences{Model: gorm.Model{ID: 1}}).Error; err != nil {
		slog.Warn("LoadAccountingPreferences", "error", err)
	}
	return prefs
}

func init() {
	lamu.RegistryAdmin.Register("p_uniquity_finance_accounts.AccountingPreferences", lamu.AdminPanel[AccountingPreferences]{
		SearchField: "",
	})
}
