package p_uniquity_finance_invoices

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/views"
	"gorm.io/gorm"
)

// draftListDatetimeRange restricts draft invoice list by optional Datetime bounds (filter form: DatetimeFrom, DatetimeTo).
type draftListDatetimeRange struct{}

func (draftListDatetimeRange) Patch(_ views.View, r *http.Request, query gorm.ChainInterface[DraftInvoice]) gorm.ChainInterface[DraftInvoice] {
	ctx := r.Context()
	if t, ok := invoiceFilterTimeFromGet(ctx, "DatetimeFrom"); ok {
		query = query.Where("datetime >= ?", t)
	}
	if t, ok := invoiceFilterTimeFromGet(ctx, "DatetimeTo"); ok {
		query = query.Where("datetime <= ?", t)
	}
	return query
}

// draftListFiscalYearEnvironment restricts drafts to the fiscal year from environment cookie.
type draftListFiscalYearEnvironment struct{}

func (draftListFiscalYearEnvironment) Patch(_ views.View, r *http.Request, query gorm.ChainInterface[DraftInvoice]) gorm.ChainInterface[DraftInvoice] {
	ctx := r.Context()
	db, err := getters.DBFromContext(ctx)
	if err != nil {
		return query
	}
	fy, restrict := selectedInvoiceListFiscalYear(db, ctx)
	if !restrict {
		return query
	}
	return query.Where("datetime >= ? AND datetime <= ?", fy.Start, fy.End)
}

// draftListExcludePosted hides drafts that already have a posted invoice.
type draftListExcludePosted struct{}

func (draftListExcludePosted) Patch(_ views.View, _ *http.Request, query gorm.ChainInterface[DraftInvoice]) gorm.ChainInterface[DraftInvoice] {
	return query.Where(
		"NOT EXISTS (SELECT 1 FROM posted_invoices p WHERE p.draft_invoice_id = draft_invoices.id AND p.deleted_at IS NULL)",
	)
}

func invoiceFilterTimeFromGet(ctx context.Context, field string) (time.Time, bool) {
	m, ok := ctx.Value("$get").(map[string]any)
	if !ok || m == nil {
		return time.Time{}, false
	}
	v, ok := m[field]
	if !ok || v == nil {
		return time.Time{}, false
	}
	switch t := v.(type) {
	case time.Time:
		if t.IsZero() {
			return time.Time{}, false
		}
		return t, true
	case string:
		s := strings.TrimSpace(t)
		if s == "" {
			return time.Time{}, false
		}
		parsed, err := parseInvoiceListFilterDatetime(s)
		if err != nil {
			return time.Time{}, false
		}
		return parsed, true
	default:
		return time.Time{}, false
	}
}

func parseInvoiceListFilterDatetime(s string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	return time.ParseInLocation("2006-01-02T15:04", s, time.Local)
}

// postedListDatetimeRange filters posted invoices by Datetime.
type postedListDatetimeRange struct{}

func (postedListDatetimeRange) Patch(_ views.View, r *http.Request, query gorm.ChainInterface[PostedInvoice]) gorm.ChainInterface[PostedInvoice] {
	ctx := r.Context()
	if t, ok := invoiceFilterTimeFromGet(ctx, "DatetimeFrom"); ok {
		query = query.Where("datetime >= ?", t)
	}
	if t, ok := invoiceFilterTimeFromGet(ctx, "DatetimeTo"); ok {
		query = query.Where("datetime <= ?", t)
	}
	return query
}

type postedListFiscalYearEnvironment struct{}

func (postedListFiscalYearEnvironment) Patch(_ views.View, r *http.Request, query gorm.ChainInterface[PostedInvoice]) gorm.ChainInterface[PostedInvoice] {
	ctx := r.Context()
	db, err := getters.DBFromContext(ctx)
	if err != nil {
		return query
	}
	fy, restrict := selectedInvoiceListFiscalYear(db, ctx)
	if !restrict {
		return query
	}
	return query.Where("datetime >= ? AND datetime <= ?", fy.Start, fy.End)
}

// postedListExcludeCancelled hides posted rows that have been cancelled.
type postedListExcludeCancelled struct{}

func (postedListExcludeCancelled) Patch(_ views.View, _ *http.Request, query gorm.ChainInterface[PostedInvoice]) gorm.ChainInterface[PostedInvoice] {
	return query.Where(
		"NOT EXISTS (SELECT 1 FROM cancelled_invoices c WHERE c.posted_invoice_id = posted_invoices.id AND c.deleted_at IS NULL)",
	)
}

// cancelledListDatetimeRange filters cancelled invoices by original Datetime.
type cancelledListDatetimeRange struct{}

func (cancelledListDatetimeRange) Patch(_ views.View, r *http.Request, query gorm.ChainInterface[CancelledInvoice]) gorm.ChainInterface[CancelledInvoice] {
	ctx := r.Context()
	if t, ok := invoiceFilterTimeFromGet(ctx, "DatetimeFrom"); ok {
		query = query.Where("datetime >= ?", t)
	}
	if t, ok := invoiceFilterTimeFromGet(ctx, "DatetimeTo"); ok {
		query = query.Where("datetime <= ?", t)
	}
	return query
}

type cancelledListFiscalYearEnvironment struct{}

func (cancelledListFiscalYearEnvironment) Patch(_ views.View, r *http.Request, query gorm.ChainInterface[CancelledInvoice]) gorm.ChainInterface[CancelledInvoice] {
	ctx := r.Context()
	db, err := getters.DBFromContext(ctx)
	if err != nil {
		return query
	}
	fy, restrict := selectedInvoiceListFiscalYear(db, ctx)
	if !restrict {
		return query
	}
	return query.Where("datetime >= ? AND datetime <= ?", fy.Start, fy.End)
}
