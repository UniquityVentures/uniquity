package p_uniquity_finance_invoices

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	finance_taxes "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_taxes"
	"github.com/lariv-in/lago/components"
	"github.com/lariv-in/lago/fields"
	"github.com/lariv-in/lago/getters"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// InvoiceHeaderTaxAmountRow is one document-level tax line in [FieldInvoiceLinesSummary].
type InvoiceHeaderTaxAmountRow struct {
	Label         string
	Amount        string
	IsWithholding bool
}

// InvoiceLinesSummary is the footer under the lines table (matches the draft line editor).
type InvoiceLinesSummary struct {
	LinesSubtotal   string
	TaxRows         []InvoiceHeaderTaxAmountRow
	WithholdingRows []InvoiceHeaderTaxAmountRow
	GrandTotal      string
}

func mergeInvoiceLineTaxIDs(into map[uint]struct{}, taxes []finance_taxes.Tax) {
	for _, t := range taxes {
		if t.ID != 0 {
			into[t.ID] = struct{}{}
		}
	}
}

// documentLevelHeaderTaxes returns invoice-level taxes that are not already applied on any line.
func documentLevelHeaderTaxes(header []finance_taxes.Tax, lineTaxIDs map[uint]struct{}) []finance_taxes.Tax {
	var out []finance_taxes.Tax
	for _, t := range header {
		if t.ID == 0 {
			continue
		}
		if _, onLine := lineTaxIDs[t.ID]; onLine {
			continue
		}
		out = append(out, t)
	}
	return out
}

func accumulateInvoiceLineTotals(lines []DraftInvoiceLine) (invoiceLinesTotals, map[uint]struct{}) {
	var totals invoiceLinesTotals
	lineTaxIDs := map[uint]struct{}{}
	for _, ln := range lines {
		u, lev, wh, _ := invoiceLineAmountBreakdown(ln.Quantity, ln.Rate, ln.Taxes)
		totals.UntaxedSubtotal = decSum(totals.UntaxedSubtotal, u)
		totals.LinesLevied = decSum(totals.LinesLevied, lev)
		totals.LinesWithholding = decSum(totals.LinesWithholding, wh)
		mergeInvoiceLineTaxIDs(lineTaxIDs, ln.Taxes)
	}
	return totals, lineTaxIDs
}

func accumulatePostedInvoiceLineTotals(lines []PostedInvoiceLine) (invoiceLinesTotals, map[uint]struct{}) {
	var totals invoiceLinesTotals
	lineTaxIDs := map[uint]struct{}{}
	for _, ln := range lines {
		u, lev, wh, _ := invoiceLineAmountBreakdown(ln.Quantity, ln.Rate, ln.Taxes)
		totals.UntaxedSubtotal = decSum(totals.UntaxedSubtotal, u)
		totals.LinesLevied = decSum(totals.LinesLevied, lev)
		totals.LinesWithholding = decSum(totals.LinesWithholding, wh)
		mergeInvoiceLineTaxIDs(lineTaxIDs, ln.Taxes)
	}
	return totals, lineTaxIDs
}

func accumulateCancelledInvoiceLineTotals(lines []CancelledInvoiceLine) (invoiceLinesTotals, map[uint]struct{}) {
	var totals invoiceLinesTotals
	lineTaxIDs := map[uint]struct{}{}
	for _, ln := range lines {
		u, lev, wh, _ := invoiceLineAmountBreakdown(ln.Quantity, ln.Rate, ln.Taxes)
		totals.UntaxedSubtotal = decSum(totals.UntaxedSubtotal, u)
		totals.LinesLevied = decSum(totals.LinesLevied, lev)
		totals.LinesWithholding = decSum(totals.LinesWithholding, wh)
		mergeInvoiceLineTaxIDs(lineTaxIDs, ln.Taxes)
	}
	return totals, lineTaxIDs
}

func headerTaxRows(untaxedSubtotal fields.DecimalSix, headerTaxes []finance_taxes.Tax, lineTaxIDs map[uint]struct{}) (levied, withholding []InvoiceHeaderTaxAmountRow) {
	for _, t := range documentLevelHeaderTaxes(headerTaxes, lineTaxIDs) {
		amt := taxAmountForTax(untaxedSubtotal, t)
		label := t.Name
		if strings.TrimSpace(label) == "" {
			label = fmt.Sprintf("Tax #%d", t.ID)
		}
		row := InvoiceHeaderTaxAmountRow{Label: label}
		if effectiveTaxKind(t) == finance_taxes.TaxKindWithholding {
			row.Amount = decimalSixDisplayWithholding(amt)
			row.IsWithholding = true
			withholding = append(withholding, row)
		} else {
			row.Amount = decimalSixDisplay(amt)
			levied = append(levied, row)
		}
	}
	return levied, withholding
}

func finishInvoiceLinesSummary(totals invoiceLinesTotals, headerTaxes []finance_taxes.Tax, lineTaxIDs map[uint]struct{}, lineCount int) InvoiceLinesSummary {
	taxRows, whHeaderRows := headerTaxRows(totals.UntaxedSubtotal, headerTaxes, lineTaxIDs)

	var withholdingRows []InvoiceHeaderTaxAmountRow
	if !decimalIsZero(totals.LinesWithholding) {
		withholdingRows = append(withholdingRows, InvoiceHeaderTaxAmountRow{
			Label:         "Withholding (lines)",
			Amount:        decimalSixDisplayWithholding(totals.LinesWithholding),
			IsWithholding: true,
		})
	}
	withholdingRows = append(withholdingRows, whHeaderRows...)

	grand := invoiceReceivableGrandTotal(totals, headerTaxes, lineTaxIDs)

	linesSub := "—"
	if lineCount > 0 {
		linesSub = decimalSixDisplay(totals.linesGrossBeforeWithholding())
	}
	grandStr := decimalSixDisplay(grand)
	if lineCount == 0 && len(taxRows) == 0 && len(withholdingRows) == 0 && decimalIsZero(grand) {
		grandStr = "—"
	}
	return InvoiceLinesSummary{LinesSubtotal: linesSub, TaxRows: taxRows, WithholdingRows: withholdingRows, GrandTotal: grandStr}
}

func decimalIsZero(d fields.DecimalSix) bool {
	return d.R == nil || d.R.Sign() == 0
}

func invoiceLinesSummaryFromDraftLines(lines []DraftInvoiceLine, headerTaxes []finance_taxes.Tax) InvoiceLinesSummary {
	totals, lineTaxIDs := accumulateInvoiceLineTotals(lines)
	return finishInvoiceLinesSummary(totals, headerTaxes, lineTaxIDs, len(lines))
}

func invoiceLinesSummaryFromPostedLines(lines []PostedInvoiceLine, headerTaxes []finance_taxes.Tax) InvoiceLinesSummary {
	totals, lineTaxIDs := accumulatePostedInvoiceLineTotals(lines)
	return finishInvoiceLinesSummary(totals, headerTaxes, lineTaxIDs, len(lines))
}

func invoiceLinesSummaryFromCancelledLines(lines []CancelledInvoiceLine, headerTaxes []finance_taxes.Tax) InvoiceLinesSummary {
	totals, lineTaxIDs := accumulateCancelledInvoiceLineTotals(lines)
	return finishInvoiceLinesSummary(totals, headerTaxes, lineTaxIDs, len(lines))
}

// InvoiceLineDisplay is one read-only row for [FieldInvoiceLines].
type InvoiceLineDisplay struct {
	Product           string
	Quantity          string
	Rate              string
	LineTaxes         string
	UntaxedAmount     string
	LeviedTaxAmount   string
	WithholdingAmount string
	LineTotal         string
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
			Td(ColSpan("8"), Class("text-center opacity-50 py-4"), Text("No lines")),
		))
	} else {
		for _, r := range rows {
			tbody = append(tbody, Tr(
				Td(Class("whitespace-nowrap max-w-md"), Text(r.Product)),
				Td(Class("whitespace-nowrap text-end tabular-nums"), Text(r.Quantity)),
				Td(Class("whitespace-nowrap text-end tabular-nums"), Text(r.Rate)),
				Td(Class("min-w-[10rem] max-w-md text-sm"), Text(r.LineTaxes)),
				Td(Class("whitespace-nowrap text-end tabular-nums"), Text(r.UntaxedAmount)),
				Td(Class("whitespace-nowrap text-end tabular-nums"), Text(r.LeviedTaxAmount)),
				Td(Class("whitespace-nowrap text-end tabular-nums"), Text(r.WithholdingAmount)),
				Td(Class("whitespace-nowrap text-end tabular-nums font-medium"), Text(r.LineTotal)),
			))
		}
	}
	return Div(
		Class(wrap),
		Div(
			Class("overflow-x-auto rounded-box border border-base-300 bg-base-100"),
			Table(
				Class("table table-sm w-full"),
				THead(Tr(
					Th(Class("whitespace-nowrap min-w-[12rem]"), Text("Product")),
					Th(Class("whitespace-nowrap w-32 text-end"), Text("Quantity")),
					Th(Class("whitespace-nowrap w-32 text-end"), Text("Rate")),
					Th(Class("whitespace-nowrap min-w-[10rem]"), Text("Line taxes")),
					Th(Class("whitespace-nowrap min-w-[7rem] text-end"), Text("Untaxed amount")),
					Th(Class("whitespace-nowrap min-w-[7rem] text-end"), Text("Levied tax")),
					Th(Class("whitespace-nowrap min-w-[7rem] text-end"), Text("Withholding")),
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
		Div(
			Class("grid grid-cols-[1fr_auto] gap-x-4 items-center px-4 py-3"),
			Div(Class("text-sm font-bold min-w-0 truncate"), Text("Lines subtotal")),
			Div(Class("text-sm tabular-nums text-end font-semibold shrink-0 min-w-[7rem]"), Text(s.LinesSubtotal)),
		),
	}
	for _, row := range s.TaxRows {
		rowNodes = append(rowNodes, Div(
			Class("grid grid-cols-[1fr_auto] gap-x-4 items-center px-4 py-3"),
			Div(Class("text-sm font-bold min-w-0 truncate"), Text(row.Label)),
			Div(Class("text-sm tabular-nums text-end font-semibold shrink-0 min-w-[7rem]"), Text(row.Amount)),
		))
	}
	for _, row := range s.WithholdingRows {
		rowNodes = append(rowNodes, Div(
			Class("grid grid-cols-[1fr_auto] gap-x-4 items-center px-4 py-3"),
			Div(Class("text-sm font-bold min-w-0 truncate"), Text(row.Label)),
			Div(Class("text-sm tabular-nums text-end font-semibold shrink-0 min-w-[7rem]"), Text(row.Amount)),
		))
	}
	rowNodes = append(rowNodes, Div(
		Class("grid grid-cols-[1fr_auto] gap-x-4 items-center px-4 py-3 bg-base-200/60"),
		Div(Class("text-sm font-bold min-w-0 truncate"), Text("Total")),
		Div(Class("text-sm tabular-nums text-end font-bold shrink-0 min-w-[7rem]"), Text(s.GrandTotal)),
	))
	return Div(
		Class(outerClass),
		Div(append([]Node{Class("rounded-box border border-base-300 bg-base-100 overflow-hidden divide-y divide-base-300")}, rowNodes...)...),
	)
}
