package p_uniquity_finance_invoices

import (
	"context"
	"strings"
	"time"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

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
			Attr: getters.FormBoostedGet(lamu.RoutePath("finance_invoices.DefaultRoute", nil)),
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
