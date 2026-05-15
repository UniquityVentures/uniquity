package p_uniquity_finance_customer

import (
	"github.com/UniquityVentures/lamu/lamu"
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
	lamu.RegistryAdmin.Register("p_uniquity_finance_customer.Customer", lamu.AdminPanel[Customer]{
		SearchField: "Name",
		ListFields:  []string{"Name", "Email", "Phone", "GSTIN", "PAN", "UpdatedAt"},
	})
}
