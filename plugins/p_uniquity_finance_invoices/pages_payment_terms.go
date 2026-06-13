package p_uniquity_finance_invoices

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
	. "maragu.dev/gomponents"
)

var paymentTermKindChoiceList = []registry.Pair[string, string]{
	{Key: PaymentTermTypeDueDate, Value: "Due on calendar date"},
	{Key: PaymentTermTypeRelative, Value: "Relative"},
}

func paymentTermKindChoices() getters.Getter[[]registry.Pair[string, string]] {
	return getters.Static(paymentTermKindChoiceList)
}

func paymentTermCreateTypeGetter() getters.Getter[registry.Pair[string, string]] {
	return registry.PairFromGetter(func(ctx context.Context) (string, error) {
		s, err := getters.Key[string]("$in.Type")(ctx)
		if err != nil {
			return "", err
		}
		if s == "" {
			if m, ok := paymentTermGetQuery(ctx); ok {
				if raw, ok := m["Type"].(string); ok && raw != "" {
					s = raw
				}
			}
		}
		if s == "" {
			s = PaymentTermTypeDueDate
		}
		return s, nil
	}, paymentTermKindChoiceList)
}

func paymentTermGetQuery(ctx context.Context) (map[string]any, bool) {
	v := ctx.Value(getters.ContextKeyGet)
	m, ok := v.(map[string]any)
	return m, ok && m != nil
}

func paymentTermCreateLinkWithQuery(typeParam string) getters.Getter[string] {
	base := lamu.RoutePath("finance_invoices.PaymentTermCreateRoute", nil)
	return func(ctx context.Context) (string, error) {
		b, err := base(ctx)
		if err != nil {
			return "", err
		}
		return b + "?Type=" + url.QueryEscape(typeParam), nil
	}
}

func paymentTermCreateDueDatetimeGetter() getters.Getter[time.Time] {
	return func(ctx context.Context) (time.Time, error) {
		t, err := getters.Key[time.Time]("$in.DueDatetime")(ctx)
		if err != nil {
			return time.Time{}, nil
		}
		return t, nil
	}
}

func paymentTermCreateDurationGetter() getters.Getter[*time.Duration] {
	return func(ctx context.Context) (*time.Duration, error) {
		d, err := getters.Key[*time.Duration]("$in.Duration")(ctx)
		if err != nil {
			return nil, nil
		}
		return d, nil
	}
}

func paymentTermRowSummaryGetter() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		id, err := getters.Key[uint]("$row.ID")(ctx)
		if err != nil {
			return "", err
		}
		typ, err := getters.Key[string]("$row.Type")(ctx)
		if err != nil {
			return "", err
		}
		bid, err := getters.Key[uint]("$row.BackingID")(ctx)
		if err != nil {
			return "", err
		}
		inst, err := ResolvePaymentTermInstance(ctx, typ, bid)
		if err != nil {
			return fmt.Sprintf("#%d — %s", id, typ), nil
		}
		return fmt.Sprintf("#%d — %s", id, inst.Summary()), nil
	}
}

func paymentTermDetailSummaryGetter() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		typ, err := getters.Key[string]("$in.Type")(ctx)
		if err != nil {
			return "", err
		}
		bid, err := getters.Key[uint]("$in.BackingID")(ctx)
		if err != nil {
			return "", err
		}
		inst, err := ResolvePaymentTermInstance(ctx, typ, bid)
		if err != nil {
			return fmt.Sprintf("%q #%d", typ, bid), nil
		}
		return inst.Summary(), nil
	}
}

func paymentTermCreateFormInputs() []components.PageInterface {
	return []components.PageInterface{
		&components.ClientData{
			Page: components.Page{Key: "finance_invoices.PaymentTermCreateForm.kindBranch"},
			Data: fmt.Sprintf("{ paymentTermType: %q }", PaymentTermTypeDueDate),
			Init: `(() => { const v = $el.querySelector('[name=Type]')?.value; if (v) paymentTermType = v })()`,
			Children: []components.PageInterface{
				&components.ContainerError{
					Error: getters.Key[error]("$error.Type"),
					Children: []components.PageInterface{
						&components.InputSelect[string]{
							Label:    "Kind",
							Name:     "Type",
							Required: true,
							Choices:  paymentTermKindChoices(),
							Getter:   paymentTermCreateTypeGetter(),
							Attr:     getters.Static(Attr("x-model", "paymentTermType")),
						},
					},
				},
				&components.ClientIf{
					Condition: fmt.Sprintf("paymentTermType === %q", PaymentTermTypeDueDate),
					Children: []components.PageInterface{
						&components.ContainerError{
							Error: getters.Key[error]("$error.DueDatetime"),
							Children: []components.PageInterface{
								&components.InputDatetime{
									Label:    "Due date & time",
									Name:     "DueDatetime",
									Required: true,
									Getter:   paymentTermCreateDueDatetimeGetter(),
								},
							},
						},
					},
				},
				&components.ClientIf{
					Condition: fmt.Sprintf("paymentTermType === %q", PaymentTermTypeRelative),
					Children: []components.PageInterface{
						&components.ContainerError{
							Error: getters.Key[error]("$error.Duration"),
							Children: []components.PageInterface{
								&components.InputDuration{
									Label:    "Offset duration (Go syntax, e.g. 720h, 30m)",
									Name:     "Duration",
									Required: true,
									Getter:   paymentTermCreateDurationGetter(),
								},
							},
						},
					},
				},
			},
		},
	}
}

func pageEntriesPaymentTermPages() []registry.Pair[string, components.PageInterface] {
	ptCreateName := getters.Static("finance_invoices.PaymentTermCreateForm")

	return []registry.Pair[string, components.PageInterface]{
		{Key: "finance_invoices.PaymentTermTable", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_accounts.MainMenu"}},
			Children: []components.PageInterface{
				&components.DataTable[PaymentTerm]{
					UID:     "finance-payment-term-table",
					Classes: "w-full",
					Data:    getters.Key[components.ObjectList[PaymentTerm]]("payment_terms"),
					Actions: []components.PageInterface{
						&components.TableButtonCreate{
							Link:        paymentTermCreateLinkWithQuery(PaymentTermTypeDueDate),
							Label:       "Due date",
							Page:        components.Page{Roles: []string{"superuser"}},
							Classes:     "btn-outline btn-sm",
							IconClasses: "mr-1",
						},
						&components.TableButtonCreate{
							Link:        paymentTermCreateLinkWithQuery(PaymentTermTypeRelative),
							Label:       "Relative",
							Page:        components.Page{Roles: []string{"superuser"}},
							Classes:     "btn-outline btn-sm",
							IconClasses: "mr-1",
						},
					},
					RowAttr: getters.RowAttrNavigate(lamu.RoutePath("finance_invoices.PaymentTermDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("$row.ID")),
					})),
					Columns: []components.TableColumn{
						{Label: "ID", Name: "ID", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Format("%d", getters.Any(getters.Key[uint]("$row.ID")))},
						}},
						{Label: "Kind", Name: "Type", Children: []components.PageInterface{
							&components.FieldText{Getter: registry.PairValueFromKey(getters.Key[string]("$row.Type"), paymentTermKindChoiceList)},
						}},
						{Label: "Summary", Name: "Summary", Children: []components.PageInterface{
							&components.FieldText{Getter: paymentTermRowSummaryGetter()},
						}},
					},
				},
			},
		}},
		{Key: "finance_invoices.PaymentTermCreateForm", Value: &components.ShellScaffold{
			Page:    components.Page{Roles: []string{"superuser"}},
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_accounts.MainMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name:      ptCreateName,
					ActionURL: lamu.RoutePath("finance_invoices.PaymentTermCreateRoute", nil),
					Children: []components.PageInterface{
						&components.FormComponent[PaymentTerm]{
							Attr:          getters.FormBubbling(ptCreateName),
							Title:         "Create payment term",
							Subtitle:      "Fixed: calendar due date. Relative: duration after invoice (Go duration syntax, e.g. 720h). Use toolbar shortcuts or change Kind.",
							ChildrenInput: paymentTermCreateFormInputs(),
							ChildrenAction: []components.PageInterface{
								&components.ButtonSubmit{Label: "Save"},
							},
						},
					},
				},
			},
		}},
		{Key: "finance_invoices.PaymentTermDetail", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_invoices.PaymentTermDetailMenu"}},
			Children: []components.PageInterface{
				&components.Detail[PaymentTerm]{
					Getter: getters.Key[PaymentTerm]("payment_term"),
					Children: []components.PageInterface{
						&components.ContainerColumn{
							Classes: "p-4",
							Children: []components.PageInterface{
								&components.LabelInline{Title: "Kind", Children: []components.PageInterface{
									&components.FieldText{Getter: getters.Key[string]("$in.Type")},
								}},
								&components.LabelInline{Title: "Summary", Children: []components.PageInterface{
									&components.FieldText{Getter: paymentTermDetailSummaryGetter()},
								}},
							},
						},
					},
				},
			},
		}},
		{Key: "finance_invoices.PaymentTermDeleteForm", Value: &components.Modal{
			Page: components.Page{Roles: []string{"superuser"}},
			UID:  "finance-payment-term-delete-modal",
			Children: []components.PageInterface{
				&components.DeleteConfirmation{
					Title:   "Delete payment term?",
					Message: "This removes the payment term and its backing definition. Invoices that reference it will be blocked until you change them.",
					Attr:    getters.FormBubbling(getters.Key[string]("$get.name")),
				},
			},
		}},
		{Key: "finance_invoices.PaymentTermDetailMenu", Value: &components.SidebarMenu{
			Title: getters.Format("Payment term #%d", getters.Any(getters.Key[uint]("payment_term.ID"))),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("All payment terms"),
				Url:   lamu.RoutePath("finance_invoices.PaymentTermListRoute", nil),
			},
			Children: []components.PageInterface{
				&components.SidebarMenuItem{
					Title: getters.Static("Detail"),
					Url: lamu.RoutePath("finance_invoices.PaymentTermDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("payment_term.ID")),
					}),
				},
				&components.SidebarMenuItem{
					Page:  components.Page{Roles: []string{"superuser"}},
					Title: getters.Static("Delete"),
					Url: lamu.RoutePath("finance_invoices.PaymentTermDeleteRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("payment_term.ID")),
					}),
				},
			},
		}},
	}
}
