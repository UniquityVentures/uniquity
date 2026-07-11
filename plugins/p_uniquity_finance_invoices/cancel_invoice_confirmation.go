package p_uniquity_finance_invoices

import (
	"context"
	"net/http"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

var _ components.FormInterface = CancelInvoiceConfirmation{}

// CancelInvoiceConfirmation is the modal body for cancelling a posted invoice (warning + reason + submit).
type CancelInvoiceConfirmation struct {
	components.Page
	Title   string
	Message string
	Classes string
	Attr    getters.Getter[Node]
}

func (e CancelInvoiceConfirmation) GetKey() string     { return e.Key }
func (e CancelInvoiceConfirmation) GetRoles() []string { return e.Roles }

func cancelInvoiceConfirmationGlobalError(ctx context.Context) Node {
	err, lookupErr := getters.Key[error]("$error._global")(ctx)
	if lookupErr != nil || err == nil {
		return nil
	}
	return Div(Class("alert alert-error my-2 text-sm"), Text(err.Error()))
}

func (e CancelInvoiceConfirmation) Build(ctx context.Context) Node {
	form := components.FormComponent[struct{}]{
		Classes: "gap-2 my-4",
		Attr:    e.Attr,
		ChildrenInput: []components.PageInterface{
			&components.InputText{Name: "Reason", Label: "Reason"},
		},
		ChildrenAction: []components.PageInterface{&components.ButtonSubmit{Label: "Cancel invoice", Classes: "btn-error my-2"}},
	}

	title := e.Title
	if title == "" {
		title = "Cancel this invoice?"
	}
	msg := e.Message
	if msg == "" {
		msg = "Cancelling creates a credit note that reverses the journal entry. The invoice is recorded as cancelled; you cannot restore the posted-only state."
	}

	return Div(Class("container mx-auto "+e.Classes),
		H2(Class("text-xl font-bold"), Text(title)),
		P(Class("text-sm text-gray-500 my-2"), Text("Creates a credit note reversing the journal entry.")),
		Div(Class("alert alert-warning text-sm my-3 leading-relaxed flex flex-col gap-2 items-stretch"),
			P(Class("font-semibold"), Text("This action cannot be reverted.")),
			P(Class("m-0"), Text(msg)),
		),
		cancelInvoiceConfirmationGlobalError(ctx),
		form.Build(ctx),
	)
}

func (e CancelInvoiceConfirmation) ParseForm(r *http.Request) (map[string]any, map[string]error, error) {
	inner := components.FormComponent[struct{}]{
		ChildrenInput: []components.PageInterface{
			&components.InputText{Name: "Reason", Label: "Reason"},
		},
		ChildrenAction: []components.PageInterface{&components.ButtonSubmit{Label: "Cancel invoice", Classes: "btn-error my-2"}},
	}
	return inner.ParseForm(r)
}
