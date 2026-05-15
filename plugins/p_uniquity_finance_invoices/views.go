package p_uniquity_finance_invoices

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
		user := p_users.UserFromContext(r.Context(), "finance_invoices.SuperuserOnlyLayer")
		if !user.IsSuperuser {
			slog.Error("finance_invoices.SuperuserOnlyLayer: forbidden", "user_id", user.ID)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func pluginViews() lamu.PluginFeatures[*views.View] {
	return lamu.PluginFeatures[*views.View]{
		Entries: []registry.Pair[string, *views.View]{
			{
				Key: "finance_invoices.InvoiceListView",
				Value: lamu.GetPageView("finance_invoices.InvoiceTable").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_invoices.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_invoices.invoice_list", views.LayerList[Invoice]{
						Key: getters.Static("invoices"),
						QueryPatchers: views.QueryPatchers[Invoice]{
							{Key: "finance_invoices.preload_customer", Value: views.QueryPatcherPreload[Invoice]{Fields: []string{"Customer", "PaymentTerm", "Taxes"}}},
						},
					}),
			},
			{
				Key: "finance_invoices.InvoiceDetailView",
				Value: lamu.GetPageView("finance_invoices.InvoiceDetail").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_invoices.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_invoices.invoice_detail", views.LayerDetail[Invoice]{
						Key:          getters.Static("invoice"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[Invoice]{
							{Key: "finance_invoices.preload_customer", Value: views.QueryPatcherPreload[Invoice]{Fields: []string{"Customer", "PaymentTerm", "Taxes", "Lines", "Lines.Product"}}},
						},
					}),
			},
			{
				Key: "finance_invoices.InvoiceCreateView",
				Value: lamu.GetPageView("finance_invoices.InvoiceCreateForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_invoices.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_invoices.invoice_create", views.LayerCreate[Invoice]{
						SuccessURL: lamu.RoutePath("finance_invoices.DefaultRoute", nil),
						FormPatchers: views.FormPatchers{
							{Key: "finance_invoices.invoice_create_lines", Value: invoiceCreateLinesPatcher{}},
						},
					}),
			},
			{
				Key: "finance_invoices.PaymentTermListView",
				Value: lamu.GetPageView("finance_invoices.PaymentTermTable").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_invoices.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_invoices.payment_term_list", views.LayerList[PaymentTerm]{
						Key: getters.Static("payment_terms"),
					}),
			},
			{
				Key: "finance_invoices.PaymentTermCreateView",
				Value: lamu.GetPageView("finance_invoices.PaymentTermCreateForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_invoices.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_invoices.payment_term_create", views.LayerCreate[PaymentTerm]{
						SuccessURL: lamu.RoutePath("finance_invoices.PaymentTermDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("$id")),
						}),
						FormPatchers: views.FormPatchers{
							{Key: "finance_invoices.payment_term_create_backing", Value: paymentTermCreateFormPatcher{}},
						},
					}),
			},
			{
				Key: "finance_invoices.PaymentTermDetailView",
				Value: lamu.GetPageView("finance_invoices.PaymentTermDetail").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_invoices.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_invoices.payment_term_detail", views.LayerDetail[PaymentTerm]{
						Key:          getters.Static("payment_term"),
						PathParamKey: getters.Static("id"),
					}),
			},
			{
				Key: "finance_invoices.PaymentTermDeleteView",
				Value: lamu.GetPageView("finance_invoices.PaymentTermDeleteForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_invoices.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_invoices.payment_term_detail", views.LayerDetail[PaymentTerm]{
						Key:          getters.Static("payment_term"),
						PathParamKey: getters.Static("id"),
					}).
					WithLayer("finance_invoices.payment_term_delete", views.LayerDelete[PaymentTerm]{
						Key:        getters.Static("payment_term"),
						SuccessURL: lamu.RoutePath("finance_invoices.PaymentTermListRoute", nil),
					}),
			},
			{
				Key: "finance_invoices.PaymentTermFkSelectView",
				Value: lamu.GetPageView("finance_invoices.PaymentTermFkSelectionTable").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_invoices.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_invoices.payment_term_fk_list", views.LayerList[PaymentTerm]{
						Key: getters.Static("payment_terms"),
					}),
			},
		},
	}
}
