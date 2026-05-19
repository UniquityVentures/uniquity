package p_uniquity_finance_invoices

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"text/template"

	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/views"
	finance_accounts "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_accounts"
	"github.com/francescoalemanno/gotypst"
	"gorm.io/gorm"
)

func draftInvoiceFromContext(ctx context.Context) (DraftInvoice, bool) {
	switch v := ctx.Value("draft_invoice").(type) {
	case DraftInvoice:
		return v, true
	case *DraftInvoice:
		if v == nil {
			return DraftInvoice{}, false
		}
		return *v, true
	default:
		return DraftInvoice{}, false
	}
}

func cancelledInvoiceFromContext(ctx context.Context) (CancelledInvoice, bool) {
	switch v := ctx.Value("cancelled_invoice").(type) {
	case CancelledInvoice:
		return v, true
	case *CancelledInvoice:
		if v == nil {
			return CancelledInvoice{}, false
		}
		return *v, true
	default:
		return CancelledInvoice{}, false
	}
}

func serveInvoicePDFFromPrefs(w http.ResponseWriter, db *gorm.DB, templateRoot map[string]any, filenameBase string, logPrefix string, entityID uint) {
	prefs := finance_accounts.LoadAccountingPreferences(db)
	tmplSrc := strings.TrimSpace(prefs.InvoicePDFTemplate)
	if tmplSrc == "" {
		http.Error(w, "Configure the invoice PDF template in Accounting preferences first.", http.StatusBadRequest)
		return
	}
	tmpl, err := template.New("invoice_pdf").Parse(tmplSrc)
	if err != nil {
		slog.Error(logPrefix+": parse preferences template", "error", err)
		http.Error(w, "Invalid invoice PDF template in preferences.", http.StatusInternalServerError)
		return
	}
	var typstBuf bytes.Buffer
	if err := tmpl.Execute(&typstBuf, templateRoot); err != nil {
		slog.Error(logPrefix+": execute template", "error", err, "id", entityID)
		http.Error(w, "Rendering invoice PDF template failed.", http.StatusInternalServerError)
		return
	}
	pdfBytes, err := gotypst.PDF(typstBuf.Bytes())
	if err != nil {
		slog.Error(logPrefix+": gotypst", "error", err, "id", entityID)
		http.Error(w, "PDF compilation failed.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.pdf"`, filenameBase))
	if _, err := w.Write(pdfBytes); err != nil {
		slog.Error(logPrefix+": write", "error", err)
	}
}

func draftInvoicePdfHandler(_ *views.View) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		draft, ok := draftInvoiceFromContext(ctx)
		if !ok {
			slog.Error("draft_invoice_pdf: missing draft_invoice in context")
			http.Error(w, "Draft invoice not found", http.StatusInternalServerError)
			return
		}
		db, err := getters.DBFromContext(ctx)
		if err != nil {
			slog.Error("draft_invoice_pdf: db", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		base := fmt.Sprintf("draft-invoice-%d", draft.ID)
		if draft.Number != nil && strings.TrimSpace(*draft.Number) != "" {
			base = sanitizeInvoicePdfFilenameBase(*draft.Number)
		}
		serveInvoicePDFFromPrefs(w, db, getters.MapFromStruct(draft), base, "draft_invoice_pdf", draft.ID)
	})
}

func cancelledInvoicePdfHandler(_ *views.View) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		inv, ok := cancelledInvoiceFromContext(ctx)
		if !ok {
			slog.Error("cancelled_invoice_pdf: missing cancelled_invoice in context")
			http.Error(w, "Cancelled invoice not found", http.StatusInternalServerError)
			return
		}
		db, err := getters.DBFromContext(ctx)
		if err != nil {
			slog.Error("cancelled_invoice_pdf: db", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		base := fmt.Sprintf("cancelled-invoice-%d", inv.ID)
		if strings.TrimSpace(inv.Number) != "" {
			base = sanitizeInvoicePdfFilenameBase(inv.Number)
		}
		serveInvoicePDFFromPrefs(w, db, getters.MapFromStruct(inv), base, "cancelled_invoice_pdf", inv.ID)
	})
}

func sanitizeInvoicePdfFilenameBase(s string) string {
	s = strings.TrimSpace(s)
	for _, ch := range []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"} {
		s = strings.ReplaceAll(s, ch, "-")
	}
	if s == "" {
		return "invoice"
	}
	return s
}
