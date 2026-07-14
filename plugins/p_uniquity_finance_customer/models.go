package p_uniquity_finance_customer

import (
	"github.com/lariv-in/lago"
	"gorm.io/gorm"
)

// Customer is a business customer record (contact and India tax identifiers).
type Customer struct {
	gorm.Model

	Name    string `gorm:"not null"`
	Address string
	GSTIN   string `gorm:"column:gstin"`
	PAN     string `gorm:"column:pan"`
	Phone   string
	Email   string
	Website string
}

func init() {
	lago.RegistryAdmin.Register("p_uniquity_finance_customer.Customer", lago.AdminPanel[Customer]{
		SearchField: "Name",
		ListFields:  []string{"Name", "Email", "Phone", "GSTIN", "PAN", "UpdatedAt"},
	})
}
