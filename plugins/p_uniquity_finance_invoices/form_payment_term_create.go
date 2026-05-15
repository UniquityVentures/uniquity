package p_uniquity_finance_invoices

import (
	"fmt"
	"net/http"
	"time"

	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/views"
)

// paymentTermCreateFormPatcher creates the backing due-date or relative row and sets Type + BackingID for [PaymentTerm].
type paymentTermCreateFormPatcher struct{}

func (paymentTermCreateFormPatcher) Patch(_ views.View, r *http.Request, formData map[string]any, formErrors map[string]error) (map[string]any, map[string]error) {
	if len(formErrors) > 0 {
		return formData, formErrors
	}
	db, err := getters.DBFromContext(r.Context())
	if err != nil {
		formErrors["_form"] = err
		return formData, formErrors
	}
	typVal, ok := formData["Type"].(string)
	if !ok || typVal == "" {
		formErrors["Type"] = fmt.Errorf("payment term type is required")
		return formData, formErrors
	}
	switch typVal {
	case PaymentTermTypeDueDate:
		delete(formData, "Duration")
		dt, ok := formData["DueDatetime"].(time.Time)
		if !ok || dt.IsZero() {
			formErrors["DueDatetime"] = fmt.Errorf("datetime is required for this payment term kind")
			return formData, formErrors
		}
		row := PaymentTermDueDate{Datetime: dt}
		if err := db.Create(&row).Error; err != nil {
			formErrors["_form"] = err
			return formData, formErrors
		}
		formData["BackingID"] = row.ID
		formData["Type"] = PaymentTermTypeDueDate
	case PaymentTermTypeRelative:
		delete(formData, "DueDatetime")
		durPtr, ok := formData["Duration"].(*time.Duration)
		if !ok || durPtr == nil || *durPtr <= 0 {
			formErrors["Duration"] = fmt.Errorf("enter a positive duration (e.g. 720h for 30 days)")
			return formData, formErrors
		}
		row := PaymentTermRelative{Duration: *durPtr}
		if err := db.Create(&row).Error; err != nil {
			formErrors["_form"] = err
			return formData, formErrors
		}
		formData["BackingID"] = row.ID
		formData["Type"] = PaymentTermTypeRelative
	default:
		formErrors["Type"] = fmt.Errorf("invalid payment term type")
		return formData, formErrors
	}
	delete(formData, "DueDatetime")
	delete(formData, "Duration")
	return formData, formErrors
}
