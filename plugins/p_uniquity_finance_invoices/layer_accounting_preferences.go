package p_uniquity_finance_invoices

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/views"
	finance_products "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_products"
	"gorm.io/gorm"
)

type accountingPreferencesInvoicePrefsLayer struct {
	inner views.Layer
}

func patchAccountingPreferencesView(v *views.View) *views.View {
	return v.PatchLayer("finance_accounts.accounting_preferences", wrapAccountingPreferencesLayer)
}

func wrapAccountingPreferencesLayer(layer views.Layer) views.Layer {
	if _, ok := layer.(accountingPreferencesInvoicePrefsLayer); ok {
		return layer
	}
	return accountingPreferencesInvoicePrefsLayer{inner: layer}
}

func (m accountingPreferencesInvoicePrefsLayer) Next(view views.View, next http.Handler) http.Handler {
	mergeOnGet := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := mergeInvoicePreferencesIntoIn(r.Context())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
	innerHandler := m.inner.Next(view, mergeOnGet)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			innerHandler.ServeHTTP(w, r)
			return
		}
		m.handlePost(view, w, r)
	})
}

func mergeInvoicePreferencesIntoIn(ctx context.Context) context.Context {
	db, err := getters.DBFromContext(ctx)
	if err != nil {
		return ctx
	}
	inMap, ok := ctx.Value(getters.ContextKeyIn).(map[string]any)
	if !ok {
		inMap = map[string]any{}
	} else {
		cloned := make(map[string]any, len(inMap))
		for k, v := range inMap {
			cloned[k] = v
		}
		inMap = cloned
	}
	for k, v := range getters.MapFromStruct(LoadInvoicePreferences(db)) {
		inMap[k] = v
	}
	for k, v := range getters.MapFromStruct(LoadPaymentPreferences(db)) {
		inMap[k] = v
	}
	return context.WithValue(ctx, getters.ContextKeyIn, inMap)
}

func (m accountingPreferencesInvoicePrefsLayer) handlePost(view views.View, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db, dberr := getters.DBFromContext(ctx)
	if dberr != nil {
		slog.Error("finance_invoices.accounting_preferences_invoice_prefs: db from context", "error", dberr)
		ctx = views.ContextWithErrorsAndValues(ctx, nil, map[string]error{"_global": dberr})
		view.RenderPage(w, r.WithContext(ctx))
		return
	}

	values, fieldErrors, err := view.ParseForm(w, r)
	if err != nil {
		slog.Error("finance_invoices.accounting_preferences_invoice_prefs: parse form", "error", err)
		ctx = views.ContextWithErrorsAndValues(ctx, values, map[string]error{"_form": err})
		view.RenderPage(w, r.WithContext(ctx))
		return
	}
	if len(fieldErrors) != 0 {
		ctx = views.ContextWithErrorsAndValues(ctx, values, fieldErrors)
		view.RenderPage(w, r.WithContext(ctx))
		return
	}

	accountingProductValues, invoiceValues, paymentValues := splitPreferenceFormValues(values)
	invoiceRegular, _ := views.SplitAssociationValues(invoiceValues)
	paymentRegular, _ := views.SplitAssociationValues(paymentValues)
	finance_products.NormalizeOptionalUintFKValues(invoiceRegular,
		InvoicePrefAccountReceivableIDField,
		InvoicePrefAccountRevenueIDField,
		InvoicePrefAccountTaxPayableIDField,
		InvoicePrefJournalIDField,
	)

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := finance_products.SaveAccountingAndProductPreferencesTx(tx, accountingProductValues); err != nil {
			return err
		}
		var invoicePrefs InvoicePreferences
		if err := tx.FirstOrCreate(&invoicePrefs, InvoicePreferences{Model: gorm.Model{ID: 1}}).Error; err != nil {
			return err
		}
		if len(invoiceRegular) > 0 {
			if err := tx.Model(&invoicePrefs).Where("id = ?", invoicePrefs.ID).Updates(invoiceRegular).Error; err != nil {
				return err
			}
		}
		paymentPrefs := LoadPaymentPreferences(tx)
		return savePaymentPreferencesTx(tx, &paymentPrefs, paymentRegular)
	})
	if err != nil {
		slog.Error("finance_invoices.accounting_preferences_invoice_prefs: transaction", "error", err)
		fieldErrors["_form"] = fmt.Errorf("%v", err)
		ctx = views.ContextWithErrorsAndValues(ctx, values, fieldErrors)
		view.RenderPage(w, r.WithContext(ctx))
		return
	}

	successURL, err := lamu.RoutePath("finance_accounts.AccountingPreferencesRoute", nil)(ctx)
	if err != nil {
		fieldErrors["_form"] = err
		ctx = views.ContextWithErrorsAndValues(ctx, values, fieldErrors)
		view.RenderPage(w, r.WithContext(ctx))
		return
	}
	views.HtmxRedirect(w, r, successURL, http.StatusSeeOther)
}

func splitPreferenceFormValues(values map[string]any) (accountingProduct, invoice, payment map[string]any) {
	accountingProduct = make(map[string]any, len(values))
	invoice = make(map[string]any, len(InvoicePreferenceFormFields()))
	payment = make(map[string]any, len(PaymentPreferenceFormFields()))
	for key, value := range values {
		if _, isInvoice := InvoicePreferenceFormFields()[key]; isInvoice {
			invoice[key] = value
			continue
		}
		if _, isPayment := PaymentPreferenceFormFields()[key]; isPayment {
			payment[key] = value
			continue
		}
		accountingProduct[key] = value
	}
	return accountingProduct, invoice, payment
}

func savePaymentPreferencesTx(tx *gorm.DB, prefs *PaymentPreferences, values map[string]any) error {
	if len(values) == 0 {
		return nil
	}
	finance_products.NormalizeOptionalUintFKValues(values, PaymentPrefAccountIDField)
	v, ok := values[PaymentPrefAccountIDField]
	if !ok {
		return nil
	}
	return paymentPreferencesDB(tx).Where("id = ?", prefs.ID).Update("payment_account_id", v).Error
}
