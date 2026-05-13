package p_uniquity_products

import (
	"github.com/UniquityVentures/lamu/lamu"
	ent "github.com/UniquityVentures/uniquity/plugins/p_uniquity_entities"
	"gorm.io/gorm"
)

func init() {
	lamu.RegistryAdmin.Register("p_uniquity_products.Product", lamu.AdminPanel[Product]{
		SearchField: "Name",
		ListFields:  []string{"Name", "Code", "UpdatedAt"},
	})
}

// Product is a minimal catalog item for invoice lines (expand as needed).
type Product struct {
	gorm.Model

	EntityID uint       `gorm:"not null;index"`
	Entity   ent.Entity `gorm:"constraint:OnDelete:CASCADE"`

	Name string `gorm:"size:255;not null"`
	Code string `gorm:"size:64"`
}
