package p_uniquity_invoices

import (
	"log/slog"
	"net/http"

	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/plugins/p_users"
	"github.com/UniquityVentures/lamu/registry"
	"github.com/UniquityVentures/lamu/views"
)

// SuperuserOnlyLayer returns 401 unless the authenticated user is a superuser.
type SuperuserOnlyLayer struct{}

func (SuperuserOnlyLayer) Next(_ views.View, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := p_users.UserFromContext(r.Context(), "invoices.SuperuserOnlyLayer")
		if !user.IsSuperuser {
			slog.Error("invoices.SuperuserOnlyLayer: forbidden", "user_id", user.ID)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func pluginViews() lamu.PluginFeatures[*views.View] {
	auth := p_users.AuthenticationLayer{}
	su := SuperuserOnlyLayer{}

	invoicePreloads := views.QueryPatchers[Invoice]{
		{Key: "invoices.invoice_preload", Value: views.QueryPatcherPreload[Invoice]{Fields: []string{
			"Entity", "Partner", "Journal", "PaymentTerm", "Currency", "Move.Journal",
		}}},
	}

	invoiceListPreloads := views.QueryPatchers[Invoice]{
		{Key: "invoices.invoice_list_preload", Value: views.QueryPatcherPreload[Invoice]{Fields: []string{
			"Entity", "Partner", "Journal", "Currency",
		}}},
	}

	linePreloads := views.QueryPatchers[InvoiceLine]{
		{Key: "invoices.line_preload", Value: views.QueryPatcherPreload[InvoiceLine]{Fields: []string{
			"Invoice", "Product", "Account", "Taxes",
		}}},
	}

	lineListUnderInvoice := views.QueryPatchers[InvoiceLine]{
		{Key: "invoices.lines_under_invoice", Value: invoiceLinesScopedToInvoiceDetailURL{}},
		{Key: "invoices.line_list_preload", Value: views.QueryPatcherPreload[InvoiceLine]{Fields: []string{
			"Product", "Account",
		}}},
	}

	normalizers := views.FormPatchers{
		{Key: "invoices.form_normalize", Value: invoiceFormNormalizerPatcher{}},
	}

	lineInvoicePatcher := views.FormPatchers{
		{Key: "invoices.line_invoice_id", Value: invoiceLineInvoiceIDFormPatcher{}},
	}

	return lamu.PluginFeatures[*views.View]{
		Entries: []registry.Pair[string, *views.View]{
			{
				Key: "invoices.InvoiceListView",
				Value: lamu.GetPageView("invoices.InvoiceTable").
					WithLayer("p_users.auth", auth).
					WithLayer("invoices.superuser", su).
					WithLayer("invoices.invoice_list", views.LayerList[Invoice]{
						Key:           getters.Static("invoices"),
						QueryPatchers: invoiceListPreloads,
					}),
			},
			{
				Key: "invoices.InvoiceDetailView",
				Value: lamu.GetPageView("invoices.InvoiceDetail").
					WithLayer("p_users.auth", auth).
					WithLayer("invoices.superuser", su).
					WithLayer("invoices.invoice_detail", views.LayerDetail[Invoice]{
						Key:           getters.Static("invoice"),
						PathParamKey:  getters.Static("id"),
						QueryPatchers: invoicePreloads,
					}).
					WithLayer("invoices.invoice_lines", views.LayerList[InvoiceLine]{
						Key:           getters.Static("invoiceLines"),
						QueryPatchers: lineListUnderInvoice,
					}),
			},
			{
				Key: "invoices.InvoiceCreateView",
				Value: lamu.GetPageView("invoices.InvoiceCreateForm").
					WithLayer("p_users.auth", auth).
					WithLayer("invoices.superuser", su).
					WithLayer("invoices.invoice_create", views.LayerCreate[Invoice]{
						FormPatchers: normalizers,
						SuccessURL: lamu.RoutePath("invoices.InvoiceDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("$id")),
						}),
					}),
			},
			{
				Key: "invoices.InvoiceUpdateView",
				Value: lamu.GetPageView("invoices.InvoiceUpdateForm").
					WithLayer("p_users.auth", auth).
					WithLayer("invoices.superuser", su).
					WithLayer("invoices.invoice_update_detail", views.LayerDetail[Invoice]{
						Key:           getters.Static("invoice"),
						PathParamKey:  getters.Static("id"),
						QueryPatchers: invoicePreloads,
					}).
					WithLayer("invoices.invoice_update", views.LayerUpdate[Invoice]{
						Key:          getters.Static("invoice"),
						FormPatchers: normalizers,
						SuccessURL: lamu.RoutePath("invoices.InvoiceDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("invoice.ID")),
						}),
					}),
			},
			{
				Key: "invoices.InvoiceDeleteView",
				Value: lamu.GetPageView("invoices.InvoiceDeleteForm").
					WithLayer("p_users.auth", auth).
					WithLayer("invoices.superuser", su).
					WithLayer("invoices.invoice_delete_detail", views.LayerDetail[Invoice]{
						Key:           getters.Static("invoice"),
						PathParamKey:  getters.Static("id"),
						QueryPatchers: invoicePreloads,
					}).
					WithLayer("invoices.invoice_delete", views.LayerDelete[Invoice]{
						Key:        getters.Static("invoice"),
						SuccessURL: lamu.RoutePath("invoices.DefaultRoute", nil),
					}),
			},
			{
				Key: "invoices.ContactSelectView",
				Value: lamu.GetPageView("invoices.ContactSelectionTable").
					WithLayer("p_users.auth", auth).
					WithLayer("invoices.superuser", su).
					WithLayer("invoices.contact_select_list", views.LayerList[Contact]{
						Key: getters.Static("contacts"),
					}),
			},
			{
				Key: "invoices.PaymentTermSelectView",
				Value: lamu.GetPageView("invoices.PaymentTermSelectionTable").
					WithLayer("p_users.auth", auth).
					WithLayer("invoices.superuser", su).
					WithLayer("invoices.payment_term_select_list", views.LayerList[PaymentTerm]{
						Key: getters.Static("paymentTerms"),
					}),
			},
			{
				Key: "invoices.InvoiceLineCreateView",
				Value: lamu.GetPageView("invoices.InvoiceLineCreateForm").
					WithLayer("p_users.auth", auth).
					WithLayer("invoices.superuser", su).
					WithLayer("invoices.invoice_line_create_parent", views.LayerDetail[Invoice]{
						Key:           getters.Static("invoice"),
						PathParamKey:  getters.Static("invoiceId"),
						QueryPatchers: invoicePreloads,
					}).
					WithLayer("invoices.invoice_line_create", views.LayerCreate[InvoiceLine]{
						FormPatchers: lineInvoicePatcher,
						SuccessURL: lamu.RoutePath("invoices.InvoiceLineDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("$id")),
						}),
					}),
			},
			{
				Key: "invoices.InvoiceLineDetailView",
				Value: lamu.GetPageView("invoices.InvoiceLineDetail").
					WithLayer("p_users.auth", auth).
					WithLayer("invoices.superuser", su).
					WithLayer("invoices.invoice_line_detail", views.LayerDetail[InvoiceLine]{
						Key:           getters.Static("invoiceLine"),
						PathParamKey:  getters.Static("id"),
						QueryPatchers: linePreloads,
					}),
			},
			{
				Key: "invoices.InvoiceLineUpdateView",
				Value: lamu.GetPageView("invoices.InvoiceLineUpdateForm").
					WithLayer("p_users.auth", auth).
					WithLayer("invoices.superuser", su).
					WithLayer("invoices.invoice_line_update_detail", views.LayerDetail[InvoiceLine]{
						Key:           getters.Static("invoiceLine"),
						PathParamKey:  getters.Static("id"),
						QueryPatchers: linePreloads,
					}).
					WithLayer("invoices.invoice_line_update", views.LayerUpdate[InvoiceLine]{
						Key: getters.Static("invoiceLine"),
						SuccessURL: lamu.RoutePath("invoices.InvoiceLineDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("invoiceLine.ID")),
						}),
					}),
			},
			{
				Key: "invoices.InvoiceLineDeleteView",
				Value: lamu.GetPageView("invoices.InvoiceLineDeleteForm").
					WithLayer("p_users.auth", auth).
					WithLayer("invoices.superuser", su).
					WithLayer("invoices.invoice_line_delete_detail", views.LayerDetail[InvoiceLine]{
						Key:           getters.Static("invoiceLine"),
						PathParamKey:  getters.Static("id"),
						QueryPatchers: linePreloads,
					}).
					WithLayer("invoices.invoice_line_delete", views.LayerDelete[InvoiceLine]{
						Key: getters.Static("invoiceLine"),
						SuccessURL: lamu.RoutePath("invoices.InvoiceDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("invoiceLine.InvoiceID")),
						}),
					}),
			},
		},
	}
}
