package p_uniquity_finance_invoices

import (
	"fmt"
	"math/big"
	"time"

	"github.com/UniquityVentures/lamu/fields"
	finance_accounts "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_accounts"
	finance_creditnotes "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_creditnotes"
	finance_taxes "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_taxes"
	"gorm.io/gorm"
)

func decMul(a, b fields.DecimalSix) fields.DecimalSix {
	if a.R == nil || b.R == nil {
		return fields.DecimalSix{R: big.NewRat(0, 1)}.NormalizeDecimals()
	}
	return fields.DecimalSix{R: new(big.Rat).Mul(a.R, b.R)}.NormalizeDecimals()
}

func decSum(a, b fields.DecimalSix) fields.DecimalSix {
	ar, br := big.NewRat(0, 1), big.NewRat(0, 1)
	if a.R != nil {
		ar.Set(a.R)
	}
	if b.R != nil {
		br.Set(b.R)
	}
	return fields.DecimalSix{R: new(big.Rat).Add(ar, br)}.NormalizeDecimals()
}

func decNeg(a fields.DecimalSix) fields.DecimalSix {
	if a.R == nil {
		return fields.DecimalSix{R: big.NewRat(0, 1)}.NormalizeDecimals()
	}
	return fields.DecimalSix{R: new(big.Rat).Neg(a.R)}.NormalizeDecimals()
}

func sumTaxPercents(taxes []finance_taxes.Tax) fields.DecimalSix {
	acc := big.NewRat(0, 1)
	for _, t := range taxes {
		if t.Percentage.R != nil {
			acc.Add(acc, t.Percentage.R)
		}
	}
	return fields.DecimalSix{R: acc}.NormalizeDecimals()
}

func taxAmountOnBase(base, pctSum fields.DecimalSix) fields.DecimalSix {
	if base.R == nil || base.R.Sign() == 0 {
		return fields.DecimalSix{R: big.NewRat(0, 1)}.NormalizeDecimals()
	}
	if pctSum.R == nil || pctSum.R.Sign() == 0 {
		return fields.DecimalSix{R: big.NewRat(0, 1)}.NormalizeDecimals()
	}
	hundred := big.NewRat(100, 1)
	r := new(big.Rat).Quo(pctSum.R, hundred)
	return decMul(base, fields.DecimalSix{R: r})
}

// invoiceLineAmountBreakdown returns untaxed (qty×rate), tax on that base, and line total.
// It matches the draft line editor and [taxAmountOnBase]/[sumTaxPercents] used when posting.
func invoiceLineAmountBreakdown(qty, rate fields.DecimalSix, taxes []finance_taxes.Tax) (untaxed, taxAmt, lineTotal fields.DecimalSix) {
	untaxed = decMul(qty, rate)
	taxAmt = taxAmountOnBase(untaxed, sumTaxPercents(taxes))
	lineTotal = decSum(untaxed, taxAmt)
	return
}

func ratSumBalance(items []finance_accounts.JournalEntryItem) *big.Rat {
	s := big.NewRat(0, 1)
	for _, it := range items {
		if it.Amount.R != nil {
			s.Add(s, it.Amount.R)
		}
	}
	return s
}

// NewPosted creates a posted invoice, journal entry, and lines from a draft inside tx.
//
// Posting (single balanced journal entry; signed amounts must sum to zero):
//   - Sales per line: line base = rate × quantity; sales tax = base × (sum of line tax percentages / 100),
//     union of header and product taxes on the line via [mergeTaxesUnique] in draft line creation.
//     For each line: Cr revenue (−base), Cr tax payable (−sales tax).
//   - Cost per line: cost base = product.BaseCost × quantity; Dr COGS (+cost base), Cr inventory (−cost base).
//     Requires product inventory and COGS accounts.
//   - Cost tax per line: if product has input-tax account and line has taxes, tax on cost base with same
//     percentages; Dr input tax (+), Cr inventory (−).
//   - Closing: Dr accounts receivable (+) for total of all line bases and sales taxes.
//   - [PostedInvoiceLine.JournalEntryItemID] points at the revenue credit line for that invoice line (first
//     of the two sales items per line in entry build order).
func (d *DraftInvoice) NewPosted(tx *gorm.DB, postedAt time.Time) (*PostedInvoice, error) {
	if d == nil || d.ID == 0 {
		return nil, fmt.Errorf("draft invoice required")
	}
	var n int64
	if err := tx.Model(&PostedInvoice{}).Where("draft_invoice_id = ? AND deleted_at IS NULL", d.ID).Count(&n).Error; err != nil {
		return nil, err
	}
	if n > 0 {
		return nil, fmt.Errorf("draft already posted")
	}
	var full DraftInvoice
	if err := tx.Preload("Taxes", nil).Preload("Customer", nil).Preload("PaymentTerm", nil).
		First(&full, d.ID).Error; err != nil {
		return nil, err
	}
	var lines []DraftInvoiceLine
	if err := tx.Where("draft_invoice_id = ? AND deleted_at IS NULL", full.ID).Order("id ASC").
		Preload("Taxes", nil).Preload("Product.Taxes", nil).Preload("Product", nil).Find(&lines).Error; err != nil {
		return nil, err
	}
	if len(lines) == 0 {
		return nil, fmt.Errorf("draft has no lines")
	}
	for i := range lines {
		p := lines[i].Product
		if p.InventoryAccountID == nil || *p.InventoryAccountID == 0 || p.CostOfSalesAcctID == nil || *p.CostOfSalesAcctID == 0 {
			return nil, fmt.Errorf("product %q must have inventory and cost-of-sales accounts for posting", p.Name)
		}
	}
	prefs := finance_accounts.LoadAccountingPreferences(tx)
	number, err := PostedInvoiceNumber(tx, &full, prefs)
	if err != nil {
		return nil, err
	}
	var dupPosted int64
	if err := tx.Model(&PostedInvoice{}).Where("number = ? AND deleted_at IS NULL", number).Count(&dupPosted).Error; err != nil {
		return nil, err
	}
	if dupPosted > 0 {
		return nil, fmt.Errorf("invoice number %q is already used by another posted invoice", number)
	}
	var dupCancelled int64
	if err := tx.Model(&CancelledInvoice{}).Where("number = ? AND deleted_at IS NULL", number).Count(&dupCancelled).Error; err != nil {
		return nil, err
	}
	if dupCancelled > 0 {
		return nil, fmt.Errorf("invoice number %q is already used by a cancelled invoice", number)
	}
	if postedAt.IsZero() {
		postedAt = time.Now()
	}
	doc := finance_accounts.SourceDoc{Type: PostedInvoiceSourceDocType, SourceDocID: 0}
	if err := tx.Create(&doc).Error; err != nil {
		return nil, err
	}
	je := finance_accounts.JournalEntry{
		Datetime:    full.Datetime,
		SourceDocID: doc.ID,
		JournalID:   full.JournalID,
	}
	if err := tx.Create(&je).Error; err != nil {
		return nil, err
	}

	var journalItems []finance_accounts.JournalEntryItem
	type itemSpec struct {
		accountID uint
		amount    fields.DecimalSix
	}
	var specs []itemSpec

	for _, line := range lines {
		lineBase := decMul(line.Rate, line.Quantity)
		sp := sumTaxPercents(line.Taxes)
		salesTax := taxAmountOnBase(lineBase, sp)
		specs = append(specs,
			itemSpec{accountID: full.AccountRevenueID, amount: decNeg(lineBase)},
			itemSpec{accountID: full.AccountTaxPayableID, amount: decNeg(salesTax)},
		)
	}

	for _, line := range lines {
		p := line.Product
		qty := line.Quantity
		costBase := decMul(p.BaseCost, qty)
		specs = append(specs,
			itemSpec{accountID: *p.CostOfSalesAcctID, amount: costBase},
			itemSpec{accountID: *p.InventoryAccountID, amount: decNeg(costBase)},
		)
		if p.InputTaxAccountID != nil && *p.InputTaxAccountID != 0 && len(line.Taxes) > 0 {
			ct := taxAmountOnBase(costBase, sumTaxPercents(line.Taxes))
			if ct.R != nil && ct.R.Sign() != 0 {
				specs = append(specs,
					itemSpec{accountID: *p.InputTaxAccountID, amount: ct},
					itemSpec{accountID: *p.InventoryAccountID, amount: decNeg(ct)},
				)
			}
		}
	}

	var totalAR fields.DecimalSix
	for _, line := range lines {
		lineBase := decMul(line.Rate, line.Quantity)
		st := taxAmountOnBase(lineBase, sumTaxPercents(line.Taxes))
		totalAR = decSum(totalAR, decSum(lineBase, st))
	}
	specs = append(specs, itemSpec{accountID: full.AccountReceivableID, amount: totalAR})

	for _, sp := range specs {
		it := finance_accounts.JournalEntryItem{
			Datetime:       full.Datetime,
			AccountID:      sp.accountID,
			Amount:         sp.amount.NormalizeDecimals(),
			JournalEntryID: je.ID,
		}
		if err := tx.Create(&it).Error; err != nil {
			return nil, err
		}
		journalItems = append(journalItems, it)
	}
	if ratSumBalance(journalItems).Sign() != 0 {
		return nil, fmt.Errorf("internal error: journal entry does not balance")
	}

	posted := PostedInvoice{
		DraftInvoiceID:      full.ID,
		PostedAt:            &postedAt,
		Number:              number,
		AccountReceivableID: full.AccountReceivableID,
		AccountRevenueID:    full.AccountRevenueID,
		AccountTaxPayableID: full.AccountTaxPayableID,
		JournalID:           full.JournalID,
		Datetime:            full.Datetime,
		CustomerID:          full.CustomerID,
		PaymentTermType:     full.PaymentTermType,
		PaymentTermID:       full.PaymentTermID,
		JournalEntryID:      je.ID,
	}
	if err := tx.Create(&posted).Error; err != nil {
		return nil, err
	}
	if err := tx.Model(&finance_accounts.SourceDoc{}).Where("id = ?", doc.ID).Update("source_doc_id", posted.ID).Error; err != nil {
		return nil, err
	}

	for i, line := range lines {
		revItemIdx := 2 * i
		if revItemIdx >= len(journalItems) {
			return nil, fmt.Errorf("internal error: revenue item index")
		}
		pLine := PostedInvoiceLine{
			PostedInvoiceID:    posted.ID,
			ProductID:          line.ProductID,
			Rate:               line.Rate,
			Quantity:           line.Quantity,
			JournalEntryItemID: journalItems[revItemIdx].ID,
		}
		if err := tx.Create(&pLine).Error; err != nil {
			return nil, err
		}
		if len(line.Taxes) > 0 {
			if err := tx.Model(&pLine).Association("Taxes").Append(line.Taxes); err != nil {
				return nil, err
			}
		}
	}
	if len(full.Taxes) > 0 {
		if err := tx.Model(&posted).Association("Taxes").Append(full.Taxes); err != nil {
			return nil, err
		}
	}
	if err := tx.First(&posted, posted.ID).Error; err != nil {
		return nil, err
	}
	return &posted, nil
}

// NewCancelled creates a credit note reversal and cancelled invoice snapshot.
func (p *PostedInvoice) NewCancelled(tx *gorm.DB, reason string, at time.Time) (*CancelledInvoice, error) {
	if p == nil || p.ID == 0 {
		return nil, fmt.Errorf("posted invoice required")
	}
	var n int64
	if err := tx.Model(&CancelledInvoice{}).Where("posted_invoice_id = ? AND deleted_at IS NULL", p.ID).Count(&n).Error; err != nil {
		return nil, err
	}
	if n > 0 {
		return nil, fmt.Errorf("already cancelled")
	}
	var full PostedInvoice
	if err := tx.Preload("Taxes", nil).Preload("Lines", func(db *gorm.DB) *gorm.DB {
		return db.Order("posted_invoice_lines.id ASC")
	}).Preload("Lines.Taxes", nil).
		First(&full, p.ID).Error; err != nil {
		return nil, err
	}
	if at.IsZero() {
		at = time.Now()
	}
	cn := finance_creditnotes.CreditNote{
		Datetime:       at,
		Reason:         reason,
		JournalEntryID: full.JournalEntryID,
	}
	if err := tx.Create(&cn).Error; err != nil {
		return nil, err
	}

	var origItems []finance_accounts.JournalEntryItem
	if err := tx.Where("journal_entry_id = ?", full.JournalEntryID).Order("id ASC").Find(&origItems).Error; err != nil {
		return nil, err
	}
	var revItems []finance_accounts.JournalEntryItem
	if err := tx.Where("journal_entry_id = ?", cn.ReversedJournalEntryID).Order("id ASC").Find(&revItems).Error; err != nil {
		return nil, err
	}
	if len(origItems) != len(revItems) {
		return nil, fmt.Errorf("reversal line count mismatch")
	}
	origToRev := make(map[uint]uint, len(origItems))
	for i := range origItems {
		origToRev[origItems[i].ID] = revItems[i].ID
	}

	cancelled := CancelledInvoice{
		PostedInvoiceID:     full.ID,
		PostedAt:            full.PostedAt,
		CancelledAt:         &at,
		Number:              full.Number,
		AccountReceivableID: full.AccountReceivableID,
		AccountRevenueID:    full.AccountRevenueID,
		AccountTaxPayableID: full.AccountTaxPayableID,
		JournalID:           full.JournalID,
		Datetime:            full.Datetime,
		CustomerID:          full.CustomerID,
		PaymentTermType:     full.PaymentTermType,
		PaymentTermID:       full.PaymentTermID,
		CreditNoteID:        cn.ID,
	}
	if err := tx.Create(&cancelled).Error; err != nil {
		return nil, err
	}

	for _, pl := range full.Lines {
		revID, ok := origToRev[pl.JournalEntryItemID]
		if !ok {
			return nil, fmt.Errorf("could not map journal line for posted invoice line %d", pl.ID)
		}
		cl := CancelledInvoiceLine{
			CancelledInvoiceID: cancelled.ID,
			ProductID:          pl.ProductID,
			Rate:               pl.Rate,
			Quantity:           pl.Quantity,
			JournalEntryItemID: revID,
		}
		if err := tx.Create(&cl).Error; err != nil {
			return nil, err
		}
		if len(pl.Taxes) > 0 {
			if err := tx.Model(&cl).Association("Taxes").Append(pl.Taxes); err != nil {
				return nil, err
			}
		}
	}
	if len(full.Taxes) > 0 {
		if err := tx.Model(&cancelled).Association("Taxes").Append(full.Taxes); err != nil {
			return nil, err
		}
	}
	if err := tx.First(&cancelled, cancelled.ID).Error; err != nil {
		return nil, err
	}
	return &cancelled, nil
}

// NewDraft creates a new draft from a cancelled invoice snapshot.
func (c *CancelledInvoice) NewDraft(tx *gorm.DB) (*DraftInvoice, error) {
	if c == nil || c.ID == 0 {
		return nil, fmt.Errorf("cancelled invoice required")
	}
	var full CancelledInvoice
	if err := tx.Preload("Taxes", nil).Preload("Lines", func(db *gorm.DB) *gorm.DB {
		return db.Order("cancelled_invoice_lines.id ASC")
	}).Preload("Lines.Taxes", nil).
		First(&full, c.ID).Error; err != nil {
		return nil, err
	}
	draft := DraftInvoice{
		Number:              nil,
		AccountReceivableID: full.AccountReceivableID,
		AccountRevenueID:    full.AccountRevenueID,
		AccountTaxPayableID: full.AccountTaxPayableID,
		JournalID:           full.JournalID,
		Datetime:            full.Datetime,
		CustomerID:          full.CustomerID,
		PaymentTermType:     full.PaymentTermType,
		PaymentTermID:       full.PaymentTermID,
	}
	if err := tx.Create(&draft).Error; err != nil {
		return nil, err
	}
	if err := tx.Model(&draft).Association("Taxes").Append(full.Taxes); err != nil {
		return nil, err
	}
	for _, cl := range full.Lines {
		line := DraftInvoiceLine{
			DraftInvoiceID: draft.ID,
			ProductID:      cl.ProductID,
			Rate:           cl.Rate,
			Quantity:       cl.Quantity,
		}
		if err := tx.Create(&line).Error; err != nil {
			return nil, err
		}
		if len(cl.Taxes) > 0 {
			if err := tx.Model(&line).Association("Taxes").Append(cl.Taxes); err != nil {
				return nil, err
			}
		}
	}
	if err := tx.First(&draft, draft.ID).Error; err != nil {
		return nil, err
	}
	return &draft, nil
}
