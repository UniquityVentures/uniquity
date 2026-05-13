package p_uniquity_tax_rates

import (
	"github.com/UniquityVentures/lamu/fields"
	"github.com/UniquityVentures/lamu/lamu"
	acct "github.com/UniquityVentures/uniquity/plugins/p_uniquity_accounting"
	ent "github.com/UniquityVentures/uniquity/plugins/p_uniquity_entities"
	"gorm.io/gorm"
)

const (
	TaxScopeSale     = "sale"
	TaxScopePurchase = "purchase"
	TaxScopeNone     = "none"
)

func init() {
	lamu.RegistryAdmin.Register("p_uniquity_tax_rates.TaxRate", lamu.AdminPanel[TaxRate]{
		SearchField: "Name",
		ListFields:  []string{"Name", "Scope", "Amount", "IsActive", "UpdatedAt"},
	})
}

// TaxRate is a percentage tax (e.g. VAT) scoped to an entity.
type TaxRate struct {
	gorm.Model

	EntityID uint       `gorm:"not null;index"`
	Entity   ent.Entity `gorm:"constraint:OnDelete:CASCADE"`

	Name  string `gorm:"size:100;not null"`
	Scope string `gorm:"size:20;not null"`

	Amount fields.DecimalSix `gorm:"type:numeric(12,4);not null;default:0"`

	AccountID *uint         `gorm:"index"`
	Account   *acct.Account `gorm:"foreignKey:AccountID;constraint:OnDelete:SET NULL"`

	RefundAccountID *uint         `gorm:"index"`
	RefundAccount   *acct.Account `gorm:"foreignKey:RefundAccountID;constraint:OnDelete:SET NULL"`

	IsActive bool `gorm:"not null;default:true"`
}
