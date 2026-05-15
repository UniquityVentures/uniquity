package p_uniquity_finance_invoices

import (
	"fmt"
	"strings"
	"time"

	"github.com/UniquityVentures/lamu/fields"
	"github.com/UniquityVentures/lamu/lamu"
	finance_customer "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_customer"
	finance_products "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_products"
	finance_taxes "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_taxes"
	"gorm.io/gorm"
)

// PaymentTermDueDate is a payment term anchored to a calendar datetime.
type PaymentTermDueDate struct {
	gorm.Model

	Datetime time.Time `gorm:"column:datetime;not null"`
}

// PaymentTermRelative is a payment term expressed as a duration after invoice date (nanoseconds).
type PaymentTermRelative struct {
	gorm.Model

	Duration time.Duration `gorm:"column:duration_ns;not null"`
}

// PaymentTerm is the canonical row invoices reference; it points at a backing due-date or relative row.
// There is no database foreign key to the backing table (polymorphic by Type + BackingID).
type PaymentTerm struct {
	gorm.Model

	Type      string `gorm:"column:payment_term_type;not null"`
	BackingID uint   `gorm:"column:backing_id;not null"`
}

// BeforeDelete removes the backing row so the pointer table and backing stay aligned.
func (p *PaymentTerm) BeforeDelete(tx *gorm.DB) error {
	switch p.Type {
	case PaymentTermTypeDueDate:
		return tx.Delete(&PaymentTermDueDate{}, p.BackingID).Error
	case PaymentTermTypeRelative:
		return tx.Delete(&PaymentTermRelative{}, p.BackingID).Error
	default:
		return fmt.Errorf("p_uniquity_finance_invoices: unknown payment term type %q", p.Type)
	}
}

// Invoice is a customer invoice linked to a payment term row and optional tax links.
type Invoice struct {
	gorm.Model

	Number     string                    `gorm:"not null"`
	Datetime   time.Time                 `gorm:"column:datetime;not null"`
	CustomerID uint                      `gorm:"not null"`
	Customer   finance_customer.Customer `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	PaymentTermID uint        `gorm:"not null"`
	PaymentTerm   PaymentTerm `gorm:"foreignKey:PaymentTermID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Taxes []finance_taxes.Tax `gorm:"many2many:invoice_taxes;"`

	Lines []InvoiceLine `gorm:"foreignKey:InvoiceID"`

	Status InvoiceStatus `gorm:"type:\"InvoiceStatus\";not null"`

	// PendingLines is set from the create-invoice form (JSON) and consumed by AfterCreate; not persisted on invoices.
	PendingLines []InvoiceLinePending `gorm:"-"`
}

// InvoiceLinePending is one line in the invoice create form JSON (product + quantity text).
type InvoiceLinePending struct {
	ProductID uint   `json:"product_id"`
	Quantity  string `json:"quantity"`
}

// AfterCreate inserts [InvoiceLine] rows from [Invoice.PendingLines] inside the same transaction as the invoice.
func (i *Invoice) AfterCreate(tx *gorm.DB) error {
	if len(i.PendingLines) == 0 {
		return nil
	}
	for _, p := range i.PendingLines {
		if p.ProductID == 0 {
			return fmt.Errorf("invoice line: product is required")
		}
		var qty fields.DecimalSix
		if err := qty.UnmarshalText([]byte(strings.TrimSpace(p.Quantity))); err != nil {
			return fmt.Errorf("invoice line product #%d quantity: %w", p.ProductID, err)
		}
		qty = qty.NormalizeDecimals()
		if qty.R == nil || qty.R.Sign() <= 0 {
			return fmt.Errorf("invoice line product #%d: quantity must be positive", p.ProductID)
		}
		line := InvoiceLine{
			InvoiceID: i.ID,
			ProductID: p.ProductID,
			Quantity:  qty,
		}
		if err := tx.Create(&line).Error; err != nil {
			return err
		}
	}
	i.PendingLines = nil
	return nil
}

// InvoiceLine is one line on an invoice.
type InvoiceLine struct {
	gorm.Model

	InvoiceID uint    `gorm:"not null"`
	Invoice   Invoice `gorm:"foreignKey:InvoiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	ProductID uint                     `gorm:"not null"`
	Product   finance_products.Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Quantity fields.DecimalSix `gorm:"type:numeric(19,6);not null"`

	Taxes []finance_taxes.Tax `gorm:"many2many:invoice_line_taxes;"`
}

func (l *InvoiceLine) BeforeCreate(_ *gorm.DB) error {
	l.Quantity = l.Quantity.NormalizeDecimals()
	return nil
}

func (l *InvoiceLine) BeforeUpdate(_ *gorm.DB) error {
	l.Quantity = l.Quantity.NormalizeDecimals()
	return nil
}

func init() {
	lamu.RegistryAdmin.Register("p_uniquity_finance_invoices.PaymentTermDueDate", lamu.AdminPanel[PaymentTermDueDate]{
		SearchField: "ID",
		ListFields:  []string{"Datetime", "UpdatedAt"},
	})
	lamu.RegistryAdmin.Register("p_uniquity_finance_invoices.PaymentTermRelative", lamu.AdminPanel[PaymentTermRelative]{
		SearchField: "ID",
		ListFields:  []string{"Duration", "UpdatedAt"},
	})
	lamu.RegistryAdmin.Register("p_uniquity_finance_invoices.PaymentTerm", lamu.AdminPanel[PaymentTerm]{
		SearchField: "Type",
		ListFields:  []string{"Type", "UpdatedAt"},
	})
	lamu.RegistryAdmin.Register("p_uniquity_finance_invoices.Invoice", lamu.AdminPanel[Invoice]{
		SearchField: "Number",
		ListFields: []string{
			"Number", "Datetime", "Customer.Name", "PaymentTermID", "Status", "UpdatedAt",
		},
		Preload: []string{"Customer", "Taxes", "PaymentTerm", "Lines", "Lines.Product"},
	})
	lamu.RegistryAdmin.Register("p_uniquity_finance_invoices.InvoiceLine", lamu.AdminPanel[InvoiceLine]{
		SearchField: "Invoice.Number",
		ListFields:  []string{"Invoice.Number", "Product.Name", "Quantity", "UpdatedAt"},
		Preload:     []string{"Invoice", "Product", "Taxes"},
	})
}
