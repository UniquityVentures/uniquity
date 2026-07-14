package p_uniquity_finance_creditnotes

import (
	"fmt"
	"math/big"
	"time"

	finance_accounts "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_accounts"
	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/fields"
	"gorm.io/gorm"
)

// CreditNoteSourceDocType is stored on [finance_accounts.SourceDoc.Type] for reversal entries tied to a credit note.
const CreditNoteSourceDocType = "p_uniquity_finance_creditnotes.CreditNote"

// CreditNote records a credit against an existing journal entry and the reversing entry created automatically.
type CreditNote struct {
	gorm.Model

	Datetime time.Time `gorm:"column:datetime;not null"`
	Reason   string    `gorm:"type:text"`

	JournalEntryID uint                          `gorm:"column:journal_entry_id;not null"`
	JournalEntry   finance_accounts.JournalEntry `gorm:"foreignKey:JournalEntryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	ReversedJournalEntryID uint                          `gorm:"column:reversed_journal_entry_id;not null"`
	ReversedJournalEntry   finance_accounts.JournalEntry `gorm:"foreignKey:ReversedJournalEntryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	pendingReversalSourceDocID uint `gorm:"-"`
}

// BeforeCreate inserts a reversing [finance_accounts.JournalEntry] with negated line amounts for the selected journal entry.
func (c *CreditNote) BeforeCreate(tx *gorm.DB) error {
	if c.JournalEntryID == 0 {
		return fmt.Errorf("journal entry is required")
	}
	if c.ReversedJournalEntryID != 0 {
		return fmt.Errorf("reversed journal entry must not be set manually")
	}
	var orig finance_accounts.JournalEntry
	if err := tx.First(&orig, c.JournalEntryID).Error; err != nil {
		return fmt.Errorf("load journal entry: %w", err)
	}
	var items []finance_accounts.JournalEntryItem
	if err := tx.Where("journal_entry_id = ?", orig.ID).Order("id ASC").Find(&items).Error; err != nil {
		return fmt.Errorf("load journal entry lines: %w", err)
	}
	if len(items) == 0 {
		return fmt.Errorf("journal entry has no lines to reverse")
	}

	doc := finance_accounts.SourceDoc{
		Type:        CreditNoteSourceDocType,
		SourceDocID: 0,
	}
	if err := tx.Create(&doc).Error; err != nil {
		return err
	}
	c.pendingReversalSourceDocID = doc.ID

	dt := c.Datetime
	if dt.IsZero() {
		dt = time.Now()
	}
	rev := finance_accounts.JournalEntry{
		Datetime:    dt,
		SourceDocID: doc.ID,
		JournalID:   orig.JournalID,
	}
	if err := tx.Create(&rev).Error; err != nil {
		return err
	}
	for _, it := range items {
		line := finance_accounts.JournalEntryItem{
			Datetime:       dt,
			AccountID:      it.AccountID,
			Amount:         negateDecimalSix(it.Amount),
			JournalEntryID: rev.ID,
		}
		if err := tx.Create(&line).Error; err != nil {
			return err
		}
	}
	c.ReversedJournalEntryID = rev.ID
	return nil
}

// AfterCreate updates the reversal [finance_accounts.SourceDoc] to point at this credit note’s primary key.
func (c *CreditNote) AfterCreate(tx *gorm.DB) error {
	if c.pendingReversalSourceDocID == 0 {
		return nil
	}
	return tx.Model(&finance_accounts.SourceDoc{}).
		Where("id = ?", c.pendingReversalSourceDocID).
		Update("source_doc_id", c.ID).Error
}

func negateDecimalSix(a fields.DecimalSix) fields.DecimalSix {
	if a.R == nil {
		return fields.DecimalSix{R: big.NewRat(0, 1)}.NormalizeDecimals()
	}
	return fields.DecimalSix{R: new(big.Rat).Neg(a.R)}.NormalizeDecimals()
}

func init() {
	lago.RegistryAdmin.Register("p_uniquity_finance_creditnotes.CreditNote", lago.AdminPanel[CreditNote]{
		SearchField: "Reason",
		ListFields: []string{
			"Datetime", "Reason", "JournalEntryID", "ReversedJournalEntryID", "UpdatedAt",
		},
	})
}
