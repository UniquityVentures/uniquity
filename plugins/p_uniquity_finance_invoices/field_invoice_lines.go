package p_uniquity_finance_invoices

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// InvoiceLineDisplay is one read-only row for [FieldInvoiceLines].
type InvoiceLineDisplay struct {
	Product  string
	Quantity string
	Rate     string
}

// FieldInvoiceLines renders stored invoice lines as a read-only table (product, quantity, rate).
type FieldInvoiceLines struct {
	components.Page
	Getter  getters.Getter[[]InvoiceLineDisplay]
	Classes string
}

func (e FieldInvoiceLines) GetKey() string { return e.Key }

func (e FieldInvoiceLines) GetRoles() []string { return e.Roles }

func (e FieldInvoiceLines) Build(ctx context.Context) Node {
	var rows []InvoiceLineDisplay
	if e.Getter != nil {
		r, err := e.Getter(ctx)
		if err != nil {
			slog.Error("FieldInvoiceLines getter failed", "error", err, "key", e.Key)
			return components.ContainerError{Error: getters.Static(err)}.Build(ctx)
		}
		rows = r
	}
	wrap := fmt.Sprintf("w-full %s", e.Classes)
	var tbody []Node
	if len(rows) == 0 {
		tbody = append(tbody, Tr(
			Td(ColSpan("3"), Class("text-center opacity-50 py-4"), Text("No lines")),
		))
	} else {
		for _, r := range rows {
			tbody = append(tbody, Tr(
				Td(Class("whitespace-nowrap max-w-md"), Text(r.Product)),
				Td(Class("whitespace-nowrap text-end tabular-nums"), Text(r.Quantity)),
				Td(Class("whitespace-nowrap text-end tabular-nums"), Text(r.Rate)),
			))
		}
	}
	return Div(Class(wrap),
		Div(Class("overflow-x-auto rounded-box border border-base-300 bg-base-100"),
			Table(Class("table table-sm w-full"),
				THead(Tr(
					Th(Class("whitespace-nowrap min-w-[12rem]"), Text("Product")),
					Th(Class("whitespace-nowrap w-32 text-end"), Text("Quantity")),
					Th(Class("whitespace-nowrap w-32 text-end"), Text("Rate")),
				)),
				TBody(tbody...),
			),
		),
	)
}
