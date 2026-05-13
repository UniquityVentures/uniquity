package p_uniquity_invoices

import (
	"time"

	"github.com/UniquityVentures/lamu/fields"
	"github.com/UniquityVentures/lamu/lamu"
	acct "github.com/UniquityVentures/uniquity/plugins/p_uniquity_accounting"
	currencies "github.com/UniquityVentures/uniquity/plugins/p_uniquity_currencies"
	ent "github.com/UniquityVentures/uniquity/plugins/p_uniquity_entities"
	"gorm.io/gorm"
)

// Invoice type values (customer invoice, vendor bill, credit notes).
const (
	InvoiceTypeOutInvoice = "out_invoice"
	InvoiceTypeOutRefund  = "out_refund"
	InvoiceTypeInInvoice  = "in_invoice"
	InvoiceTypeInRefund   = "in_refund"
)

// Invoice workflow states.
const (
	InvoiceStateDraft     = "draft"
	InvoiceStatePosted    = "posted"
	InvoiceStatePaid      = "paid"
	InvoiceStateCancelled = "cancelled"
)

func init() {
	lamu.RegistryAdmin.Register("p_uniquity_invoices.Invoice", lamu.AdminPanel[Invoice]{
		SearchField: "Reference",
		ListFields:  []string{"Number", "InvoiceType", "State", "InvoiceDate", "AmountTotal", "AmountResidual", "UpdatedAt"},
	})
	lamu.RegistryAdmin.Register("p_uniquity_invoices.Contact", lamu.AdminPanel[Contact]{
		SearchField: "Name",
		ListFields:  []string{"Name", "UpdatedAt"},
	})
	lamu.RegistryAdmin.Register("p_uniquity_invoices.InvoiceLine", lamu.AdminPanel[InvoiceLine]{
		SearchField: "Label",
		ListFields:  []string{"Label", "Quantity", "PriceUnit", "PriceSubtotal", "UpdatedAt"},
	})
	lamu.RegistryAdmin.Register("p_uniquity_invoices.PaymentTerm", lamu.AdminPanel[PaymentTerm]{
		SearchField: "Name",
		ListFields:  []string{"Name", "UpdatedAt"},
	})
}

// Invoice is a subledger document: customer invoices, vendor bills, and credit/debit notes.
type Invoice struct {
	gorm.Model

	EntityID uint       `gorm:"not null;index"`
	Entity   ent.Entity `gorm:"constraint:OnDelete:CASCADE"`

	Number *string `gorm:"size:50;index"`

	PartnerID uint    `gorm:"not null;index"`
	Partner   Contact `gorm:"constraint:OnDelete:RESTRICT"`

	JournalID uint         `gorm:"not null;index"`
	Journal   acct.Journal `gorm:"constraint:OnDelete:RESTRICT"`

	InvoiceType string `gorm:"size:20;not null;index"`
	State       string `gorm:"size:20;not null;default:draft"`

	Reference string `gorm:"size:100"`

	InvoiceDate time.Time `gorm:"type:date;not null"`

	PaymentTermID *uint        `gorm:"index"`
	PaymentTerm   *PaymentTerm `gorm:"constraint:OnDelete:SET NULL"`

	DueDate *time.Time `gorm:"type:date"`

	CurrencyID uint                 `gorm:"not null;index"`
	Currency   currencies.Currency `gorm:"constraint:OnDelete:RESTRICT"`

	// Move is the posted accounting entry (one invoice ↔ at most one journal entry).
	MoveID *uint               `gorm:"uniqueIndex"`
	Move   *acct.JournalEntry `gorm:"constraint:OnDelete:SET NULL"`

	AmountUntaxed  fields.DecimalSix `gorm:"type:numeric(20,2);not null;default:0"`
	AmountTax      fields.DecimalSix `gorm:"type:numeric(20,2);not null;default:0"`
	AmountTotal    fields.DecimalSix `gorm:"type:numeric(20,2);not null;default:0"`
	AmountResidual fields.DecimalSix `gorm:"type:numeric(20,2);not null;default:0"`

	Lines []InvoiceLine `gorm:"foreignKey:InvoiceID"`
}
