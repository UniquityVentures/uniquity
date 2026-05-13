package p_uniquity_accounting

import (
	"github.com/UniquityVentures/lamu/fields"
	"gorm.io/gorm"
)

type Journal struct {
	gorm.Model

	Name string `gorm:"size:255;not null"`
}

type JournalEntry struct {
	gorm.Model

	JournalID uint    `gorm:"not null;index"`
	Journal   Journal `gorm:"not null"`
}

type JournalEntryItem struct {
	gorm.Model

	Amount    fields.DecimalSix `gorm:"type:numeric(19,6);not null"`
	AccountID uint    `gorm:"not null"`
	Account   Account `gorm:"not null"`

	JournalEntryID uint         `gorm:"not null;index"`
	JournalEntry   JournalEntry `gorm:"not null"`
}
