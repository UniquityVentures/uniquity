package p_uniquity_finance_invoices

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/registry"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// InputInvoiceLinesDraft edits invoice lines as JSON in a hidden field (Alpine + submit capture).
// Each row uses the same HTMX + fk-select pattern as [components.InputForeignKey]; [ProductPickURL]
// must point at a route that reads target_input (see finance_products.ProductFkSelectRoute).
// Per-line taxes use fk-multi-select with [TaxPickURL] (see finance_taxes.TaxMultiSelectRoute).
type InputInvoiceLinesDraft struct {
	components.Page
	Label   string
	Name    string
	Choices getters.Getter[[]registry.Pair[uint, string]] // optional fallback when Preview is nil
	// Preview returns JSON: {"products":[...], "tax_pct_by_id":{...}, "all_taxes":[{"id","name"},...]}.
	Preview        getters.Getter[string]
	ProductPickURL getters.Getter[string]
	TaxPickURL     getters.Getter[string]
	Defaults       getters.Getter[string]
	Classes        string
}

func (e InputInvoiceLinesDraft) GetKey() string { return e.Key }

func (e InputInvoiceLinesDraft) GetRoles() []string { return e.Roles }

type invoiceLineProductOpt struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	SalesPrice string `json:"sales_price,omitempty"`
	TaxIDs     []uint `json:"tax_ids,omitempty"`
}

const invoiceLinesDraftAlpineMethods = `allocFkSlot() {
	if (typeof crypto !== 'undefined' && crypto.randomUUID) return 'InvoiceLineProduct_' + crypto.randomUUID();
	return 'InvoiceLineProduct_' + Math.random().toString(36).slice(2) + '_' + Date.now().toString(36);
},
productPickHref(slot) {
	const b = this.product_pick_base || '';
	if (!b || !slot) return b || '#';
	const sep = b.indexOf('?') >= 0 ? '&' : '?';
	return b + sep + 'target_input=' + encodeURIComponent(String(slot));
},
lineTaxPickHref(fkSlot) {
	const b = this.tax_pick_base || '';
	if (!b || !fkSlot) return b || '#';
	const sep = b.indexOf('?') >= 0 ? '&' : '?';
	const name = 'InvoiceLineTaxes_' + String(fkSlot);
	return b + sep + 'target_input=' + encodeURIComponent(name);
},
formatDec(n) {
	if (typeof n !== 'number' || !isFinite(n)) return '—';
	let s = n.toFixed(6);
	s = s.replace(/0+$/, '');
	s = s.replace(/\.$/, '');
	return s || '0';
},
lineUntaxedNumber(line) {
	const q = parseFloat(String(line.quantity).replace(/,/g, '.')) || 0;
	const rate = parseFloat(String(line.rate ?? '').trim().replace(/,/g, '.')) || 0;
	return q * rate;
},
taxKindForId(id) {
	const k = this.tax_kind_by_id && this.tax_kind_by_id[id];
	return k === 'withholding' ? 'withholding' : 'levied';
},
lineTaxAmountForKind(line, kind) {
	const base = this.lineUntaxedNumber(line);
	if (!Array.isArray(line.line_taxes) || line.line_taxes.length === 0) return 0;
	let sum = 0;
	for (const t of line.line_taxes) {
		const id = String(t.Key);
		if (this.taxKindForId(id) !== kind) continue;
		const pctStr = this.tax_pct_by_id[id];
		const pct = pctStr != null && pctStr !== '' ? parseFloat(String(pctStr)) : NaN;
		if (!isNaN(pct)) {
			sum += base * (pct / 100);
		}
	}
	return sum;
},
lineLeviedTaxNumber(line) {
	return this.lineTaxAmountForKind(line, 'levied');
},
lineWithholdingTaxNumber(line) {
	return this.lineTaxAmountForKind(line, 'withholding');
},
lineUntaxedDisplay(line) {
	const u = this.lineUntaxedNumber(line);
	if (!line.product_id && u === 0) return '—';
	return this.formatDec(u);
},
lineLeviedTaxDisplay(line) {
	const u = this.lineUntaxedNumber(line);
	const tax = this.lineLeviedTaxNumber(line);
	if (!line.product_id && u === 0 && tax === 0) return '—';
	return this.formatDec(tax);
},
lineWithholdingDisplay(line) {
	const wh = this.lineWithholdingTaxNumber(line);
	if (wh === 0) return '—';
	return '(' + this.formatDec(wh) + ')';
},
lineTotal(line) {
	const u = this.lineUntaxedNumber(line);
	const lev = this.lineLeviedTaxNumber(line);
	const wh = this.lineWithholdingTaxNumber(line);
	const tot = u + lev - wh;
	if (!line.product_id && tot === 0) return '—';
	return this.formatDec(tot);
},
lineTotalNumber(line) {
	return this.lineUntaxedNumber(line) + this.lineLeviedTaxNumber(line) - this.lineWithholdingTaxNumber(line);
},
linesUntaxedSubtotalNumber() {
	if (!Array.isArray(this.lines)) return 0;
	let sum = 0;
	for (const line of this.lines) {
		sum += this.lineUntaxedNumber(line);
	}
	return sum;
},
linesLeviedSubtotalNumber() {
	if (!Array.isArray(this.lines)) return 0;
	let sum = 0;
	for (const line of this.lines) {
		sum += this.lineLeviedTaxNumber(line);
	}
	return sum;
},
linesWithholdingSubtotalNumber() {
	if (!Array.isArray(this.lines)) return 0;
	let sum = 0;
	for (const line of this.lines) {
		sum += this.lineWithholdingTaxNumber(line);
	}
	return sum;
},
linesSubtotalNumber() {
	return this.linesUntaxedSubtotalNumber() + this.linesLeviedSubtotalNumber();
},
linesSubtotalDisplay() {
	if (!Array.isArray(this.lines) || this.lines.length === 0) return '—';
	return this.formatDec(this.linesSubtotalNumber());
},
invoiceTaxLabel(item) {
	if (item.Value != null && String(item.Value).trim() !== '') {
		return String(item.Value);
	}
	return 'Tax #' + String(item.Key);
},
invoiceHeaderTaxAmountForKind(kind) {
	const base = this.linesUntaxedSubtotalNumber();
	const st = typeof Alpine !== 'undefined' && Alpine.store && Alpine.store('m2mSelections');
	const sel = st && st.Taxes;
	if (!sel || !Array.isArray(sel)) {
		return 0;
	}
	let sum = 0;
	for (const item of sel) {
		const id = String(item.Key);
		if (this.taxKindForId(id) !== kind) continue;
		const pctStr = this.tax_pct_by_id[id];
		const pct = pctStr != null && pctStr !== '' ? parseFloat(String(pctStr)) : NaN;
		if (!isNaN(pct)) {
			sum += base * (pct / 100);
		}
	}
	return sum;
},
invoiceTaxAmountDisplay(item) {
	const base = this.linesUntaxedSubtotalNumber();
	const id = String(item.Key);
	const kind = this.taxKindForId(id);
	const pctStr = this.tax_pct_by_id[id];
	const pct = pctStr != null && pctStr !== '' ? parseFloat(String(pctStr)) : NaN;
	if (isNaN(pct)) {
		return '—';
	}
	const amt = base * (pct / 100);
	if (kind === 'withholding') {
		return '(' + this.formatDec(amt) + ')';
	}
	return this.formatDec(amt);
},
invoiceGrandTotalDisplay() {
	const sub = this.linesSubtotalNumber();
	const headerLevied = this.invoiceHeaderTaxAmountForKind('levied');
	const headerWithholding = this.invoiceHeaderTaxAmountForKind('withholding');
	const lineWh = this.linesWithholdingSubtotalNumber();
	const total = sub + headerLevied - lineWh - headerWithholding;
	const st = typeof Alpine !== 'undefined' && Alpine.store && Alpine.store('m2mSelections');
	const sel = st && st.Taxes;
	const hasTaxes = sel && Array.isArray(sel) && sel.length > 0;
	const hasLines = Array.isArray(this.lines) && this.lines.length > 0;
	if (!hasLines && !hasTaxes && total === 0) {
		return '—';
	}
	return this.formatDec(total);
}`

func (e InputInvoiceLinesDraft) Build(ctx context.Context) Node {
	var productsJSON []byte
	var taxPctJSON = []byte("{}")
	var taxKindJSON = []byte("{}")
	var allTaxesJSON = []byte("[]")
	if e.Preview != nil {
		if s, err := e.Preview(ctx); err != nil {
			slog.Error("InputInvoiceLinesDraft Preview failed", "error", err, "key", e.Key)
		} else if strings.TrimSpace(s) != "" {
			var prev struct {
				Products    []invoiceLineProductOpt `json:"products"`
				TaxPctByID  map[string]string       `json:"tax_pct_by_id"`
				TaxKindByID map[string]string       `json:"tax_kind_by_id"`
				AllTaxes    []invoiceLineTaxMeta    `json:"all_taxes"`
			}
			if err := json.Unmarshal([]byte(s), &prev); err != nil {
				slog.Error("InputInvoiceLinesDraft Preview unmarshal", "error", err, "key", e.Key)
			} else {
				b, err := json.Marshal(prev.Products)
				if err != nil {
					slog.Error("InputInvoiceLinesDraft products marshal", "error", err, "key", e.Key)
				} else {
					productsJSON = b
				}
				if prev.TaxPctByID != nil {
					if b, err := json.Marshal(prev.TaxPctByID); err == nil {
						taxPctJSON = b
					}
				}
				if prev.TaxKindByID != nil {
					if b, err := json.Marshal(prev.TaxKindByID); err == nil {
						taxKindJSON = b
					}
				}
				if prev.AllTaxes != nil {
					if b, err := json.Marshal(prev.AllTaxes); err == nil {
						allTaxesJSON = b
					}
				}
			}
		}
	}
	if len(productsJSON) == 0 && e.Choices != nil {
		opts := []invoiceLineProductOpt{}
		pairs, err := e.Choices(ctx)
		if err != nil {
			slog.Error("InputInvoiceLinesDraft Choices failed", "error", err, "key", e.Key)
		} else {
			for _, p := range pairs {
				opts = append(opts, invoiceLineProductOpt{ID: p.Key, Name: p.Value})
			}
		}
		productsJSON, err = json.Marshal(opts)
		if err != nil {
			productsJSON = []byte("[]")
		}
	}
	if len(productsJSON) == 0 {
		productsJSON = []byte("[]")
	}
	defaults := `[{"product_id":0,"quantity":"1","rate":"","product_label":"","fk_slot":"line-slot-0","tax_ids":[]}]`
	if e.Defaults != nil {
		if s, err := e.Defaults(ctx); err == nil && strings.TrimSpace(s) != "" {
			defaults = strings.TrimSpace(s)
		}
	}
	pickBase := ""
	if e.ProductPickURL != nil {
		if u, err := e.ProductPickURL(ctx); err != nil {
			slog.Error("InputInvoiceLinesDraft ProductPickURL failed", "error", err, "key", e.Key)
		} else {
			pickBase = u
		}
	}
	pickBaseJSON, err := json.Marshal(pickBase)
	if err != nil {
		pickBaseJSON = []byte(`""`)
	}
	taxPickBase := ""
	if e.TaxPickURL != nil {
		if u, err := e.TaxPickURL(ctx); err != nil {
			slog.Error("InputInvoiceLinesDraft TaxPickURL failed", "error", err, "key", e.Key)
		} else {
			taxPickBase = u
		}
	}
	taxPickBaseJSON, err := json.Marshal(taxPickBase)
	if err != nil {
		taxPickBaseJSON = []byte(`""`)
	}
	alpineData := fmt.Sprintf("{ lines: %s, products: %s, tax_pct_by_id: %s, tax_kind_by_id: %s, all_taxes: %s, product_pick_base: %s, tax_pick_base: %s, %s }",
		defaults, string(productsJSON), string(taxPctJSON), string(taxKindJSON), string(allTaxesJSON), string(pickBaseJSON), string(taxPickBaseJSON), strings.TrimSpace(invoiceLinesDraftAlpineMethods))
	initJS := fmt.Sprintf(`
if (typeof Alpine !== 'undefined' && Alpine.store && !Alpine.store('m2mSelections')) {
	Alpine.store('m2mSelections', {});
}
(function () {
	const d = Alpine.$data($el);
	if (!d || !Array.isArray(d.lines) || typeof d.allocFkSlot !== 'function') return;
	for (const line of d.lines) {
		if (line.product_label == null) line.product_label = '';
		if (!line.fk_slot) line.fk_slot = d.allocFkSlot();
		if (!Array.isArray(line.line_taxes)) line.line_taxes = [];
		const ids = line.tax_ids;
		if (Array.isArray(ids) && ids.length > 0 && line.line_taxes.length === 0 && Array.isArray(d.all_taxes)) {
			for (const tid of ids) {
				const t = d.all_taxes.find(x => x.id === tid);
				if (t) line.line_taxes.push({ Key: String(t.id), Value: t.name });
			}
		}
		delete line.tax_ids;
	}
})();
$nextTick(() => { if (window.htmx) window.htmx.process($el); });
$el.closest('form').addEventListener('submit', (ev) => {
	const d = Alpine.$data($el);
	if (!d || !Array.isArray(d.lines)) return;
	const h = $el.querySelector('input[type="hidden"][name=%s]');
	if (!h) return;
	const strip = (l) => ({
		product_id: l.product_id,
		quantity: l.quantity,
		rate: l.rate,
		product_label: l.product_label,
		fk_slot: l.fk_slot,
		tax_ids: (l.line_taxes || []).map(t => parseInt(String(t.Key), 10)).filter(id => !isNaN(id) && id > 0),
	});
	h.value = JSON.stringify(d.lines.map(strip));
}, true);
`, strconv.Quote(e.Name))

	wrapClass := fmt.Sprintf("my-1 w-full %s", e.Classes)
	// Match [components.InputText] field chrome (label + control), but use <div> not <label> so the table is not a labeled control.
	return Div(Class(wrapClass),
		Div(Class("label text-sm font-bold flex flex-col items-start gap-1 w-full"),
			Text(e.Label),
			Div(
				Class("w-full"),
				Attr("data-invoice-lines-root", ""),
				Attr("x-data", alpineData),
				Attr("x-init", initJS),
				// Alpine adds x-bind:hx-get after first paint; HTMX only scans once unless we process again.
				Attr("x-effect", "lines.length; $nextTick(() => { if (window.htmx) window.htmx.process($el); })"),
				Attr("@fk-select.window", `if (!$event.detail) return;
	const n = $event.detail.name;
	const v = $event.detail.value;
	const disp = $event.detail.display || '';
	for (const line of lines) {
		if (!line.fk_slot || line.fk_slot !== n) continue;
		const pid = parseInt(String(v), 10) || 0;
		line.product_id = pid;
		line.product_label = disp;
		if (!pid) { line.rate = ''; line.line_taxes = []; continue; }
		const prod = products.find(p => p.id === pid);
		const sp = prod && prod.sales_price != null && String(prod.sales_price).trim() !== '' ? String(prod.sales_price).trim() : '';
		line.rate = sp;
		break;
	}`),
				Attr("@fk-multi-select.window", `if (!$event.detail) return;
	const n = String($event.detail.name || '');
	const v = $event.detail.value;
	const disp = $event.detail.display || '';
	for (const line of lines) {
		const expected = 'InvoiceLineTaxes_' + String(line.fk_slot || '');
		if (expected !== n) continue;
		const value = String(v);
		const items = line.line_taxes || (line.line_taxes = []);
		const idx = items.findIndex(x => x.Key === value);
		if (idx >= 0) items.splice(idx, 1);
		else items.push({ Key: value, Value: String(disp || value) });
		break;
	}`),
				Div(Class("overflow-x-auto rounded-box border border-base-300 bg-base-100"),
					Table(Class("table table-sm w-full"),
						THead(Tr(
							Th(Class("whitespace-nowrap min-w-[12rem]"), Text("Product")),
							Th(Class("whitespace-nowrap w-32"), Text("Quantity")),
							Th(Class("whitespace-nowrap w-32"), Text("Rate")),
							Th(Class("whitespace-nowrap min-w-[10rem]"), Text("Line taxes")),
							Th(Class("whitespace-nowrap min-w-[7rem] text-end"), Text("Untaxed amount")),
							Th(Class("whitespace-nowrap min-w-[7rem] text-end"), Text("Levied tax")),
							Th(Class("whitespace-nowrap min-w-[7rem] text-end"), Text("Withholding")),
							Th(Class("whitespace-nowrap min-w-[7rem] text-end"), Text("Line total")),
							Th(Class("whitespace-nowrap w-24"), Text("Actions")),
						)),
						TBody(
							Template(
								Attr("x-for", "(line, i) in lines"),
								Attr(":key", "line.fk_slot"),
								Tr(
									Td(Class("align-middle min-w-[12rem] max-w-md"),
										Div(Class("my-1 relative w-full"),
											Div(Class("flex w-full items-stretch gap-1"),
												Div(Class("input input-bordered flex-1 flex items-center cursor-pointer min-w-0"),
													Attr(":class", "line.product_label ? '' : 'opacity-50'"),
													Attr("x-bind:hx-get", "productPickHref(line.fk_slot)"),
													Attr("hx-target", components.HTMXTargetBodyModal),
													Attr("hx-swap", components.HTMXSwapBodyModal),
													Attr("hx-push-url", "false"),
													Span(Class("text-sm truncate"), Attr("x-text", "line.product_label || 'Select…'")),
												),
												Button(
													Type("button"),
													Class("btn btn-ghost btn-square shrink-0"),
													Attr("@click.stop", "line.product_id = 0; line.product_label = ''; line.rate = ''; line.line_taxes = []"),
													Attr("x-show", "line.product_id"),
													Attr("aria-label", "Clear product selection"),
													components.Render(components.Icon{Name: "x-mark"}, ctx),
												),
											),
										),
									),
									Td(Class("align-middle"),
										Input(
											Type("text"),
											Class("input input-bordered w-full"),
											Attr("x-model", "line.quantity"),
											Attr("inputmode", "decimal"),
										),
									),
									Td(Class("align-middle"),
										Input(
											Type("text"),
											Class("input input-bordered w-full"),
											Attr("x-model", "line.rate"),
											Attr("inputmode", "decimal"),
										),
									),
									Td(Class("align-middle min-w-[10rem] max-w-xs"),
										Div(Class("my-1"),
											Div(Class("input input-bordered min-h-10 w-full flex flex-wrap items-center gap-1 cursor-pointer py-1 px-2"),
												Attr(":class", "(line.line_taxes && line.line_taxes.length) ? '' : 'opacity-50'"),
												Attr("x-bind:hx-get", "lineTaxPickHref(line.fk_slot)"),
												Attr("hx-target", components.HTMXTargetBodyModal),
												Attr("hx-swap", components.HTMXSwapBodyModal),
												Attr("hx-push-url", "false"),
												Span(
													Class("text-sm"),
													Attr("x-show", "!line.line_taxes || line.line_taxes.length === 0"),
													Text("Select taxes…"),
												),
												Template(
													Attr("x-for", "ltItem in (line.line_taxes || [])"),
													Attr(":key", "ltItem.Key"),
													Div(
														Class("flex items-center gap-1 rounded-lg bg-base-200 pl-2 pr-1 py-0.5 max-w-full"),
														Attr("@click", "$event.stopPropagation()"),
														Span(Class("text-xs truncate max-w-[8rem]"), Attr("x-text", "ltItem.Value")),
														Button(
															Type("button"),
															Class("btn btn-ghost btn-square btn-xs shrink-0"),
															Attr("@click.stop", "line.line_taxes = (line.line_taxes || []).filter(it => it.Key !== ltItem.Key)"),
															Attr("aria-label", "Remove tax"),
															components.Render(components.Icon{Name: "x-mark"}, ctx),
														),
													),
												),
											),
										),
									),
									Td(Class("align-middle text-end tabular-nums whitespace-nowrap"),
										Span(
											Class("text-sm"),
											Attr("x-text", "lineUntaxedDisplay(line)"),
										),
									),
									Td(Class("align-middle text-end tabular-nums whitespace-nowrap"),
										Span(
											Class("text-sm"),
											Attr("x-text", "lineLeviedTaxDisplay(line)"),
										),
									),
									Td(Class("align-middle text-end tabular-nums whitespace-nowrap"),
										Span(
											Class("text-sm"),
											Attr("x-text", "lineWithholdingDisplay(line)"),
										),
									),
									Td(Class("align-middle text-end tabular-nums whitespace-nowrap"),
										Span(
											Class("text-sm"),
											Attr("x-text", "lineTotal(line)"),
										),
									),
									Td(Class("align-middle w-24"),
										Button(
											Type("button"),
											Class("btn btn-ghost btn-sm"),
											Attr("@click", "lines.splice(i, 1); if (lines.length === 0) lines.push({ product_id: 0, quantity: '1', rate: '', product_label: '', fk_slot: allocFkSlot(), line_taxes: [] }); $nextTick(() => { const r = $el.closest('[data-invoice-lines-root]'); if (r && window.htmx) window.htmx.process(r) })"),
											Text("Remove"),
										),
									),
								),
							),
						),
					),
				),
				Button(
					Type("button"),
					Class("btn btn-outline btn-sm mt-2 w-full sm:w-auto"),
					Attr("@click", "lines.push({ product_id: 0, quantity: '1', rate: '', product_label: '', fk_slot: allocFkSlot(), line_taxes: [] }); $nextTick(() => { const r = $el.closest('[data-invoice-lines-root]'); if (r && window.htmx) window.htmx.process(r) })"),
					Text("Add line"),
				),
				Div(Class("mt-3 w-full rounded-box border border-base-300 bg-base-100 overflow-hidden divide-y divide-base-300"),
					Div(Class("grid grid-cols-[1fr_auto] gap-x-4 items-center px-4 py-3"),
						Div(Class("text-sm font-bold min-w-0 truncate"), Text("Lines subtotal")),
						Div(
							Class("text-sm tabular-nums text-end font-semibold shrink-0 min-w-[7rem]"),
							Attr("x-text", "linesSubtotalDisplay()"),
						),
					),
					Template(
						Attr("x-show", "linesWithholdingSubtotalNumber() > 0"),
						Div(Class("grid grid-cols-[1fr_auto] gap-x-4 items-center px-4 py-3"),
							Div(Class("text-sm font-bold min-w-0 truncate"), Text("Withholding (lines)")),
							Div(
								Class("text-sm tabular-nums text-end font-semibold shrink-0 min-w-[7rem]"),
								Attr("x-text", "'(' + formatDec(linesWithholdingSubtotalNumber()) + ')'"),
							),
						),
					),
					Template(
						Attr("x-for", "invTaxItem in (($store.m2mSelections && $store.m2mSelections.Taxes) ? $store.m2mSelections.Taxes : [])"),
						Attr(":key", "invTaxItem.Key"),
						Div(Class("grid grid-cols-[1fr_auto] gap-x-4 items-center px-4 py-3"),
							Div(Class("text-sm font-bold min-w-0 truncate"), Attr("x-text", "invoiceTaxLabel(invTaxItem)")),
							Div(
								Class("text-sm tabular-nums text-end font-semibold shrink-0 min-w-[7rem]"),
								Attr("x-text", "invoiceTaxAmountDisplay(invTaxItem)"),
							),
						),
					),
					Div(Class("grid grid-cols-[1fr_auto] gap-x-4 items-center px-4 py-3 bg-base-200/60"),
						Div(Class("text-sm font-bold min-w-0 truncate"), Text("Total")),
						Div(
							Class("text-sm tabular-nums text-end font-bold shrink-0 min-w-[7rem]"),
							Attr("x-text", "(($store.m2mSelections && $store.m2mSelections.Taxes), invoiceGrandTotalDisplay())"),
						),
					),
				),
				Input(Type("hidden"), Name(e.Name)),
			),
		),
	)
}

func (e InputInvoiceLinesDraft) Parse(v any, _ context.Context) (any, error) {
	vals, _ := v.([]string)
	if len(vals) == 0 || strings.TrimSpace(vals[0]) == "" {
		return "[]", nil
	}
	return strings.TrimSpace(vals[0]), nil
}

func (e InputInvoiceLinesDraft) GetName() string { return e.Name }
