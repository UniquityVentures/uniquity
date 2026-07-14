package p_uniquity_finance_invoices

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/components"
	"github.com/lariv-in/lago/getters"
	"github.com/lariv-in/lago/registry"
)

// invoiceListFilterFormTargetGetter preserves ?tab= when applying filters on the invoice hub.
func invoiceListFilterFormTargetGetter() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		base, err := lago.RoutePath("finance_invoices.DefaultRoute", nil)(ctx)
		if err != nil {
			return "", err
		}
		r, _ := ctx.Value("$request").(*http.Request)
		if r == nil {
			return base, nil
		}
		tab := strings.TrimSpace(r.URL.Query().Get("tab"))
		if tab == "" {
			return base, nil
		}
		sep := "?"
		if strings.Contains(base, "?") {
			sep = "&"
		}
		return base + sep + "tab=" + url.QueryEscape(tab), nil
	}
}

// InvoiceDateFilter holds list filter GET params (not persisted).
type InvoiceDateFilter struct {
	DatetimeFrom time.Time
	DatetimeTo   time.Time
}

func invoiceFilterGETTime(field string) getters.Getter[time.Time] {
	return func(ctx context.Context) (time.Time, error) {
		m, ok := ctx.Value("$get").(map[string]any)
		if !ok || m == nil {
			return time.Time{}, nil
		}
		v, ok := m[field]
		if !ok || v == nil {
			return time.Time{}, nil
		}
		switch t := v.(type) {
		case time.Time:
			return t, nil
		case string:
			s := strings.TrimSpace(t)
			if s == "" {
				return time.Time{}, nil
			}
			return parseInvoiceListFilterDatetime(s)
		default:
			return time.Time{}, nil
		}
	}
}

func pageEntriesInvoiceFilterPage() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "finance_invoices.InvoiceFilter", Value: &components.FormComponent[InvoiceDateFilter]{
			Attr: getters.FormBoostedGet(invoiceListFilterFormTargetGetter()),
			ChildrenInput: []components.PageInterface{
				&components.InputDatetime{
					Name:   "DatetimeFrom",
					Label:  "Invoice date from",
					Getter: invoiceFilterGETTime("DatetimeFrom"),
				},
				&components.InputDatetime{
					Name:   "DatetimeTo",
					Label:  "Invoice date to",
					Getter: invoiceFilterGETTime("DatetimeTo"),
				},
			},
			ChildrenAction: []components.PageInterface{
				&components.ContainerRow{
					Classes: "flex gap-2",
					Children: []components.PageInterface{
						&components.ButtonSubmit{Label: "Apply filters"},
						&components.ButtonClear{Label: "Clear"},
					},
				},
			},
		}},
	}
}
