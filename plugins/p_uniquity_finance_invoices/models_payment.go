package p_uniquity_finance_invoices

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	finance_accounts "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_accounts"
	finance_products "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_products"
	finance_taxes "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_taxes"
	"github.com/lariv-in/lago/fields"
	"gorm.io/gorm"
)

// PaymentSourceDocType is stored on [finance_accounts.SourceDoc.Type] for invoice receipts.
const PaymentSourceDocType = "p_uniquity_finance_invoices.Payment"

// Payment records settlement applied to a posted invoice AR balance (immutable after create).
// Amount is the settlement credited to AR; bank receipt is Amount minus payment-time withholding.
type Payment struct {
	gorm.Model

	PostedInvoiceID uint          `gorm:"column:posted_invoice_id;not null"`
	PostedInvoice   PostedInvoice `gorm:"foreignKey:PostedInvoiceID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Amount fields.DecimalSix `gorm:"column:amount;type:numeric(19,6);not null"`

	Taxes []finance_taxes.Tax `gorm:"many2many:payment_taxes"`

	AccountID uint                     `gorm:"column:account_id;not null"`
	Account   finance_accounts.Account `gorm:"foreignKey:AccountID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	Datetime time.Time `gorm:"column:datetime;not null"`

	JournalEntryID uint                          `gorm:"column:journal_entry_id;not null"`
	JournalEntry   finance_accounts.JournalEntry `gorm:"foreignKey:JournalEntryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	pendingSourceDocID uint `gorm:"-"`
	isFullPaymentHook  bool `gorm:"-"`
}

// postedInvoiceReceivableTotal returns the invoice amount due (line totals plus document-level taxes).
func postedInvoiceReceivableTotal(tx *gorm.DB, postedID uint) (fields.DecimalSix, error) {
	var posted PostedInvoice
	if err := tx.Preload("Taxes").Preload("Lines.Taxes").First(&posted, postedID).Error; err != nil {
		return fields.DecimalSix{}, err
	}
	totals, lineTaxIDs := accumulatePostedInvoiceLineTotals(posted.Lines)
	return invoiceReceivableGrandTotal(totals, posted.Taxes, lineTaxIDs), nil
}

// postedInvoiceOpenBalance returns invoice amount due minus payments already recorded.
func postedInvoiceOpenBalance(tx *gorm.DB, postedID uint) (fields.DecimalSix, error) {
	invTotal, err := postedInvoiceReceivableTotal(tx, postedID)
	if err != nil {
		return fields.DecimalSix{}, err
	}
	var posted PostedInvoice
	if err := tx.Select("id").First(&posted, postedID).Error; err != nil {
		return fields.DecimalSix{}, err
	}
	appliedSum, err := sumPostedInvoicePayments(tx, posted.ID)
	if err != nil {
		return fields.DecimalSix{}, err
	}
	open := decSub(invTotal, appliedSum).NormalizeDecimals()
	if open.R != nil && open.R.Sign() < 0 {
		return fields.DecimalSix{R: big.NewRat(0, 1)}.NormalizeDecimals(), nil
	}
	return open, nil
}

func sumPostedInvoicePayments(tx *gorm.DB, postedInvoiceID uint) (fields.DecimalSix, error) {
	type row struct {
		S fields.DecimalSix `gorm:"column:s"`
	}
	var out row
	err := tx.Raw(`
		SELECT COALESCE(SUM(amount), 0) AS s FROM payments
		WHERE posted_invoice_id = ? AND deleted_at IS NULL
	`, postedInvoiceID).Scan(&out).Error
	return out.S.NormalizeDecimals(), err
}

func decCmpPayment(a, b fields.DecimalSix) int {
	ar, br := big.NewRat(0, 1), big.NewRat(0, 1)
	if a.R != nil {
		ar.Set(a.R)
	}
	if b.R != nil {
		br.Set(b.R)
	}
	return ar.Cmp(br)
}

// BeforeCreate validates the payment and inserts the balancing journal entry:
// Dr bank (settlement − withholding), Dr withholding tax accounts, Cr AR (settlement).
func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p == nil {
		return fmt.Errorf("payment required")
	}
	if p.PostedInvoiceID == 0 {
		return fmt.Errorf("posted invoice is required")
	}
	if p.JournalEntryID != 0 {
		return fmt.Errorf("journal entry must not be set manually")
	}
	if p.Amount.R == nil || p.Amount.R.Sign() <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	paymentPrefs := LoadPaymentPreferences(tx)
	if err := ValidatePaymentPreferencesForCreate(tx, &paymentPrefs); err != nil {
		return err
	}
	if p.AccountID == 0 {
		p.AccountID = finance_products.OptionalUintValue(paymentPrefs.PaymentAccountID)
	}
	if err := finance_accounts.ValidateLeafAccountBalanceType(tx, p.AccountID, finance_accounts.BalanceTypeDebit, "payment account"); err != nil {
		return err
	}

	var posted PostedInvoice
	if err := tx.First(&posted, p.PostedInvoiceID).Error; err != nil {
		return fmt.Errorf("posted invoice: %w", err)
	}

	var cancelledCount int64
	if err := tx.Model(&CancelledInvoice{}).
		Where("posted_invoice_id = ? AND deleted_at IS NULL", posted.ID).
		Count(&cancelledCount).Error; err != nil {
		return err
	}
	if cancelledCount > 0 {
		return fmt.Errorf("cannot pay a cancelled invoice")
	}

	var paidCount int64
	if err := tx.Model(&PaidInvoice{}).
		Where("posted_invoice_id = ? AND deleted_at IS NULL", posted.ID).
		Count(&paidCount).Error; err != nil {
		return err
	}
	if paidCount > 0 {
		return fmt.Errorf("invoice is already fully paid")
	}

	invTotal, err := postedInvoiceReceivableTotal(tx, posted.ID)
	if err != nil {
		return err
	}

	appliedSum, err := sumPostedInvoicePayments(tx, posted.ID)
	if err != nil {
		return err
	}
	totalAfter := decSum(appliedSum, p.Amount).NormalizeDecimals()
	if decCmpPayment(totalAfter, invTotal) > 0 {
		return fmt.Errorf("payment exceeds open balance")
	}
	p.isFullPaymentHook = decCmpPayment(totalAfter, invTotal) == 0

	dt := p.Datetime
	if dt.IsZero() {
		dt = time.Now()
	}

	doc := finance_accounts.SourceDoc{Type: PaymentSourceDocType, SourceDocID: 0}
	if err := tx.Create(&doc).Error; err != nil {
		return err
	}
	p.pendingSourceDocID = doc.ID

	je := finance_accounts.JournalEntry{
		Datetime:    dt,
		SourceDocID: doc.ID,
		JournalID:   posted.JournalID,
	}
	if err := tx.Create(&je).Error; err != nil {
		return err
	}

	taxes := paymentTaxesFromContext(tx.Statement.Context)
	if err := validatePaymentTaxes(taxes); err != nil {
		return err
	}
	settlement := p.Amount.NormalizeDecimals()
	bankAmt := paymentBankAmount(settlement, taxes)
	if bankAmt.R != nil && bankAmt.R.Sign() < 0 {
		return fmt.Errorf("withholding exceeds settlement amount")
	}

	items := []finance_accounts.JournalEntryItem{
		{
			Datetime:       dt,
			AccountID:      p.AccountID,
			Amount:         bankAmt,
			JournalEntryID: je.ID,
		},
		{
			Datetime:       dt,
			AccountID:      posted.AccountReceivableID,
			Amount:         decNeg(settlement),
			JournalEntryID: je.ID,
		},
	}
	for _, tax := range taxesWithholding(taxes) {
		whAmt := taxAmountForTax(settlement, tax)
		if whAmt.R == nil || whAmt.R.Sign() == 0 {
			continue
		}
		acctID, err := withholdingTaxAccountID(tax)
		if err != nil {
			return err
		}
		items = append(items, finance_accounts.JournalEntryItem{
			Datetime:       dt,
			AccountID:      acctID,
			Amount:         whAmt,
			JournalEntryID: je.ID,
		})
	}
	if err := tx.Create(&items).Error; err != nil {
		return err
	}

	var loaded []finance_accounts.JournalEntryItem
	if err := tx.Where("journal_entry_id = ?", je.ID).Order("id ASC").Find(&loaded).Error; err != nil {
		return err
	}
	if ratSumBalance(loaded).Sign() != 0 {
		return fmt.Errorf("internal error: payment journal entry does not balance")
	}
	p.JournalEntryID = je.ID
	return nil
}

// AfterCreate links the source document and inserts PaidInvoice or PartiallyPaidInvoice.
func (p *Payment) AfterCreate(tx *gorm.DB) error {
	if p.pendingSourceDocID != 0 {
		if err := tx.Model(&finance_accounts.SourceDoc{}).
			Where("id = ?", p.pendingSourceDocID).
			Update("source_doc_id", p.ID).Error; err != nil {
			return err
		}
	}

	var priorID *uint
	var prev PartiallyPaidInvoice
	err := tx.Where("posted_invoice_id = ? AND deleted_at IS NULL", p.PostedInvoiceID).
		Order("id DESC").Take(&prev).Error
	if err == nil {
		pid := prev.ID
		priorID = &pid
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if p.isFullPaymentHook {
		row := PaidInvoice{
			PaymentID:                   p.ID,
			PostedInvoiceID:             p.PostedInvoiceID,
			PriorPartiallyPaidInvoiceID: priorID,
		}
		return tx.Create(&row).Error
	}
	row := PartiallyPaidInvoice{
		PaymentID:                   p.ID,
		PostedInvoiceID:             p.PostedInvoiceID,
		PriorPartiallyPaidInvoiceID: priorID,
	}
	return tx.Create(&row).Error
}
