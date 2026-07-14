package p_uniquity_finance_invoices

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/lariv-in/lago/getters"
	"github.com/lariv-in/lago/views"
)

// paymentCreateQueryDefaultsLayer merges ?PostedInvoiceID= into $in on GET so selecting a
// posted invoice on the payment create form pre-fills Amount with the open balance.
type paymentCreateQueryDefaultsLayer struct{}

func paymentCreateDefaultsFromPostedInvoiceID(ctx context.Context, postedID uint) map[string]any {
	vals := map[string]any{"PostedInvoiceID": postedID}
	db, err := getters.DBFromContext(ctx)
	if err != nil {
		slog.Error("paymentCreateDefaultsFromPostedInvoiceID: db from context", "error", err)
		return vals
	}
	open, err := postedInvoiceOpenBalance(db, postedID)
	if err != nil {
		slog.Error("paymentCreateDefaultsFromPostedInvoiceID: open balance", "error", err, "postedID", postedID)
		return vals
	}
	if open.R != nil && open.R.Sign() > 0 {
		vals["Amount"] = open
	}
	return vals
}

func (paymentCreateQueryDefaultsLayer) Next(_ views.View, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}
		pid := r.URL.Query().Get("PostedInvoiceID")
		if pid == "" {
			next.ServeHTTP(w, r)
			return
		}
		id64, err := strconv.ParseUint(pid, 10, 32)
		if err != nil || id64 == 0 {
			next.ServeHTTP(w, r)
			return
		}
		vals := paymentCreateDefaultsFromPostedInvoiceID(r.Context(), uint(id64))
		ctx := views.ContextWithErrorsAndValues(r.Context(), vals, nil)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
