package p_uniquity_finance_products

import (
	"github.com/UniquityVentures/lamu/fields"
	"github.com/UniquityVentures/lamu/lamu"
	finance_accounts "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_accounts"
	finance_taxes "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_taxes"
	"gorm.io/gorm"
)

// Product is a sellable item with optional tax links, cost/price, and HSN (India) code.
type Product struct {
	gorm.Model

	Name       string              `gorm:"not null"`
	Taxes      []finance_taxes.Tax `gorm:"many2many:product_taxes;"`
	BaseCost   fields.DecimalSix   `gorm:"column:base_cost;type:numeric(19,6);not null"`
	SalesPrice fields.DecimalSix   `gorm:"column:sales_price;type:numeric(19,6);not null"`
	HSNCode    int64               `gorm:"column:hsn_code;not null"`

	InventoryAccountID *uint                     `gorm:"column:inventory_account_id"`
	InventoryAccount   *finance_accounts.Account `gorm:"foreignKey:InventoryAccountID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CostOfSalesAcctID  *uint                     `gorm:"column:cost_of_sales_account_id"`
	CostOfSalesAccount *finance_accounts.Account `gorm:"foreignKey:CostOfSalesAcctID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	InputTaxAccountID  *uint                     `gorm:"column:input_tax_account_id"`
	InputTaxAccount    *finance_accounts.Account `gorm:"foreignKey:InputTaxAccountID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
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
		ListFields:  []string{"Name", "BaseCost", "SalesPrice", "HSNCode", "InventoryAccountID", "CostOfSalesAcctID", "InputTaxAccountID", "UpdatedAt"},
		Preload:     []string{"InventoryAccount", "CostOfSalesAccount", "InputTaxAccount"},
	})
}
