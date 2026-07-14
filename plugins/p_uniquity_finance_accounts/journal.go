package p_uniquity_finance_accounts

import (
	"time"

	"github.com/lariv-in/lago/fields"
	"gorm.io/gorm"
)

// Journal is a named journal (book) scoped to one currency.
type Journal struct {
	gorm.Model

	Name       string      `gorm:"not null"`
	IsActive   bool        `gorm:"column:is_active;not null"`
	CurrencyID uint        `gorm:"not null"`
	Currency   Currency    `gorm:"foreignKey:CurrencyID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Type       JournalType `gorm:"column:journal_type;type:\"JournalType\";not null"`
}

// JournalEntry is a dated entry linked to a source document and a journal.
type JournalEntry struct {
	gorm.Model

	Datetime time.Time `gorm:"column:datetime;not null"`

	SourceDocID uint      `gorm:"not null"`
	SourceDoc   SourceDoc `gorm:"foreignKey:SourceDocID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	JournalID uint    `gorm:"not null"`
	Journal   Journal `gorm:"foreignKey:JournalID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// JournalEntryItem is one line belonging to a journal entry.
type JournalEntryItem struct {
	gorm.Model

	Datetime time.Time `gorm:"column:datetime;not null"`

	AccountID uint    `gorm:"not null"`
	Account   Account `gorm:"foreignKey:AccountID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Amount fields.DecimalSix `gorm:"type:numeric(19,6);not null"`

	JournalEntryID uint         `gorm:"not null"`
	JournalEntry   JournalEntry `gorm:"foreignKey:JournalEntryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}
