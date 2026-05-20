package p_uniquity_finance_invoices

import (
	"time"

	"github.com/UniquityVentures/lamu/fields"
	finance_accounts "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_accounts"
	finance_creditnotes "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_creditnotes"
	finance_customer "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_customer"
	finance_products "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_products"
	finance_taxes "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_taxes"
	"gorm.io/gorm"
)

// DraftInvoice is an editable sales invoice before posting.
type DraftInvoice struct {
	gorm.Model

	Number *string `gorm:"column:number"`

	Datetime time.Time `gorm:"column:datetime;not null"`

	CustomerID uint                      `gorm:"column:customer_id;not null"`
	Customer   finance_customer.Customer `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	PaymentTermType string      `gorm:"column:payment_term_type;not null"`
	PaymentTermID   uint        `gorm:"column:payment_term_id;not null"`
	PaymentTerm     PaymentTerm `gorm:"foreignKey:PaymentTermID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Taxes []finance_taxes.Tax `gorm:"many2many:draft_invoice_taxes;"`
	Lines []DraftInvoiceLine  `gorm:"foreignKey:DraftInvoiceID"`

	PendingLines []DraftLinePending `gorm:"-"`
}

// DraftInvoiceLine is one draft line with optional per-line taxes (M2M).
type DraftInvoiceLine struct {
	gorm.Model

	DraftInvoiceID uint `gorm:"column:draft_invoice_id;not null"`
	DraftInvoice   DraftInvoice

	ProductID uint                     `gorm:"column:product_id;not null"`
	Product   finance_products.Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Rate     fields.DecimalSix `gorm:"column:rate;type:numeric(19,6);not null"`
	Quantity fields.DecimalSix `gorm:"column:quantity;type:numeric(19,6);not null"`

	Taxes []finance_taxes.Tax `gorm:"many2many:draft_invoice_line_taxes;"`
}

// PostedInvoice is an immutable posted document linked to one draft and one journal entry.
type PostedInvoice struct {
	gorm.Model

	DraftInvoiceID uint         `gorm:"column:draft_invoice_id;not null"`
	DraftInvoice   DraftInvoice `gorm:"foreignKey:DraftInvoiceID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	PostedAt *time.Time `gorm:"column:posted_at"`

	Number string `gorm:"column:number;not null"`

	AccountReceivableID uint                     `gorm:"column:account_receivable_id;not null"`
	AccountReceivable   finance_accounts.Account `gorm:"foreignKey:AccountReceivableID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AccountRevenueID    uint                     `gorm:"column:account_revenue_id;not null"`
	AccountRevenue      finance_accounts.Account `gorm:"foreignKey:AccountRevenueID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AccountTaxPayableID uint                     `gorm:"column:account_tax_payable_id;not null"`
	AccountTaxPayable   finance_accounts.Account `gorm:"foreignKey:AccountTaxPayableID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	JournalID           uint                     `gorm:"column:journal_id;not null"`
	Journal             finance_accounts.Journal `gorm:"foreignKey:JournalID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Datetime time.Time `gorm:"column:datetime;not null"`

	CustomerID uint                      `gorm:"column:customer_id;not null"`
	Customer   finance_customer.Customer `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	PaymentTermType string      `gorm:"column:payment_term_type;not null"`
	PaymentTermID   uint        `gorm:"column:payment_term_id;not null"`
	PaymentTerm     PaymentTerm `gorm:"foreignKey:PaymentTermID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	JournalEntryID uint                          `gorm:"column:journal_entry_id;not null"`
	JournalEntry   finance_accounts.JournalEntry `gorm:"foreignKey:JournalEntryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Taxes []finance_taxes.Tax `gorm:"many2many:posted_invoice_taxes;"`
	Lines []PostedInvoiceLine `gorm:"foreignKey:PostedInvoiceID"`
}

// PostedInvoiceLine is a posted line referencing the revenue journal item for that line.
type PostedInvoiceLine struct {
	gorm.Model

	PostedInvoiceID uint `gorm:"column:posted_invoice_id;not null"`
	PostedInvoice   PostedInvoice

	ProductID uint                     `gorm:"column:product_id;not null"`
	Product   finance_products.Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Rate     fields.DecimalSix `gorm:"column:rate;type:numeric(19,6);not null"`
	Quantity fields.DecimalSix `gorm:"column:quantity;type:numeric(19,6);not null"`

	JournalEntryItemID uint                              `gorm:"column:journal_entry_item_id;not null"`
	JournalEntryItem   finance_accounts.JournalEntryItem `gorm:"foreignKey:JournalEntryItemID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Taxes []finance_taxes.Tax `gorm:"many2many:posted_invoice_line_taxes;"`
}

// CancelledInvoice snapshots a posted invoice that was reversed via credit note.
type CancelledInvoice struct {
	gorm.Model

	PostedInvoiceID uint          `gorm:"column:posted_invoice_id;not null"`
	PostedInvoice   PostedInvoice `gorm:"foreignKey:PostedInvoiceID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	PostedAt    *time.Time `gorm:"column:posted_at"`
	CancelledAt *time.Time `gorm:"column:cancelled_at"`

	Number string `gorm:"column:number;not null"`

	AccountReceivableID uint                     `gorm:"column:account_receivable_id;not null"`
	AccountReceivable   finance_accounts.Account `gorm:"foreignKey:AccountReceivableID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AccountRevenueID    uint                     `gorm:"column:account_revenue_id;not null"`
	AccountRevenue      finance_accounts.Account `gorm:"foreignKey:AccountRevenueID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AccountTaxPayableID uint                     `gorm:"column:account_tax_payable_id;not null"`
	AccountTaxPayable   finance_accounts.Account `gorm:"foreignKey:AccountTaxPayableID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	JournalID           uint                     `gorm:"column:journal_id;not null"`
	Journal             finance_accounts.Journal `gorm:"foreignKey:JournalID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Datetime time.Time `gorm:"column:datetime;not null"`

	CustomerID uint                      `gorm:"column:customer_id;not null"`
	Customer   finance_customer.Customer `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	PaymentTermType string      `gorm:"column:payment_term_type;not null"`
	PaymentTermID   uint        `gorm:"column:payment_term_id;not null"`
	PaymentTerm     PaymentTerm `gorm:"foreignKey:PaymentTermID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	CreditNoteID uint                           `gorm:"column:credit_note_id;not null"`
	CreditNote   finance_creditnotes.CreditNote `gorm:"foreignKey:CreditNoteID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Taxes []finance_taxes.Tax    `gorm:"many2many:cancelled_invoice_taxes;"`
	Lines []CancelledInvoiceLine `gorm:"foreignKey:CancelledInvoiceID"`
}

// CancelledInvoiceLine references a line on the reversing journal entry.
type CancelledInvoiceLine struct {
	gorm.Model

	CancelledInvoiceID uint `gorm:"column:cancelled_invoice_id;not null"`
	CancelledInvoice   CancelledInvoice

	ProductID uint                     `gorm:"column:product_id;not null"`
	Product   finance_products.Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Rate     fields.DecimalSix `gorm:"column:rate;type:numeric(19,6);not null"`
	Quantity fields.DecimalSix `gorm:"column:quantity;type:numeric(19,6);not null"`

	JournalEntryItemID uint                              `gorm:"column:journal_entry_item_id;not null"`
	JournalEntryItem   finance_accounts.JournalEntryItem `gorm:"foreignKey:JournalEntryItemID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Taxes []finance_taxes.Tax `gorm:"many2many:cancelled_invoice_line_taxes;"`
}

// DraftLinePending is submitted JSON for new draft lines before persistence.
type DraftLinePending struct {
	ProductID    uint   `json:"product_id"`
	Rate         string `json:"rate"`
	Quantity     string `json:"quantity"`
	ProductLabel string `json:"product_label,omitempty"`
	FkSlot       string `json:"fk_slot,omitempty"`
	TaxIDs       []uint `json:"tax_ids,omitempty"`
}
