package p_uniquity_finance_invoices

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	finance_products "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_products"
	finance_taxes "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_taxes"
	"github.com/lariv-in/lago/fields"
	"github.com/lariv-in/lago/getters"
	"github.com/lariv-in/lago/views"
)

// invoiceUpdateLinesPatcher validates [InvoiceLinesJSON] and sets [DraftInvoice.PendingLines] for [DraftInvoice.AfterUpdate].
type invoiceUpdateLinesPatcher struct{}

func (invoiceUpdateLinesPatcher) Patch(_ views.View, r *http.Request, formData map[string]any, formErrors map[string]error) (map[string]any, map[string]error) {
	if len(formErrors) > 0 {
		return formData, formErrors
	}
	raw, ok := formData["InvoiceLinesJSON"].(string)
	if !ok || strings.TrimSpace(raw) == "" {
		formErrors["InvoiceLinesJSON"] = fmt.Errorf("add at least one invoice line")
		return formData, formErrors
	}
	var rows []DraftLinePending
	if err := json.Unmarshal([]byte(raw), &rows); err != nil {
		formErrors["InvoiceLinesJSON"] = fmt.Errorf("invalid lines data: %w", err)
		return formData, formErrors
	}
	if len(rows) == 0 {
		formErrors["InvoiceLinesJSON"] = fmt.Errorf("add at least one invoice line")
		return formData, formErrors
	}
	db, err := getters.DBFromContext(r.Context())
	if err != nil {
		formErrors["_form"] = err
		return formData, formErrors
	}
	for idx := range rows {
		row := &rows[idx]
		if row.ProductID == 0 {
			formErrors["InvoiceLinesJSON"] = fmt.Errorf("line %d: choose a product", idx+1)
			return formData, formErrors
		}
		var cnt int64
		if err := db.WithContext(r.Context()).Model(&finance_products.Product{}).Where("id = ?", row.ProductID).Count(&cnt).Error; err != nil {
			formErrors["_form"] = err
			return formData, formErrors
		}
		if cnt != 1 {
			formErrors["InvoiceLinesJSON"] = fmt.Errorf("line %d: unknown product #%d", idx+1, row.ProductID)
			return formData, formErrors
		}
		var qty fields.DecimalSix
		if err := qty.UnmarshalText([]byte(strings.TrimSpace(row.Quantity))); err != nil {
			formErrors["InvoiceLinesJSON"] = fmt.Errorf("line %d: invalid quantity: %w", idx+1, err)
			return formData, formErrors
		}
		qty = qty.NormalizeDecimals()
		if qty.R == nil || qty.R.Sign() <= 0 {
			formErrors["InvoiceLinesJSON"] = fmt.Errorf("line %d: quantity must be positive", idx+1)
			return formData, formErrors
		}
		if strings.TrimSpace(row.Rate) != "" {
			var rate fields.DecimalSix
			if err := rate.UnmarshalText([]byte(strings.TrimSpace(row.Rate))); err != nil {
				formErrors["InvoiceLinesJSON"] = fmt.Errorf("line %d: invalid rate: %w", idx+1, err)
				return formData, formErrors
			}
			rate = rate.NormalizeDecimals()
			if rate.R != nil && rate.R.Sign() < 0 {
				formErrors["InvoiceLinesJSON"] = fmt.Errorf("line %d: rate must be non-negative", idx+1)
				return formData, formErrors
			}
		}
		if len(row.TaxIDs) > 0 {
			uniq := make([]uint, 0, len(row.TaxIDs))
			seen := map[uint]struct{}{}
			for _, id := range row.TaxIDs {
				if id == 0 {
					continue
				}
				if _, ok := seen[id]; ok {
					continue
				}
				seen[id] = struct{}{}
				uniq = append(uniq, id)
			}
			row.TaxIDs = uniq
			if len(row.TaxIDs) > 0 {
				var taxCnt int64
				if err := db.WithContext(r.Context()).Model(&finance_taxes.Tax{}).Where("id IN ?", row.TaxIDs).Count(&taxCnt).Error; err != nil {
					formErrors["_form"] = err
					return formData, formErrors
				}
				if taxCnt != int64(len(row.TaxIDs)) {
					formErrors["InvoiceLinesJSON"] = fmt.Errorf("line %d: one or more tax selections are invalid", idx+1)
					return formData, formErrors
				}
			}
		}
	}
	formData["PendingLines"] = rows
	delete(formData, "InvoiceLinesJSON")
	return formData, formErrors
}
