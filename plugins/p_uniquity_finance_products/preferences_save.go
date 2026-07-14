package p_uniquity_finance_products

import (
	"context"
	"strings"

	finance_accounts "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_accounts"
	"github.com/lariv-in/lago/components"
	"github.com/lariv-in/lago/getters"
	"github.com/lariv-in/lago/views"
	"gorm.io/gorm"
)

// OptionalUintValue returns 0 when p is nil.
func OptionalUintValue(p *uint) uint {
	if p == nil {
		return 0
	}
	return *p
}

// NormalizeOptionalUintFKValues converts empty or zero FK form values to nil so GORM writes NULL.
func NormalizeOptionalUintFKValues(values map[string]any, fields ...string) {
	for _, field := range fields {
		v, ok := values[field]
		if !ok {
			continue
		}
		switch t := v.(type) {
		case nil:
			values[field] = nil
		case string:
			if strings.TrimSpace(t) == "" {
				values[field] = nil
			}
		case uint:
			if t == 0 {
				values[field] = nil
			}
		case int:
			if t == 0 {
				values[field] = nil
			}
		case int64:
			if t == 0 {
				values[field] = nil
			}
		case float64:
			if t == 0 {
				values[field] = nil
			}
		}
	}
}

// OptionalPrefUintGetter reads an optional *uint preference from $in for foreign-key inputs.
func OptionalPrefUintGetter(field string) getters.Getter[uint] {
	return func(ctx context.Context) (uint, error) {
		ptr, err := getters.Key[*uint]("$in." + field)(ctx)
		if err != nil || ptr == nil {
			return 0, nil
		}
		return *ptr, nil
	}
}

// ProductPreferenceFormFields returns form field names owned by [ProductPreferences].
func ProductPreferenceFormFields() map[string]struct{} {
	return map[string]struct{}{
		ProductPrefInventoryAccountIDField: {},
		ProductPrefCostOfSalesAcctIDField:  {},
	}
}

// SplitProductPreferenceFormValues separates accounting preference fields from product preference fields.
func SplitProductPreferenceFormValues(values map[string]any) (accounting map[string]any, product map[string]any) {
	accounting = make(map[string]any, len(values))
	product = make(map[string]any, len(ProductPreferenceFormFields()))
	for key, value := range values {
		if _, isProduct := ProductPreferenceFormFields()[key]; isProduct {
			product[key] = value
			continue
		}
		accounting[key] = value
	}
	return accounting, product
}

// SaveAccountingAndProductPreferencesTx persists accounting and product singleton rows from a combined form post.
func SaveAccountingAndProductPreferencesTx(tx *gorm.DB, values map[string]any) error {
	accountingValues, productValues := SplitProductPreferenceFormValues(values)
	accountingRegular, accountingAssoc := views.SplitAssociationValues(accountingValues)
	productRegular, _ := views.SplitAssociationValues(productValues)

	NormalizeOptionalUintFKValues(
		productRegular,
		ProductPrefInventoryAccountIDField,
		ProductPrefCostOfSalesAcctIDField,
	)

	var accounting finance_accounts.AccountingPreferences
	if err := tx.FirstOrCreate(&accounting, finance_accounts.AccountingPreferences{Model: gorm.Model{ID: 1}}).Error; err != nil {
		return err
	}
	if len(accountingRegular) > 0 {
		if err := tx.Model(&accounting).Where("id = ?", accounting.ID).Updates(accountingRegular).Error; err != nil {
			return err
		}
	}
	if err := applyAccountingAssociationUpdates(tx, &accounting, accountingAssoc); err != nil {
		return err
	}

	var productPrefs ProductPreferences
	if err := tx.FirstOrCreate(&productPrefs, ProductPreferences{Model: gorm.Model{ID: 1}}).Error; err != nil {
		return err
	}
	if len(productRegular) == 0 {
		return nil
	}
	return tx.Model(&productPrefs).Where("id = ?", productPrefs.ID).Updates(productRegular).Error
}

func applyAccountingAssociationUpdates(tx *gorm.DB, record *finance_accounts.AccountingPreferences, associations map[string]components.AssociationIDs) error {
	if len(associations) == 0 {
		return nil
	}
	for field, ids := range associations {
		assoc := tx.Model(record).Association(field)
		if assoc.Error != nil {
			return assoc.Error
		}
		if len(ids.IDs) == 0 {
			if err := assoc.Clear(); err != nil {
				return err
			}
			continue
		}
		replaceIDs := make([]any, len(ids.IDs))
		for i, id := range ids.IDs {
			replaceIDs[i] = id
		}
		if err := assoc.Replace(replaceIDs); err != nil {
			return err
		}
	}
	return nil
}
