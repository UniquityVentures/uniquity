package p_uniquity_invoices

import (
	"net/http"
	"strconv"

	"github.com/UniquityVentures/lamu/views"
	"gorm.io/gorm"
)

// invoiceLinesScopedToInvoiceDetailURL filters invoice lines to the invoice id
// in the path (same {id} as [views.LayerDetail] for [Invoice]).
type invoiceLinesScopedToInvoiceDetailURL struct{}

func (invoiceLinesScopedToInvoiceDetailURL) Patch(_ views.View, r *http.Request, query gorm.ChainInterface[InvoiceLine]) gorm.ChainInterface[InvoiceLine] {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return query.Where("1 = 0")
	}
	return query.Where("invoice_id = ?", uint(id))
}
