package p_uniquity_finance_invoices

import (
	"context"
	"net/http"

	"github.com/lariv-in/lago/components"
	"github.com/lariv-in/lago/getters"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

var _ components.FormInterface = PostInvoiceConfirmation{}

type postInvoiceConfirmSubmitBtn struct {
	components.Page
}

func (e postInvoiceConfirmSubmitBtn) GetKey() string     { return e.Key }
func (e postInvoiceConfirmSubmitBtn) GetRoles() []string { return e.Roles }

func (postInvoiceConfirmSubmitBtn) Build(context.Context) Node {
	return Button(Type("submit"), Class("btn btn-primary my-2"), Text("Post invoice"))
}

// PostInvoiceConfirmation is the modal body for posting a draft (irreversibility warning + submit).
type PostInvoiceConfirmation struct {
	components.Page
	Title   string
	Message string
	Classes string
	Attr    getters.Getter[Node]
}

func (e PostInvoiceConfirmation) GetKey() string     { return e.Key }
func (e PostInvoiceConfirmation) GetRoles() []string { return e.Roles }

func postInvoiceConfirmationGlobalError(ctx context.Context) Node {
	err, lookupErr := getters.Key[error]("$error._global")(ctx)
	if lookupErr != nil || err == nil {
		return nil
	}
	return Div(Class("alert alert-error my-2 text-sm"), Text(err.Error()))
}

func (e PostInvoiceConfirmation) Build(ctx context.Context) Node {
	form := components.FormComponent[struct{}]{
		Classes:        "gap-2 my-4",
		Attr:           e.Attr,
		ChildrenAction: []components.PageInterface{postInvoiceConfirmSubmitBtn{}},
	}

	title := e.Title
	if title == "" {
		title = "Post this invoice?"
	}
	msg := e.Message
	if msg == "" {
		msg = "Posting creates a permanent journal entry and locks this invoice. You cannot undo posting or return to draft. Reversal requires cancelling the posted invoice later (credit note), which creates additional accounting entries."
	}

	return Div(
		Class("container mx-auto "+e.Classes),
		H2(Class("text-xl font-bold"), Text(title)),
		Div(
			Class("alert alert-warning text-sm my-3 leading-relaxed flex flex-col gap-2 items-stretch"),
			P(Class("font-semibold"), Text("This action cannot be reverted.")),
			P(Class("m-0"), Text(msg)),
		),
		postInvoiceConfirmationGlobalError(ctx),
		form.Build(ctx),
	)
}

func (e PostInvoiceConfirmation) ParseForm(r *http.Request) (map[string]any, map[string]error, error) {
	inner := components.FormComponent[struct{}]{
		ChildrenAction: []components.PageInterface{postInvoiceConfirmSubmitBtn{}},
	}
	return inner.ParseForm(r)
}
