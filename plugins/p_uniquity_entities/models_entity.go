package p_uniquity_entities

import (
	currencies "github.com/UniquityVentures/uniquity/plugins/p_uniquity_currencies"
	"github.com/UniquityVentures/lamu/lamu"
	"gorm.io/gorm"
)

func init() {
	lamu.RegistryAdmin.Register("p_uniquity_entities.Entity", lamu.AdminPanel[Entity]{
		SearchField: "Name",
		ListFields:  []string{"Name", "Slug", "TaxID", "Email", "Phone", "UpdatedAt"},
	})
}

// Entity is a legal entity (company) for accounting and operations.
type Entity struct {
	gorm.Model

	Name    string `gorm:"size:255;not null"`
	Slug    string `gorm:"size:255"`
	TaxID   string `gorm:"size:50"`
	Address string  `gorm:"type:text"`
	Phone   string  `gorm:"size:50"`
	Mobile1 string  `gorm:"size:50"`
	Mobile2 string  `gorm:"size:50"`
	Email   string  `gorm:"size:255"`
	Website string  `gorm:"size:512"`

	CurrencyID uint                `gorm:"not null;index"`
	Currency   currencies.Currency `gorm:"constraint:OnDelete:RESTRICT"`

	// LogoPath stores an opaque filesystem blob key or path (see p_filesystem).
	LogoPath string `gorm:"size:512"`
}
