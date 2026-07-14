package p_uniquity_finance_invoices

import (
	"context"
	"strconv"
	"strings"
	"time"

	finance_fiscal_year "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_fiscal_year"
	"github.com/lariv-in/lago/getters"
	"github.com/lariv-in/lago/registry"
	"gorm.io/gorm"
)

// FinanceInvoicesEnvironmentFiscalYearKey is the environment cookie key for the invoice list fiscal year scope.
const FinanceInvoicesEnvironmentFiscalYearKey = "finance_invoices_fiscal_year"

func fiscalYearsEnvironmentOptionsGetter(ctx context.Context) ([]registry.Pair[uint, string], error) {
	db, err := getters.DBFromContext(ctx)
	if err != nil {
		return nil, err
	}
	var rows []finance_fiscal_year.FiscalYear
	if err := db.WithContext(ctx).Order("starts_at DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]registry.Pair[uint, string], 0, len(rows))
	for _, fy := range rows {
		out = append(out, registry.Pair[uint, string]{Key: fy.ID, Value: fy.Code + " — " + fy.Name})
	}
	return out, nil
}

func invoiceFiscalYearEnvironmentDefault(ctx context.Context) (uint, error) {
	db, err := getters.DBFromContext(ctx)
	if err != nil {
		return 0, nil
	}
	fy, err := invoiceDefaultFiscalYearWhenUnset(db, ctx)
	if err != nil {
		return 0, nil
	}
	return fy.ID, nil
}

// invoiceDefaultFiscalYearWhenUnset picks the fiscal year that contains "now" (starts_at <= now <= ends_at).
// If none overlap the current instant, falls back to the newest active fiscal year by starts_at.
func invoiceDefaultFiscalYearWhenUnset(db *gorm.DB, ctx context.Context) (finance_fiscal_year.FiscalYear, error) {
	var fy finance_fiscal_year.FiscalYear
	now := time.Now()
	err := db.WithContext(ctx).
		Where("starts_at <= ? AND ends_at >= ?", now, now).
		Order("starts_at DESC").
		First(&fy).Error
	if err == nil {
		return fy, nil
	}
	err = db.WithContext(ctx).Where("is_active = ?", true).Order("starts_at DESC").First(&fy).Error
	if err != nil {
		return finance_fiscal_year.FiscalYear{}, err
	}
	return fy, nil
}

// selectedInvoiceListFiscalYear resolves the fiscal year scope for the invoice list from the environment cookie.
// If the user explicitly clears the selector ("—"), restrict is false (all invoices).
// If the cookie has no selection yet, the current fiscal year (today within Start–End) is used when one exists.
func selectedInvoiceListFiscalYear(db *gorm.DB, ctx context.Context) (fy finance_fiscal_year.FiscalYear, restrict bool) {
	envMap, ok := ctx.Value("$environment").(map[string]string)
	var raw string
	var hasKey bool
	if ok {
		raw, hasKey = envMap[FinanceInvoicesEnvironmentFiscalYearKey]
	}
	if hasKey && strings.TrimSpace(raw) == "" {
		return finance_fiscal_year.FiscalYear{}, false
	}
	var id uint
	if s := strings.TrimSpace(raw); s != "" {
		parsed, err := strconv.ParseUint(s, 10, 64)
		if err == nil && parsed > 0 {
			id = uint(parsed)
		}
	}
	if id == 0 {
		defFY, err := invoiceDefaultFiscalYearWhenUnset(db, ctx)
		if err != nil {
			return finance_fiscal_year.FiscalYear{}, false
		}
		return defFY, true
	}
	var chosen finance_fiscal_year.FiscalYear
	if err := db.WithContext(ctx).First(&chosen, id).Error; err != nil {
		return finance_fiscal_year.FiscalYear{}, false
	}
	return chosen, true
}
