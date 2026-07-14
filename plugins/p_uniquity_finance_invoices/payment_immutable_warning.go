package p_uniquity_finance_invoices

import (
	"context"

	"github.com/lariv-in/lago/components"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type paymentImmutableWarning struct {
	components.Page
}

func (e paymentImmutableWarning) GetKey() string     { return e.Key }
func (e paymentImmutableWarning) GetRoles() []string { return e.Roles }

func (paymentImmutableWarning) Build(context.Context) Node {
	return Div(
		Class("alert alert-warning text-sm my-3 leading-relaxed"),
		P(Class("font-semibold m-0"), Text("Payments cannot be edited or deleted.")),
		P(Class("m-0 mt-1"), Text("Once saved, this record is permanent. Verify the invoice, amount, and account before submitting.")),
	)
}
