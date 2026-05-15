package p_uniquity_finance_products

import (
	"github.com/UniquityVentures/lamu/fields"
	"github.com/UniquityVentures/lamu/lamu"
	finance_taxes "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_taxes"
	"gorm.io/gorm"
)

// Product is a sellable item with optional tax links, cost/price, and HSN (India) code.
type Product struct {
	gorm.Model

	Name       string `gorm:"not null"`
	Taxes      []finance_taxes.Tax `gorm:"many2many:product_taxes;"`
	BaseCost   fields.DecimalSix   `gorm:"column:base_cost;type:numeric(19,6);not null"`
	SalesPrice fields.DecimalSix   `gorm:"column:sales_price;type:numeric(19,6);not null"`
	HSNCode    int64               `gorm:"column:hsn_code;not null"`
}

func (p *Product) BeforeCreate(_ *gorm.DB) error {
	p.BaseCost = p.BaseCost.NormalizeDecimals()
	p.SalesPrice = p.SalesPrice.NormalizeDecimals()
	return nil
}

func (p *Product) BeforeUpdate(_ *gorm.DB) error {
	p.BaseCost = p.BaseCost.NormalizeDecimals()
	p.SalesPrice = p.SalesPrice.NormalizeDecimals()
	return nil
}

func init() {
	lamu.RegistryAdmin.Register("p_uniquity_finance_products.Product", lamu.AdminPanel[Product]{
		SearchField: "Name",
		ListFields:  []string{"Name", "BaseCost", "SalesPrice", "HSNCode", "UpdatedAt"},
	})
}
