package p_uniquity_finance_invoices

import (
	"fmt"
	"math/big"
	"time"

	finance_accounts "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_accounts"
	finance_creditnotes "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_creditnotes"
	finance_products "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_products"
	finance_taxes "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_taxes"
	"github.com/lariv-in/lago/fields"
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

func decSub(a, b fields.DecimalSix) fields.DecimalSix {
	return decSum(a, decNeg(b))
}

func decAbs(a fields.DecimalSix) fields.DecimalSix {
	if a.R == nil || a.R.Sign() >= 0 {
		return a.NormalizeDecimals()
	}
	return decNeg(a)
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
//   - Sales per line: Cr revenue (−base); levied taxes Cr tax payable (−); withholding Dr [Tax.AccountID] (+).
//   - AR is the net of base + levied − withholding (line and document-level).
//   - Cost per line: Dr COGS (+cost), Cr inventory (−); input tax uses levied percentages only on cost base.
//   - Document-level header taxes not already on lines follow the same levied vs withholding rules.
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
	allTaxes := append(collectTaxesFromLines(lines), full.Taxes...)
	if err := validateWithholdingTaxAccounts(allTaxes); err != nil {
		return nil, err
	}
	productPrefs := finance_products.LoadProductPreferences(tx)
	if finance_products.OptionalUintValue(productPrefs.InventoryAccountID) == 0 || finance_products.OptionalUintValue(productPrefs.CostOfSalesAcctID) == 0 {
		return nil, fmt.Errorf("product preferences must have inventory and cost-of-sales accounts for posting")
	}
	invoicePrefs := LoadInvoicePreferences(tx)
	if err := ValidateInvoicePreferencesForPosting(tx, &invoicePrefs); err != nil {
		return nil, err
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
		JournalID:   finance_products.OptionalUintValue(invoicePrefs.JournalID),
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
	revItemSpecIndex := make([]int, len(lines))

	for i, line := range lines {
		lineBase := decMul(line.Rate, line.Quantity)
		leviedTax := taxAmountOnBase(lineBase, sumTaxPercents(taxesLevied(line.Taxes)))
		revItemSpecIndex[i] = len(specs)
		specs = append(specs, itemSpec{accountID: finance_products.OptionalUintValue(invoicePrefs.AccountRevenueID), amount: decNeg(lineBase)})
		if leviedTax.R != nil && leviedTax.R.Sign() != 0 {
			specs = append(specs, itemSpec{accountID: finance_products.OptionalUintValue(invoicePrefs.AccountTaxPayableID), amount: decNeg(leviedTax)})
		}
		for _, tax := range taxesWithholding(line.Taxes) {
			whAmt := taxAmountForTax(lineBase, tax)
			if whAmt.R == nil || whAmt.R.Sign() == 0 {
				continue
			}
			acctID, err := withholdingTaxAccountID(tax)
			if err != nil {
				return nil, err
			}
			specs = append(specs, itemSpec{accountID: acctID, amount: whAmt})
		}
	}

	for _, line := range lines {
		p := line.Product
		qty := line.Quantity
		costBase := decMul(p.BaseCost, qty)
		specs = append(
			specs,
			itemSpec{accountID: finance_products.OptionalUintValue(productPrefs.CostOfSalesAcctID), amount: costBase},
			itemSpec{accountID: finance_products.OptionalUintValue(productPrefs.InventoryAccountID), amount: decNeg(costBase)},
		)
	}

	var lineTotals invoiceLinesTotals
	lineTaxIDs := map[uint]struct{}{}
	for _, line := range lines {
		u, lev, wh, _ := invoiceLineAmountBreakdown(line.Quantity, line.Rate, line.Taxes)
		lineTotals.UntaxedSubtotal = decSum(lineTotals.UntaxedSubtotal, u)
		lineTotals.LinesLevied = decSum(lineTotals.LinesLevied, lev)
		lineTotals.LinesWithholding = decSum(lineTotals.LinesWithholding, wh)
		mergeInvoiceLineTaxIDs(lineTaxIDs, line.Taxes)
	}
	for _, tax := range documentLevelHeaderTaxes(full.Taxes, lineTaxIDs) {
		amt := taxAmountForTax(lineTotals.UntaxedSubtotal, tax)
		if amt.R == nil || amt.R.Sign() == 0 {
			continue
		}
		if effectiveTaxKind(tax) == finance_taxes.TaxKindWithholding {
			acctID, err := withholdingTaxAccountID(tax)
			if err != nil {
				return nil, err
			}
			specs = append(specs, itemSpec{accountID: acctID, amount: amt})
			continue
		}
		specs = append(specs, itemSpec{accountID: finance_products.OptionalUintValue(invoicePrefs.AccountTaxPayableID), amount: decNeg(amt)})
	}
	totalAR := invoiceReceivableGrandTotal(lineTotals, full.Taxes, lineTaxIDs)
	specs = append(specs, itemSpec{accountID: finance_products.OptionalUintValue(invoicePrefs.AccountReceivableID), amount: totalAR})

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
		AccountReceivableID: finance_products.OptionalUintValue(invoicePrefs.AccountReceivableID),
		AccountRevenueID:    finance_products.OptionalUintValue(invoicePrefs.AccountRevenueID),
		AccountTaxPayableID: finance_products.OptionalUintValue(invoicePrefs.AccountTaxPayableID),
		JournalID:           finance_products.OptionalUintValue(invoicePrefs.JournalID),
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
		revIdx := revItemSpecIndex[i]
		if revIdx >= len(journalItems) {
			return nil, fmt.Errorf("internal error: revenue item index")
		}
		pLine := PostedInvoiceLine{
			PostedInvoiceID:    posted.ID,
			ProductID:          line.ProductID,
			Rate:               line.Rate,
			Quantity:           line.Quantity,
			JournalEntryItemID: journalItems[revIdx].ID,
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
		Number:          nil,
		Datetime:        full.Datetime,
		CustomerID:      full.CustomerID,
		PaymentTermType: full.PaymentTermType,
		PaymentTermID:   full.PaymentTermID,
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
