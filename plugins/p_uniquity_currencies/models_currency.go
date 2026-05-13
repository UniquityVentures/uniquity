package p_uniquity_currencies

import (
	"github.com/UniquityVentures/lamu/lamu"
	"gorm.io/gorm"
)

func init() {
	lamu.RegistryAdmin.Register("p_uniquity_currencies.Currency", lamu.AdminPanel[Currency]{
		SearchField: "Code",
		ListFields:  []string{"Code", "Name", "Symbol", "IsActive", "UpdatedAt"},
	})
}

// Currency is an ISO 4217 currency record.
type Currency struct {
	gorm.Model

	Code     string `gorm:"size:3;not null;uniqueIndex"`
	Name     string `gorm:"size:50;not null"`
	Symbol   string `gorm:"size:10"`
	IsActive bool   `gorm:"not null;default:true"`
}
