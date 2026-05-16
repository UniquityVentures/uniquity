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
type InputInvoiceLinesDraft struct {
	components.Page
	Label    string
	Name     string
	Choices  getters.Getter[[]registry.Pair[uint, string]]
	Defaults getters.Getter[string]
	Classes  string
}

func (e InputInvoiceLinesDraft) GetKey() string { return e.Key }

func (e InputInvoiceLinesDraft) GetRoles() []string { return e.Roles }

type invoiceLineProductOpt struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func (e InputInvoiceLinesDraft) Build(ctx context.Context) Node {
	opts := []invoiceLineProductOpt{}
	if e.Choices != nil {
		pairs, err := e.Choices(ctx)
		if err != nil {
			slog.Error("InputInvoiceLinesDraft Choices failed", "error", err, "key", e.Key)
		} else {
			for _, p := range pairs {
				opts = append(opts, invoiceLineProductOpt{ID: p.Key, Name: p.Value})
			}
		}
	}
	productsJSON, err := json.Marshal(opts)
	if err != nil {
		productsJSON = []byte("[]")
	}
	defaults := `[{"product_id":0,"quantity":"1","rate":""}]`
	if e.Defaults != nil {
		if s, err := e.Defaults(ctx); err == nil && strings.TrimSpace(s) != "" {
			defaults = strings.TrimSpace(s)
		}
	}
	alpineData := fmt.Sprintf(`{ lines: %s, products: %s }`, defaults, string(productsJSON))
	initJS := fmt.Sprintf(`
$el.closest('form').addEventListener('submit', (ev) => {
	const d = Alpine.$data($el);
	if (!d || !Array.isArray(d.lines)) return;
	const h = $el.querySelector('input[type="hidden"][name=%s]');
	if (h) h.value = JSON.stringify(d.lines);
}, true);
`, strconv.Quote(e.Name))

	wrapClass := fmt.Sprintf("my-1 w-full %s", e.Classes)
	return Div(Class(wrapClass),
		Label(Class("label text-sm font-bold flex flex-col items-stretch gap-1 w-full"),
			Text(e.Label),
			Div(
				Class("w-full"),
				Attr("x-data", alpineData),
				Attr("x-init", initJS),
				Div(Class("overflow-x-auto rounded-box border border-base-300 bg-base-100"),
					Table(Class("table table-sm w-full"),
						THead(Tr(
							Th(Class("whitespace-nowrap min-w-[12rem]"), Text("Product")),
							Th(Class("whitespace-nowrap w-32"), Text("Quantity")),
							Th(Class("whitespace-nowrap w-32"), Text("Rate")),
							Th(Class("whitespace-nowrap w-24"), Text("Actions")),
						)),
						TBody(
							Template(
								Attr("x-for", "(line, i) in lines"),
								Attr(":key", "i"),
								Tr(
									Td(Class("align-middle"),
										Select(
											Class("select select-bordered w-full max-w-md"),
											Attr("x-model.number", "line.product_id"),
											Option(Value("0"), Text("—")),
											Template(
												Attr("x-for", "p in products"),
												Attr(":key", "p.id"),
												Option(Attr(":value", "p.id"), Attr("x-text", "p.name")),
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
									Td(Class("align-middle w-24"),
										Button(
											Type("button"),
											Class("btn btn-ghost btn-sm"),
											Attr("@click", "lines.splice(i, 1); if (lines.length === 0) lines.push({ product_id: 0, quantity: '1', rate: '' })"),
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
					Attr("@click", "lines.push({ product_id: 0, quantity: '1', rate: '' })"),
					Text("Add line"),
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
