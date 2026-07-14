package p_uniquity_finance_invoices

import (
	"fmt"
	"net/http"

	finance_taxes "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_taxes"
	"github.com/lariv-in/lago/components"
	"github.com/lariv-in/lago/getters"
	"github.com/lariv-in/lago/views"
	"gorm.io/gorm"
)

// paymentCreateFormPatcher loads payment taxes for journal posting before LayerCreate applies M2M.
type paymentCreateFormPatcher struct{}

func (paymentCreateFormPatcher) Patch(_ views.View, r *http.Request, formData map[string]any, formErrors map[string]error) (map[string]any, map[string]error) {
	if len(formErrors) > 0 {
		return formData, formErrors
	}
	db, err := getters.DBFromContext(r.Context())
	if err != nil {
		formErrors["_form"] = err
		return formData, formErrors
	}
	taxes, err := loadPaymentTaxesFromForm(db, formData)
	if err != nil {
		formErrors["Taxes"] = err
		return formData, formErrors
	}
	if err := validatePaymentTaxes(taxes); err != nil {
		formErrors["Taxes"] = err
		return formData, formErrors
	}
	*r = *r.WithContext(contextWithPaymentCreateTaxes(r.Context(), taxes))
	return formData, formErrors
}

func loadPaymentTaxesFromForm(db *gorm.DB, formData map[string]any) ([]finance_taxes.Tax, error) {
	ids := associationIDsFromForm(formData, "Taxes")
	if len(ids) == 0 {
		return nil, nil
	}
	var taxes []finance_taxes.Tax
	if err := db.Where("id IN ?", ids).Find(&taxes).Error; err != nil {
		return nil, err
	}
	if len(taxes) != len(ids) {
		return nil, fmt.Errorf("one or more selected taxes are invalid")
	}
	return taxes, nil
}

func associationIDsFromForm(formData map[string]any, field string) []uint {
	raw, ok := formData[field]
	if !ok {
		return nil
	}
	switch v := raw.(type) {
	case components.AssociationIDs:
		return v.IDs
	case *components.AssociationIDs:
		if v == nil {
			return nil
		}
		return v.IDs
	default:
		return nil
	}
}
