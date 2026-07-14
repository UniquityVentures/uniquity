package p_uniquity_finance_invoices

import (
	"context"
	"net/http"

	"github.com/lariv-in/lago/components"
	"github.com/lariv-in/lago/getters"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

var _ components.FormInterface = NewDraftFromCancelledConfirmation{}

type newDraftFromCancelledSubmitBtn struct {
	components.Page
}

func (e newDraftFromCancelledSubmitBtn) GetKey() string     { return e.Key }
func (e newDraftFromCancelledSubmitBtn) GetRoles() []string { return e.Roles }

func (newDraftFromCancelledSubmitBtn) Build(context.Context) Node {
	return Button(Type("submit"), Class("btn btn-primary my-2"), Text("Create draft"))
}

// NewDraftFromCancelledConfirmation is the modal body for creating a draft from a cancelled invoice.
type NewDraftFromCancelledConfirmation struct {
	components.Page
	Title   string
	Message string
	Classes string
	Attr    getters.Getter[Node]
}

func (e NewDraftFromCancelledConfirmation) GetKey() string     { return e.Key }
func (e NewDraftFromCancelledConfirmation) GetRoles() []string { return e.Roles }

func newDraftFromCancelledGlobalError(ctx context.Context) Node {
	err, lookupErr := getters.Key[error]("$error._global")(ctx)
	if lookupErr != nil || err == nil {
		return nil
	}
	return Div(Class("alert alert-error my-2 text-sm"), Text(err.Error()))
}

func (e NewDraftFromCancelledConfirmation) Build(ctx context.Context) Node {
	form := components.FormComponent[struct{}]{
		Classes:        "gap-2 my-4",
		Attr:           e.Attr,
		ChildrenAction: []components.PageInterface{newDraftFromCancelledSubmitBtn{}},
	}

	title := e.Title
	if title == "" {
		title = "Create new draft?"
	}
	msg := e.Message
	if msg == "" {
		msg = "A new editable draft invoice will be copied from this cancelled invoice. The cancelled record is unchanged."
	}

	return Div(
		Class("container mx-auto "+e.Classes),
		H2(Class("text-xl font-bold"), Text(title)),
		P(Class("text-sm text-gray-500 my-2"), Text(msg)),
		newDraftFromCancelledGlobalError(ctx),
		form.Build(ctx),
	)
}

func (e NewDraftFromCancelledConfirmation) ParseForm(r *http.Request) (map[string]any, map[string]error, error) {
	inner := components.FormComponent[struct{}]{
		ChildrenAction: []components.PageInterface{newDraftFromCancelledSubmitBtn{}},
	}
	return inner.ParseForm(r)
}
