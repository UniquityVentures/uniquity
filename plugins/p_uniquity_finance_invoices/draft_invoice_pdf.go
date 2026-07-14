package p_uniquity_finance_invoices

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"text/template"

	finance_accounts "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_accounts"
	"github.com/francescoalemanno/gotypst"
	"github.com/lariv-in/lago/getters"
	"github.com/lariv-in/lago/views"
	"gorm.io/gorm"
)

var postedInvoicePdfPreload = []string{"Customer", "PaymentTerm", "Taxes", "Lines", "Lines.Product", "Lines.Taxes"}

func servePostedInvoicePDF(w http.ResponseWriter, db *gorm.DB, posted PostedInvoice, logPrefix string) {
	base := fmt.Sprintf("invoice-%d", posted.ID)
	if strings.TrimSpace(posted.Number) != "" {
		base = sanitizeInvoicePdfFilenameBase(posted.Number)
	}
	serveInvoicePDFFromPrefs(w, db, getters.MapFromStruct(posted), base, logPrefix, posted.ID)
}

func serveInvoicePDFFromPrefs(w http.ResponseWriter, db *gorm.DB, templateRoot map[string]any, filenameBase string, logPrefix string, entityID uint) {
	var payments []Payment
	if logPrefix != "draft_invoice_pdf" {
		var postedID uint
		if logPrefix == "cancelled_invoice_pdf" {
			var cancelled CancelledInvoice
			if err := db.First(&cancelled, entityID).Error; err == nil {
				postedID = cancelled.PostedInvoiceID
			}
		} else {
			postedID = entityID
		}
		if postedID != 0 {
			if err := db.Where("posted_invoice_id = ?", postedID).Order("datetime ASC").Find(&payments).Error; err != nil {
				slog.Error(logPrefix+": load payments", "error", err, "posted_id", postedID)
			}
		}
	}
	templateRoot["Payments"] = payments

	prefs := finance_accounts.LoadAccountingPreferences(db)
	tmplSrc := strings.TrimSpace(prefs.InvoicePDFTemplate)
	if tmplSrc == "" {
		http.Error(w, "Configure the invoice PDF template in Accounting preferences first.", http.StatusBadRequest)
		return
	}
	tmpl, err := template.New("invoice_pdf").Funcs(invoicePDFTemplateFuncs()).Parse(tmplSrc)
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
		draft, err := views.GetValueFromContext[string, DraftInvoice](ctx, "draft_invoice")
		if err != nil {
			slog.Error("draft_invoice_pdf: missing draft_invoice in context", "error", err)
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
		inv, err := views.GetValueFromContext[string, CancelledInvoice](ctx, "cancelled_invoice")
		if err != nil {
			slog.Error("cancelled_invoice_pdf: missing cancelled_invoice in context", "error", err)
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

func postedInvoicePdfHandler(_ *views.View) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		posted, err := views.GetValueFromContext[string, PostedInvoice](ctx, "posted_invoice")
		if err != nil {
			slog.Error("posted_invoice_pdf: missing posted_invoice in context", "error", err)
			http.Error(w, "Posted invoice not found", http.StatusInternalServerError)
			return
		}
		db, err := getters.DBFromContext(ctx)
		if err != nil {
			slog.Error("posted_invoice_pdf: db", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		servePostedInvoicePDF(w, db, posted, "posted_invoice_pdf")
	})
}

func paidInvoicePdfHandler(_ *views.View) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		paid, err := views.GetValueFromContext[string, PaidInvoice](ctx, "paid_invoice")
		if err != nil {
			slog.Error("paid_invoice_pdf: missing paid_invoice in context", "error", err)
			http.Error(w, "Paid invoice not found", http.StatusInternalServerError)
			return
		}
		db, err := getters.DBFromContext(ctx)
		if err != nil {
			slog.Error("paid_invoice_pdf: db", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		servePostedInvoicePDF(w, db, paid.PostedInvoice, "paid_invoice_pdf")
	})
}

func partiallyPaidInvoicePdfHandler(_ *views.View) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		partial, err := views.GetValueFromContext[string, PartiallyPaidInvoice](ctx, "partially_paid_invoice")
		if err != nil {
			slog.Error("partially_paid_invoice_pdf: missing partially_paid_invoice in context", "error", err)
			http.Error(w, "Partially paid invoice not found", http.StatusInternalServerError)
			return
		}
		db, err := getters.DBFromContext(ctx)
		if err != nil {
			slog.Error("partially_paid_invoice_pdf: db", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		servePostedInvoicePDF(w, db, partial.PostedInvoice, "partially_paid_invoice_pdf")
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
