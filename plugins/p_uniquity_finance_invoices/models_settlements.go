package p_uniquity_finance_invoices

import "gorm.io/gorm"

// PartiallyPaidInvoice records a payment that does not fully settle the posted invoice total.
type PartiallyPaidInvoice struct {
	gorm.Model

	PaymentID uint    `gorm:"column:payment_id;not null"`
	Payment   Payment `gorm:"foreignKey:PaymentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	PostedInvoiceID uint          `gorm:"column:posted_invoice_id;not null"`
	PostedInvoice   PostedInvoice `gorm:"foreignKey:PostedInvoiceID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	PriorPartiallyPaidInvoiceID *uint               `gorm:"column:prior_partially_paid_invoice_id"`
	PriorPartial                *PartiallyPaidInvoice `gorm:"foreignKey:PriorPartiallyPaidInvoiceID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// PaidInvoice records that the posted invoice is fully paid (typically via the final payment).
type PaidInvoice struct {
	gorm.Model

	PaymentID uint    `gorm:"column:payment_id;not null"`
	Payment   Payment `gorm:"foreignKey:PaymentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	PostedInvoiceID uint          `gorm:"column:posted_invoice_id;not null"`
	PostedInvoice   PostedInvoice `gorm:"foreignKey:PostedInvoiceID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	PriorPartiallyPaidInvoiceID *uint               `gorm:"column:prior_partially_paid_invoice_id"`
	PriorPartial                *PartiallyPaidInvoice `gorm:"foreignKey:PriorPartiallyPaidInvoiceID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}
