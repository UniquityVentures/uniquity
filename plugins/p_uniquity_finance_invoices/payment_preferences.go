package p_uniquity_finance_invoices

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/UniquityVentures/lamu/lamu"
	finance_accounts "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_accounts"
	finance_products "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_products"
	"gorm.io/gorm"
)

const paymentPreferencesSingletonID uint = 1
const paymentPreferencesTable = "payment_preferences"

// Payment preference form field names shared with the accounting preferences page patch.
const PaymentPrefAccountIDField = "PaymentAccountID"

// PaymentPreferences is the singleton row for payment-wide GL settings (one row, typically id = 1).
type PaymentPreferences struct {
	gorm.Model

	PaymentAccountID *uint                     `gorm:"column:payment_account_id"`
	PaymentAccount   *finance_accounts.Account `gorm:"foreignKey:PaymentAccountID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (PaymentPreferences) TableName() string { return "payment_preferences" }

func paymentPreferencesDB(db *gorm.DB) *gorm.DB {
	// NewDB avoids inheriting the active model/table from Payment create hooks on the same tx.
	return db.Session(&gorm.Session{NewDB: true}).Table(paymentPreferencesTable)
}

// PaymentPreferenceFormFields returns form field names owned by [PaymentPreferences].
func PaymentPreferenceFormFields() map[string]struct{} {
	return map[string]struct{}{
		PaymentPrefAccountIDField: {},
	}
}

// LoadPaymentPreferences returns the singleton preferences row, creating id 1 if missing.
func LoadPaymentPreferences(db *gorm.DB) PaymentPreferences {
	var prefs PaymentPreferences
	err := paymentPreferencesDB(db).Where("id = ?", paymentPreferencesSingletonID).First(&prefs).Error
	if err == nil {
		return prefs
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		slog.Warn("LoadPaymentPreferences", "error", err)
		return prefs
	}
	prefs = PaymentPreferences{Model: gorm.Model{ID: paymentPreferencesSingletonID}}
	if err := paymentPreferencesDB(db).Create(&prefs).Error; err != nil {
		slog.Warn("LoadPaymentPreferences", "error", err)
	}
	return prefs
}

// ValidatePaymentPreferencesForCreate ensures the bank/cash account is configured for payments.
func ValidatePaymentPreferencesForCreate(tx *gorm.DB, prefs *PaymentPreferences) error {
	if prefs == nil {
		return fmt.Errorf("payment preferences required")
	}
	accountID := finance_products.OptionalUintValue(prefs.PaymentAccountID)
	if accountID == 0 {
		return fmt.Errorf("payment account is required in payment preferences")
	}
	return finance_accounts.ValidateLeafAccountBalanceType(tx, accountID, finance_accounts.BalanceTypeDebit, "payment account")
}

func init() {
	lamu.RegistryAdmin.Register("p_uniquity_finance_invoices.PaymentPreferences", lamu.AdminPanel[PaymentPreferences]{
		SearchField: "",
	})
}
