package p_uniquity_finance_invoices

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/fields"
	"github.com/UniquityVentures/lamu/getters"
	finance_taxes "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_taxes"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// InvoiceHeaderTaxAmountRow is one document-level tax line in [FieldInvoiceLinesSummary].
type InvoiceHeaderTaxAmountRow struct {
	Label  string
	Amount string
}

// InvoiceLinesSummary is the footer under the lines table (matches the draft line editor).
type InvoiceLinesSummary struct {
	LinesSubtotal string
	TaxRows       []InvoiceHeaderTaxAmountRow
	GrandTotal    string
}

func finishInvoiceLinesSummary(untaxedSubtotal, linesGrossTotal fields.DecimalSix, headerTaxes []finance_taxes.Tax, lineCount int) InvoiceLinesSummary {
	var rows []InvoiceHeaderTaxAmountRow
	var headerSum fields.DecimalSix
	for _, t := range headerTaxes {
		amt := taxAmountOnBase(untaxedSubtotal, sumTaxPercents([]finance_taxes.Tax{t}))
		label := t.Name
		if strings.TrimSpace(label) == "" {
			label = fmt.Sprintf("Tax #%d", t.ID)
		}
		rows = append(rows, InvoiceHeaderTaxAmountRow{Label: label, Amount: decimalSixDisplay(amt)})
		headerSum = decSum(headerSum, amt)
	}
	grand := decSum(linesGrossTotal, headerSum)

	linesSub := "—"
	if lineCount > 0 {
		linesSub = decimalSixDisplay(linesGrossTotal)
	}
	grandStr := decimalSixDisplay(grand)
	if lineCount == 0 && len(headerTaxes) == 0 && decimalIsZero(grand) {
		grandStr = "—"
	}
	return InvoiceLinesSummary{LinesSubtotal: linesSub, TaxRows: rows, GrandTotal: grandStr}
}

func decimalIsZero(d fields.DecimalSix) bool {
	return d.R == nil || d.R.Sign() == 0
}

func invoiceLinesSummaryFromDraftLines(lines []DraftInvoiceLine, headerTaxes []finance_taxes.Tax) InvoiceLinesSummary {
	var untaxedSubtotal, linesGrossTotal fields.DecimalSix
	for _, ln := range lines {
		u, _, tot := invoiceLineAmountBreakdown(ln.Quantity, ln.Rate, ln.Taxes)
		untaxedSubtotal = decSum(untaxedSubtotal, u)
		linesGrossTotal = decSum(linesGrossTotal, tot)
	}
	return finishInvoiceLinesSummary(untaxedSubtotal, linesGrossTotal, headerTaxes, len(lines))
}

func invoiceLinesSummaryFromPostedLines(lines []PostedInvoiceLine, headerTaxes []finance_taxes.Tax) InvoiceLinesSummary {
	var untaxedSubtotal, linesGrossTotal fields.DecimalSix
	for _, ln := range lines {
		u, _, tot := invoiceLineAmountBreakdown(ln.Quantity, ln.Rate, ln.Taxes)
		untaxedSubtotal = decSum(untaxedSubtotal, u)
		linesGrossTotal = decSum(linesGrossTotal, tot)
	}
	return finishInvoiceLinesSummary(untaxedSubtotal, linesGrossTotal, headerTaxes, len(lines))
}

func invoiceLinesSummaryFromCancelledLines(lines []CancelledInvoiceLine, headerTaxes []finance_taxes.Tax) InvoiceLinesSummary {
	var untaxedSubtotal, linesGrossTotal fields.DecimalSix
	for _, ln := range lines {
		u, _, tot := invoiceLineAmountBreakdown(ln.Quantity, ln.Rate, ln.Taxes)
		untaxedSubtotal = decSum(untaxedSubtotal, u)
		linesGrossTotal = decSum(linesGrossTotal, tot)
	}
	return finishInvoiceLinesSummary(untaxedSubtotal, linesGrossTotal, headerTaxes, len(lines))
}

// InvoiceLineDisplay is one read-only row for [FieldInvoiceLines].
type InvoiceLineDisplay struct {
	Product       string
	Quantity      string
	Rate          string
	LineTaxes     string
	UntaxedAmount string
	TaxedAmount   string
	LineTotal     string
}

func invoiceLineTaxesLabel(taxes []finance_taxes.Tax) string {
	if len(taxes) == 0 {
		return "—"
	}
	names := make([]string, 0, len(taxes))
	for _, t := range taxes {
		if t.Name != "" {
			names = append(names, t.Name)
		} else {
			names = append(names, fmt.Sprintf("#%d", t.ID))
		}
	}
	return strings.Join(names, ", ")
}

// decimalSixDisplay trims trailing zeros like the draft line editor's formatDec.
func decimalSixDisplay(d fields.DecimalSix) string {
	s := d.String()
	if !strings.Contains(s, ".") {
		return s
	}
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	if s == "" {
		return "0"
	}
	return s
}

// FieldInvoiceLines renders stored invoice lines as a read-only table (mirrors draft editor columns).
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
			Td(ColSpan("7"), Class("text-center opacity-50 py-4"), Text("No lines")),
		))
	} else {
		for _, r := range rows {
			tbody = append(tbody, Tr(
				Td(Class("whitespace-nowrap max-w-md"), Text(r.Product)),
				Td(Class("whitespace-nowrap text-end tabular-nums"), Text(r.Quantity)),
				Td(Class("whitespace-nowrap text-end tabular-nums"), Text(r.Rate)),
				Td(Class("min-w-[10rem] max-w-md text-sm"), Text(r.LineTaxes)),
				Td(Class("whitespace-nowrap text-end tabular-nums"), Text(r.UntaxedAmount)),
				Td(Class("whitespace-nowrap text-end tabular-nums"), Text(r.TaxedAmount)),
				Td(Class("whitespace-nowrap text-end tabular-nums font-medium"), Text(r.LineTotal)),
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
					Th(Class("whitespace-nowrap min-w-[10rem]"), Text("Line taxes")),
					Th(Class("whitespace-nowrap min-w-[7rem] text-end"), Text("Untaxed amount")),
					Th(Class("whitespace-nowrap min-w-[7rem] text-end"), Text("Taxed amount")),
					Th(Class("whitespace-nowrap min-w-[7rem] text-end"), Text("Line total")),
				)),
				TBody(tbody...),
			),
		),
	)
}

// FieldInvoiceLinesSummary renders Lines subtotal, document-level taxes, and total (matches draft editor footer).
type FieldInvoiceLinesSummary struct {
	components.Page
	Getter  getters.Getter[InvoiceLinesSummary]
	Classes string
}

func (e FieldInvoiceLinesSummary) GetKey() string { return e.Key }

func (e FieldInvoiceLinesSummary) GetRoles() []string { return e.Roles }

func (e FieldInvoiceLinesSummary) Build(ctx context.Context) Node {
	var s InvoiceLinesSummary
	if e.Getter != nil {
		var err error
		s, err = e.Getter(ctx)
		if err != nil {
			slog.Error("FieldInvoiceLinesSummary getter failed", "error", err, "key", e.Key)
			return components.ContainerError{Error: getters.Static(err)}.Build(ctx)
		}
	}
	outerClass := "mt-3 w-full"
	if c := strings.TrimSpace(e.Classes); c != "" {
		outerClass += " " + c
	}
	rowNodes := []Node{
		Div(Class("grid grid-cols-[1fr_auto] gap-x-4 items-center px-4 py-3"),
			Div(Class("text-sm font-bold min-w-0 truncate"), Text("Lines subtotal")),
			Div(Class("text-sm tabular-nums text-end font-semibold shrink-0 min-w-[7rem]"), Text(s.LinesSubtotal)),
		),
	}
	for _, row := range s.TaxRows {
		rowNodes = append(rowNodes, Div(Class("grid grid-cols-[1fr_auto] gap-x-4 items-center px-4 py-3"),
			Div(Class("text-sm font-bold min-w-0 truncate"), Text(row.Label)),
			Div(Class("text-sm tabular-nums text-end font-semibold shrink-0 min-w-[7rem]"), Text(row.Amount)),
		))
	}
	rowNodes = append(rowNodes, Div(Class("grid grid-cols-[1fr_auto] gap-x-4 items-center px-4 py-3 bg-base-200/60"),
		Div(Class("text-sm font-bold min-w-0 truncate"), Text("Total")),
		Div(Class("text-sm tabular-nums text-end font-bold shrink-0 min-w-[7rem]"), Text(s.GrandTotal)),
	))
	return Div(Class(outerClass),
		Div(append([]Node{Class("rounded-box border border-base-300 bg-base-100 overflow-hidden divide-y divide-base-300")}, rowNodes...)...),
	)
}
