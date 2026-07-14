package p_uniquity_finance_invoices

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/getters"
	"github.com/lariv-in/lago/views"
	"gorm.io/gorm"
)

// layerPostDraftInvoice runs [DraftInvoice.NewPosted] on POST and redirects to the posted detail page.
type layerPostDraftInvoice struct{}

func (layerPostDraftInvoice) Next(_ views.View, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			next.ServeHTTP(w, r)
			return
		}
		ctx := r.Context()
		draft, ok := ctx.Value("draft_invoice").(DraftInvoice)
		if !ok {
			slog.Error("layerPostDraftInvoice: draft missing from context")
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		db, err := getters.DBFromContext(ctx)
		if err != nil {
			slog.Error("layerPostDraftInvoice: db", "error", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		var posted *PostedInvoice
		err = db.Transaction(func(tx *gorm.DB) error {
			var d DraftInvoice
			if err := tx.First(&d, draft.ID).Error; err != nil {
				return err
			}
			p, err := (&d).NewPosted(tx, time.Now())
			posted = p
			return err
		})
		if err != nil {
			ctx = views.ContextWithErrorsAndValues(ctx, nil, map[string]error{
				"_global": fmt.Errorf("post invoice: %v", err),
			})
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		dest, err := lago.RoutePath("finance_invoices.PostedInvoiceDetailRoute", map[string]getters.Getter[any]{
			"id": getters.Any(getters.Static(posted.ID)),
		})(ctx)
		if err != nil {
			slog.Error("layerPostDraftInvoice: route", "error", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		views.HtmxRedirect(w, r, dest, http.StatusSeeOther)
	})
}

// layerCancelPostedInvoice runs [PostedInvoice.NewCancelled] on POST.
type layerCancelPostedInvoice struct{}

func (layerCancelPostedInvoice) Next(_ views.View, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			next.ServeHTTP(w, r)
			return
		}
		ctx := r.Context()
		posted, ok := ctx.Value("posted_invoice").(PostedInvoice)
		if !ok {
			slog.Error("layerCancelPostedInvoice: posted_invoice missing")
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		_ = r.ParseForm()
		reason := r.FormValue("Reason")
		db, err := getters.DBFromContext(ctx)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		var cancelled *CancelledInvoice
		err = db.Transaction(func(tx *gorm.DB) error {
			var p PostedInvoice
			if err := tx.First(&p, posted.ID).Error; err != nil {
				return err
			}
			c, err := (&p).NewCancelled(tx, reason, time.Now())
			cancelled = c
			return err
		})
		if err != nil {
			ctx = views.ContextWithErrorsAndValues(ctx, nil, map[string]error{
				"_global": fmt.Errorf("cancel invoice: %v", err),
			})
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		dest, err := lago.RoutePath("finance_invoices.CancelledInvoiceDetailRoute", map[string]getters.Getter[any]{
			"id": getters.Any(getters.Static(cancelled.ID)),
		})(ctx)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		views.HtmxRedirect(w, r, dest, http.StatusSeeOther)
	})
}

// layerNewDraftFromCancelled runs [CancelledInvoice.NewDraft] on POST.
type layerNewDraftFromCancelled struct{}

func (layerNewDraftFromCancelled) Next(_ views.View, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			next.ServeHTTP(w, r)
			return
		}
		ctx := r.Context()
		cinv, ok := ctx.Value("cancelled_invoice").(CancelledInvoice)
		if !ok {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		db, err := getters.DBFromContext(ctx)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		var draft *DraftInvoice
		err = db.Transaction(func(tx *gorm.DB) error {
			var c CancelledInvoice
			if err := tx.First(&c, cinv.ID).Error; err != nil {
				return err
			}
			d, err := (&c).NewDraft(tx)
			draft = d
			return err
		})
		if err != nil {
			ctx = views.ContextWithErrorsAndValues(ctx, nil, map[string]error{
				"_global": fmt.Errorf("new draft: %v", err),
			})
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		dest, err := lago.RoutePath("finance_invoices.DraftInvoiceDetailRoute", map[string]getters.Getter[any]{
			"id": getters.Any(getters.Static(draft.ID)),
		})(ctx)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		views.HtmxRedirect(w, r, dest, http.StatusSeeOther)
	})
}
