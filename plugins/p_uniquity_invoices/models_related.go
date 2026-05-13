package p_uniquity_invoices

import (
	"gorm.io/gorm"
)

// Contact is an invoice partner / counterparty (person/supplier/customer).
type Contact struct {
	gorm.Model

	Name string `gorm:"size:255"`
}

// PaymentTerm names standard settlement rules (e.g. Net 30).
type PaymentTerm struct {
	gorm.Model

	Name string `gorm:"size:255"`
}
