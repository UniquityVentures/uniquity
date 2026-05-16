package p_uniquity_finance_invoices

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"

	finance_accounts "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_accounts"
	finance_fiscal_year "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_fiscal_year"
	"gorm.io/gorm"
)

// PostedInvoiceNumberTemplateData is the context for formatting posted invoice numbers (text/template).
type PostedInvoiceNumberTemplateData struct {
	FISCAL_CODE string
	YY          string
	YYYY        string
	POSTED_SEQ  int64
}

// FormatPostedInvoiceNumber executes AccountingPreferences.InvoiceNumberFormat as a Go text/template.
func FormatPostedInvoiceNumber(tx *gorm.DB, format string, invoiceDatetime time.Time, postedSeq int64) (string, error) {
	if format == "" {
		format = `INV-{{.YYYY}}-{{.POSTED_SEQ}}`
	}
	var fy finance_fiscal_year.FiscalYear
	err := tx.Where("starts_at <= ? AND ends_at >= ?", invoiceDatetime, invoiceDatetime).
		Order("starts_at DESC").First(&fy).Error
	if err != nil {
		err = tx.Where("is_active = ?", true).Order("starts_at DESC").First(&fy).Error
	}
	fiscalCode := ""
	if err == nil && fy.Code != "" {
		fiscalCode = fy.Code
	}
	data := PostedInvoiceNumberTemplateData{
		FISCAL_CODE: fiscalCode,
		YY:          invoiceDatetime.Format("06"),
		YYYY:        invoiceDatetime.Format("2006"),
		POSTED_SEQ:  postedSeq,
	}
	tpl, err := template.New("invoiceNumber").Parse(format)
	if err != nil {
		return "", fmt.Errorf("invoice number template: %w", err)
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("invoice number template execute: %w", err)
	}
	return buf.String(), nil
}

// NextPostedInvoiceSeq returns max(posted_invoices.id)+1 within tx (approximate sequence for new row).
func NextPostedInvoiceSeq(tx *gorm.DB) (int64, error) {
	var seq int64
	raw := tx.Raw("SELECT COALESCE(MAX(id), 0) FROM posted_invoices WHERE deleted_at IS NULL")
	if err := raw.Scan(&seq).Error; err != nil {
		return 0, err
	}
	return seq + 1, nil
}

// PostedInvoiceNumber resolves the final posted number: draft number if set, else template from prefs.
func PostedInvoiceNumber(tx *gorm.DB, draft *DraftInvoice, prefs finance_accounts.AccountingPreferences) (string, error) {
	if draft.Number != nil && strings.TrimSpace(*draft.Number) != "" {
		return strings.TrimSpace(*draft.Number), nil
	}
	seq, err := NextPostedInvoiceSeq(tx)
	if err != nil {
		return "", err
	}
	return FormatPostedInvoiceNumber(tx, prefs.InvoiceNumberFormat, draft.Datetime, seq)
}
