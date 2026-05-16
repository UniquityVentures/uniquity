package p_uniquity_finance_invoices

import (
	"time"

	"gorm.io/gorm"
)

// PaymentTerm is the umbrella row for a polymorphic payment term (backing row identified by Type + BackingID).
type PaymentTerm struct {
	gorm.Model

	Type      string `gorm:"column:type;not null"`
	BackingID uint   `gorm:"column:backing_id;not null"`
}

// PaymentTermDueDate backs a calendar due date term.
type PaymentTermDueDate struct {
	gorm.Model

	Datetime time.Time `gorm:"column:datetime;not null"`
}

// PaymentTermRelative backs a duration-based term.
type PaymentTermRelative struct {
	gorm.Model

	Duration time.Duration `gorm:"column:duration;not null"`
}
