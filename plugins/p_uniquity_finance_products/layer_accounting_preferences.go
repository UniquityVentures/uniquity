package p_uniquity_finance_products

import (
	"context"
	"net/http"

	"github.com/lariv-in/lago/getters"
	"github.com/lariv-in/lago/views"
)

type accountingPreferencesProductPrefsLayer struct {
	inner views.Layer
}

func patchAccountingPreferencesView(v *views.View) *views.View {
	return v.PatchLayer("finance_accounts.accounting_preferences", wrapAccountingPreferencesLayer)
}

func wrapAccountingPreferencesLayer(layer views.Layer) views.Layer {
	if _, ok := layer.(accountingPreferencesProductPrefsLayer); ok {
		return layer
	}
	return accountingPreferencesProductPrefsLayer{inner: layer}
}

func (m accountingPreferencesProductPrefsLayer) Next(view views.View, next http.Handler) http.Handler {
	mergeOnGet := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := mergeProductPreferencesIntoIn(r.Context())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
	return m.inner.Next(view, mergeOnGet)
}

func mergeProductPreferencesIntoIn(ctx context.Context) context.Context {
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
	prefs := LoadProductPreferences(db)
	for k, v := range getters.MapFromStruct(prefs) {
		inMap[k] = v
	}
	return context.WithValue(ctx, getters.ContextKeyIn, inMap)
}
