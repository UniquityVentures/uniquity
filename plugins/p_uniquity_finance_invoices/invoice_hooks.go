package p_uniquity_finance_invoices

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/UniquityVentures/lamu/fields"
	finance_products "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_products"
	finance_taxes "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_taxes"
	"gorm.io/gorm"
)

func (d *DraftInvoice) BeforeSave(tx *gorm.DB) error {
	if d.Number != nil {
		t := strings.TrimSpace(*d.Number)
		if t == "" {
			d.Number = nil
		} else {
			d.Number = &t
		}
	}
	if d.PaymentTermID == 0 {
		return fmt.Errorf("payment term is required")
	}
	var pt PaymentTerm
	if err := tx.First(&pt, d.PaymentTermID).Error; err != nil {
		return fmt.Errorf("load payment term: %w", err)
	}
	d.PaymentTermType = pt.Type
	return nil
}

func loadHeaderTaxes(tx *gorm.DB, d *DraftInvoice) ([]finance_taxes.Tax, error) {
	var headerTaxes []finance_taxes.Tax
	if ctx := tx.Statement.Context; ctx != nil {
		if r, _ := ctx.Value("$request").(*http.Request); r != nil {
			_ = r.ParseForm()
			taxIDsStr := r.Form["Taxes"]
			if len(taxIDsStr) > 0 {
				var taxIDs []uint
				for _, s := range taxIDsStr {
					if id, err := strconv.ParseUint(s, 10, 64); err == nil {
						taxIDs = append(taxIDs, uint(id))
					}
				}
				if len(taxIDs) > 0 {
					if err := tx.Where("id IN ?", taxIDs).Find(&headerTaxes).Error; err != nil {
						return nil, err
					}
				}
			}
		}
	}
	if len(headerTaxes) == 0 && d.ID > 0 {
		if err := tx.Model(d).Association("Taxes").Find(&headerTaxes); err != nil {
			return nil, err
		}
	}
	return headerTaxes, nil
}

func (d *DraftInvoice) AfterCreate(tx *gorm.DB) error {
	if len(d.PendingLines) == 0 {
		return nil
	}
	headerTaxes, err := loadHeaderTaxes(tx, d)
	if err != nil {
		return err
	}
	for _, row := range d.PendingLines {
		if row.ProductID == 0 {
			continue
		}
		ln, err := buildDraftLineFromPending(tx, d.ID, row, headerTaxes)
		if err != nil {
			return err
		}
		if err := tx.Create(ln.line).Error; err != nil {
			return err
		}
		if len(ln.taxesToAssociate) > 0 {
			if err := tx.Model(ln.line).Association("Taxes").Append(ln.taxesToAssociate); err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *DraftInvoice) AfterUpdate(tx *gorm.DB) error {
	if len(d.PendingLines) == 0 {
		return nil
	}
	if err := tx.Where("draft_invoice_id = ?", d.ID).Delete(&DraftInvoiceLine{}).Error; err != nil {
		return err
	}
	headerTaxes, err := loadHeaderTaxes(tx, d)
	if err != nil {
		return err
	}
	for _, row := range d.PendingLines {
		if row.ProductID == 0 {
			continue
		}
		ln, err := buildDraftLineFromPending(tx, d.ID, row, headerTaxes)
		if err != nil {
			return err
		}
		if err := tx.Create(ln.line).Error; err != nil {
			return err
		}
		if len(ln.taxesToAssociate) > 0 {
			if err := tx.Model(ln.line).Association("Taxes").Append(ln.taxesToAssociate); err != nil {
				return err
			}
		}
	}
	return nil
}

type draftLineBuild struct {
	line             *DraftInvoiceLine
	taxesToAssociate []finance_taxes.Tax
}

func buildDraftLineFromPending(tx *gorm.DB, draftID uint, row DraftLinePending, headerTaxes []finance_taxes.Tax) (*draftLineBuild, error) {
	var qty fields.DecimalSix
	if err := qty.UnmarshalText([]byte(strings.TrimSpace(row.Quantity))); err != nil {
		return nil, fmt.Errorf("quantity: %w", err)
	}
	qty = qty.NormalizeDecimals()
	if qty.R == nil || qty.R.Sign() <= 0 {
		return nil, fmt.Errorf("quantity must be positive")
	}
	var prod finance_products.Product
	if err := tx.Preload("Taxes", nil).First(&prod, row.ProductID).Error; err != nil {
		return nil, fmt.Errorf("load product %d: %w", row.ProductID, err)
	}
	var rate fields.DecimalSix
	rateStr := strings.TrimSpace(row.Rate)
	if rateStr == "" {
		rate = prod.SalesPrice.NormalizeDecimals()
	} else {
		if err := rate.UnmarshalText([]byte(rateStr)); err != nil {
			return nil, fmt.Errorf("rate: %w", err)
		}
		rate = rate.NormalizeDecimals()
		if rate.R == nil || rate.R.Sign() < 0 {
			return nil, fmt.Errorf("rate must be non-negative")
		}
	}
	var merged []finance_taxes.Tax
	if row.TaxIDs != nil {
		if len(row.TaxIDs) > 0 {
			if err := tx.Where("id IN ?", row.TaxIDs).Find(&merged).Error; err != nil {
				return nil, fmt.Errorf("load line taxes: %w", err)
			}
			seen := map[uint]struct{}{}
			for _, id := range row.TaxIDs {
				seen[id] = struct{}{}
			}
			if len(merged) != len(seen) {
				return nil, fmt.Errorf("one or more line tax ids are invalid")
			}
		}
	} else {
		merged = mergeTaxesUnique(append([]finance_taxes.Tax{}, headerTaxes...), prod.Taxes)
	}
	line := &DraftInvoiceLine{
		DraftInvoiceID: draftID,
		ProductID:      row.ProductID,
		Rate:           rate,
		Quantity:       qty,
	}
	return &draftLineBuild{line: line, taxesToAssociate: merged}, nil
}

func mergeTaxesUnique(a, b []finance_taxes.Tax) []finance_taxes.Tax {
	seen := map[uint]struct{}{}
	var out []finance_taxes.Tax
	for _, t := range append(a, b...) {
		if t.ID == 0 {
			continue
		}
		if _, ok := seen[t.ID]; ok {
			continue
		}
		seen[t.ID] = struct{}{}
		out = append(out, t)
	}
	return out
}

func (d *DraftInvoice) BeforeUpdate(tx *gorm.DB) error {
	return errIfDraftSealed(tx, d.ID)
}

func (d *DraftInvoice) BeforeDelete(tx *gorm.DB) error {
	return errIfDraftSealed(tx, d.ID)
}

func errIfDraftSealed(tx *gorm.DB, draftID uint) error {
	if draftID == 0 {
		return nil
	}
	var n int64
	if err := tx.Model(&PostedInvoice{}).Where("draft_invoice_id = ? AND deleted_at IS NULL", draftID).Count(&n).Error; err != nil {
		return err
	}
	if n > 0 {
		return errors.New("draft invoice is posted and cannot be changed")
	}
	return nil
}

func (p *PostedInvoice) BeforeUpdate(tx *gorm.DB) error {
	return errIfPostedSealed(tx, p.ID)
}

func (p *PostedInvoice) BeforeDelete(tx *gorm.DB) error {
	return errIfPostedSealed(tx, p.ID)
}

func errIfPostedSealed(tx *gorm.DB, postedID uint) error {
	if postedID == 0 {
		return nil
	}
	var n int64
	if err := tx.Model(&CancelledInvoice{}).Where("posted_invoice_id = ? AND deleted_at IS NULL", postedID).Count(&n).Error; err != nil {
		return err
	}
	if n > 0 {
		return errors.New("posted invoice is cancelled and cannot be changed")
	}
	return nil
}

func (d *DraftInvoiceLine) BeforeCreate(_ *gorm.DB) error {
	d.Rate = d.Rate.NormalizeDecimals()
	d.Quantity = d.Quantity.NormalizeDecimals()
	return nil
}

func (d *DraftInvoiceLine) BeforeUpdate(_ *gorm.DB) error {
	d.Rate = d.Rate.NormalizeDecimals()
	d.Quantity = d.Quantity.NormalizeDecimals()
	return nil
}
