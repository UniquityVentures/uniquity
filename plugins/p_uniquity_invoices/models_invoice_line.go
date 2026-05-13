package p_uniquity_invoices

import (
	"github.com/UniquityVentures/lamu/fields"
	acct "github.com/UniquityVentures/uniquity/plugins/p_uniquity_accounting"
	prod "github.com/UniquityVentures/uniquity/plugins/p_uniquity_products"
	tax "github.com/UniquityVentures/uniquity/plugins/p_uniquity_tax_rates"
	"gorm.io/gorm"
)

// InvoiceLine is a single line on an invoice.
type InvoiceLine struct {
	gorm.Model

	InvoiceID uint    `gorm:"not null;index"`
	Invoice   Invoice `gorm:"constraint:OnDelete:CASCADE"`

	ProductID *uint         `gorm:"index"`
	Product   *prod.Product `gorm:"constraint:OnDelete:SET NULL"`

	Label string `gorm:"size:200"`

	Quantity      fields.DecimalSix `gorm:"type:numeric(12,2);not null;default:0"`
	PriceUnit     fields.DecimalSix `gorm:"type:numeric(12,2);not null;default:0"`
	Discount      fields.DecimalSix `gorm:"type:numeric(5,2);not null;default:0"`
	PriceSubtotal fields.DecimalSix `gorm:"type:numeric(16,2);not null;default:0"`

	AccountID uint         `gorm:"not null;index"`
	Account   acct.Account `gorm:"constraint:OnDelete:RESTRICT"`

	Taxes []tax.TaxRate `gorm:"many2many:invoice_line_tax_rates;"`
}
