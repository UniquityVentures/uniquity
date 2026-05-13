package p_uniquity_invoices

import (
	"net/http"
	"strings"
	"time"

	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/views"
)

type invoiceLineInvoiceIDFormPatcher struct{}

func (invoiceLineInvoiceIDFormPatcher) Patch(_ views.View, r *http.Request, formData map[string]any, formErrors map[string]error) (map[string]any, map[string]error) {
	inv, err := getters.Key[Invoice]("invoice")(r.Context())
	if err == nil && inv.ID != 0 {
		formData["InvoiceID"] = inv.ID
	}
	return formData, formErrors
}

type invoiceFormNormalizerPatcher struct{}

func (invoiceFormNormalizerPatcher) Patch(_ views.View, _ *http.Request, formData map[string]any, formErrors map[string]error) (map[string]any, map[string]error) {
	if v, ok := formData["Number"]; ok {
		if s, ok := v.(string); ok && strings.TrimSpace(s) == "" {
			formData["Number"] = nil
		}
	}
	if v, ok := formData["DueDate"]; ok {
		if t, ok := v.(time.Time); ok && t.IsZero() {
			formData["DueDate"] = nil
		}
	}
	return formData, formErrors
}
