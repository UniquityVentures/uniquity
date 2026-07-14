package p_uniquity_finance_taxes

import (
	"fmt"

	finance_accounts "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_accounts"
	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/fields"
	"gorm.io/gorm"
)

// Tax is a named tax rate (percentage with six decimal places).
type Tax struct {
	gorm.Model

	Name       string            `gorm:"not null"`
	Percentage fields.DecimalSix `gorm:"type:numeric(19,6);not null"`
	TaxType    TaxKind           `gorm:"column:tax_type;type:\"TaxKind\";not null"`

	AccountID *uint                     `gorm:"column:account_id"`
	Account   *finance_accounts.Account `gorm:"foreignKey:AccountID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (t *Tax) BeforeCreate(_ *gorm.DB) error {
	t.Percentage = t.Percentage.NormalizeDecimals()
	if t.TaxType != TaxKindLevied && t.TaxType != TaxKindWithholding {
		return fmt.Errorf("tax type is required")
	}
	if t.TaxType == TaxKindWithholding && (t.AccountID == nil || *t.AccountID == 0) {
		return fmt.Errorf("withholding tax requires a ledger account")
	}
	return nil
}

func (t *Tax) BeforeUpdate(_ *gorm.DB) error {
	t.Percentage = t.Percentage.NormalizeDecimals()
	if t.TaxType == TaxKindWithholding && (t.AccountID == nil || *t.AccountID == 0) {
		return fmt.Errorf("withholding tax requires a ledger account")
	}
	return nil
}

func init() {
	lago.RegistryAdmin.Register("p_uniquity_finance_taxes.Tax", lago.AdminPanel[Tax]{
		SearchField: "Name",
		ListFields:  []string{"Name", "Percentage", "TaxType", "AccountID", "UpdatedAt"},
		Preload:     []string{"Account"},
	})
}
