package p_uniquity_finance_invoices

// Query parameters and context keys for invoice hub DataTable column visibility (layer: views.LayerTableToggleColumns).
const (
	invoiceDraftColsParam      = "draft_invoice_cols"
	invoiceDraftColsCtxKey     = "finance_invoices.enabled_cols.draft"
	invoicePostedColsParam     = "posted_invoice_cols"
	invoicePostedColsCtxKey    = "finance_invoices.enabled_cols.posted"
	invoiceCancelledColsParam  = "cancelled_invoice_cols"
	invoiceCancelledColsCtxKey = "finance_invoices.enabled_cols.cancelled"
)
