package p_uniquity_finance_taxes

import (
	"github.com/UniquityVentures/lamu/fields"
	"github.com/UniquityVentures/lamu/lamu"
	"gorm.io/gorm"
)

// Tax is a named tax rate (percentage with six decimal places).
type Tax struct {
	gorm.Model

	Name       string            `gorm:"not null"`
	Percentage fields.DecimalSix `gorm:"type:numeric(19,6);not null"`
}

func (t *Tax) BeforeCreate(_ *gorm.DB) error {
	t.Percentage = t.Percentage.NormalizeDecimals()
	return nil
}

func (t *Tax) BeforeUpdate(_ *gorm.DB) error {
	t.Percentage = t.Percentage.NormalizeDecimals()
	return nil
}

func init() {
	lamu.RegistryAdmin.Register("p_uniquity_finance_taxes.Tax", lamu.AdminPanel[Tax]{
		SearchField: "Name",
		ListFields:  []string{"Name", "Percentage", "UpdatedAt"},
	})
}
