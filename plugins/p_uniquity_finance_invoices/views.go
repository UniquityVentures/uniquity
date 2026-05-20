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
				Key: "finance_invoices.DraftInvoiceListView",
				Value: lamu.GetPageView("finance_invoices.InvoiceListHub").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_invoices.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_invoices.draft_invoice_list", views.LayerList[DraftInvoice]{
						Key: getters.Static("draft_invoices"),
						QueryPatchers: views.QueryPatchers[DraftInvoice]{
							registry.Pair[string, views.QueryPatcher[DraftInvoice]]{Key: "finance_invoices.list_fiscal_year_environment", Value: draftListFiscalYearEnvironment{}},
							registry.Pair[string, views.QueryPatcher[DraftInvoice]]{Key: "finance_invoices.list_datetime_range", Value: draftListDatetimeRange{}},
							registry.Pair[string, views.QueryPatcher[DraftInvoice]]{Key: "finance_invoices.list_exclude_posted", Value: draftListExcludePosted{}},
							registry.Pair[string, views.QueryPatcher[DraftInvoice]]{Key: "finance_invoices.preload_draft_list", Value: views.QueryPatcherPreload[DraftInvoice]{Fields: []string{"Customer", "PaymentTerm", "Taxes"}}},
						},
					}).
					WithLayer("finance_invoices.posted_invoice_list", views.LayerList[PostedInvoice]{
						Key: getters.Static("posted_invoices"),
						QueryPatchers: views.QueryPatchers[PostedInvoice]{
							registry.Pair[string, views.QueryPatcher[PostedInvoice]]{Key: "finance_invoices.posted_list_fy", Value: postedListFiscalYearEnvironment{}},
							registry.Pair[string, views.QueryPatcher[PostedInvoice]]{Key: "finance_invoices.posted_list_dt", Value: postedListDatetimeRange{}},
							registry.Pair[string, views.QueryPatcher[PostedInvoice]]{Key: "finance_invoices.posted_list_exclude_cancelled", Value: postedListExcludeCancelled{}},
							registry.Pair[string, views.QueryPatcher[PostedInvoice]]{Key: "finance_invoices.posted_list_exclude_paid", Value: postedListExcludeFullyPaid{}},
							registry.Pair[string, views.QueryPatcher[PostedInvoice]]{Key: "finance_invoices.preload_posted_list", Value: views.QueryPatcherPreload[PostedInvoice]{Fields: []string{"Customer"}}},
						},
					}).
					WithLayer("finance_invoices.cancelled_invoice_list", views.LayerList[CancelledInvoice]{
						Key: getters.Static("cancelled_invoices"),
						QueryPatchers: views.QueryPatchers[CancelledInvoice]{
							registry.Pair[string, views.QueryPatcher[CancelledInvoice]]{Key: "finance_invoices.cancelled_list_fy", Value: cancelledListFiscalYearEnvironment{}},
							registry.Pair[string, views.QueryPatcher[CancelledInvoice]]{Key: "finance_invoices.cancelled_list_dt", Value: cancelledListDatetimeRange{}},
							registry.Pair[string, views.QueryPatcher[CancelledInvoice]]{Key: "finance_invoices.preload_cancelled_list", Value: views.QueryPatcherPreload[CancelledInvoice]{Fields: []string{"Customer"}}},
						},
					}).
					WithLayer("finance_invoices.paid_invoice_list", views.LayerList[PaidInvoice]{
						Key: getters.Static("paid_invoices"),
						QueryPatchers: views.QueryPatchers[PaidInvoice]{
							registry.Pair[string, views.QueryPatcher[PaidInvoice]]{Key: "finance_invoices.preload_paid_invoice_list", Value: views.QueryPatcherPreload[PaidInvoice]{Fields: []string{"Payment", "PostedInvoice", "PostedInvoice.Customer"}}},
						},
					}).
					WithLayer("finance_invoices.partially_paid_invoice_list", views.LayerList[PartiallyPaidInvoice]{
						Key: getters.Static("partially_paid_invoices"),
						QueryPatchers: views.QueryPatchers[PartiallyPaidInvoice]{
							registry.Pair[string, views.QueryPatcher[PartiallyPaidInvoice]]{Key: "finance_invoices.preload_partially_paid_invoice_list", Value: views.QueryPatcherPreload[PartiallyPaidInvoice]{Fields: []string{"Payment", "PostedInvoice", "PostedInvoice.Customer"}}},
						},
					}).
					WithLayer("finance_invoices.toggle_draft_invoice_cols", views.LayerTableToggleColumns{
						QueryParam: getters.Static(invoiceDraftColsParam),
						ContextKey: getters.Static(invoiceDraftColsCtxKey),
					}).
					WithLayer("finance_invoices.toggle_posted_invoice_cols", views.LayerTableToggleColumns{
						QueryParam: getters.Static(invoicePostedColsParam),
						ContextKey: getters.Static(invoicePostedColsCtxKey),
					}).
					WithLayer("finance_invoices.toggle_cancelled_invoice_cols", views.LayerTableToggleColumns{
						QueryParam: getters.Static(invoiceCancelledColsParam),
						ContextKey: getters.Static(invoiceCancelledColsCtxKey),
					}),
			},
			{
				Key: "finance_invoices.DraftInvoiceDetailView",
				Value: lamu.GetPageView("finance_invoices.DraftInvoiceDetail").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_invoices.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_invoices.draft_invoice_detail", views.LayerDetail[DraftInvoice]{
						Key:          getters.Static("draft_invoice"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[DraftInvoice]{
							registry.Pair[string, views.QueryPatcher[DraftInvoice]]{Key: "finance_invoices.preload_draft_detail", Value: views.QueryPatcherPreload[DraftInvoice]{Fields: []string{"Customer", "PaymentTerm", "Taxes", "Lines", "Lines.Product", "Lines.Taxes"}}},
						},
					}),
			},
			{
				Key: "finance_invoices.DraftInvoicePdfView",
				Value: lamu.GetPageView("finance_invoices.DraftInvoiceDetail").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_invoices.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_invoices.draft_invoice_pdf_detail", views.LayerDetail[DraftInvoice]{
						Key:          getters.Static("draft_invoice"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[DraftInvoice]{
							registry.Pair[string, views.QueryPatcher[DraftInvoice]]{Key: "finance_invoices.preload_draft_detail", Value: views.QueryPatcherPreload[DraftInvoice]{Fields: []string{"Customer", "PaymentTerm", "Taxes", "Lines", "Lines.Product", "Lines.Taxes"}}},
						},
					}).
					WithLayer("finance_invoices.draft_invoice_pdf", views.MethodLayer{
						Method:  http.MethodGet,
						Handler: draftInvoicePdfHandler,
					}),
			},
			{
				Key: "finance_invoices.DraftInvoiceCreateView",
				Value: lamu.GetPageView("finance_invoices.DraftInvoiceCreateForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_invoices.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_invoices.draft_invoice_create", views.LayerCreate[DraftInvoice]{
						SuccessURL: lamu.RoutePath("finance_invoices.DefaultRoute", nil),
						FormPatchers: views.FormPatchers{
							registry.Pair[string, views.FormPatcher]{Key: "finance_invoices.draft_invoice_create_lines", Value: invoiceCreateLinesPatcher{}},
						},
					}),
			},
			{
				Key: "finance_invoices.DraftInvoiceUpdateView",
				Value: lamu.GetPageView("finance_invoices.DraftInvoiceUpdateForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_invoices.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_invoices.draft_invoice_detail", views.LayerDetail[DraftInvoice]{
						Key:          getters.Static("draft_invoice"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[DraftInvoice]{
							registry.Pair[string, views.QueryPatcher[DraftInvoice]]{Key: "finance_invoices.preload_draft_detail", Value: views.QueryPatcherPreload[DraftInvoice]{Fields: []string{"Customer", "PaymentTerm", "Taxes", "Lines", "Lines.Product", "Lines.Taxes"}}},
						},
					}).
					WithLayer("finance_invoices.draft_invoice_update", views.LayerUpdate[DraftInvoice]{
						Key: getters.Static("draft_invoice"),
						SuccessURL: lamu.RoutePath("finance_invoices.DraftInvoiceDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("draft_invoice.ID")),
						}),
						FormPatchers: views.FormPatchers{
							registry.Pair[string, views.FormPatcher]{Key: "finance_invoices.draft_invoice_update_lines", Value: invoiceCreateLinesPatcher{}},
						},
					}),
			},
			{
				Key: "finance_invoices.DraftInvoiceDeleteView",
				Value: lamu.GetPageView("finance_invoices.DraftInvoiceDeleteForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_invoices.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_invoices.draft_invoice_detail", views.LayerDetail[DraftInvoice]{
						Key:          getters.Static("draft_invoice"),
						PathParamKey: getters.Static("id"),
					}).
					WithLayer("finance_invoices.draft_invoice_delete", views.LayerDelete[DraftInvoice]{
						Key:        getters.Static("draft_invoice"),
						SuccessURL: lamu.RoutePath("finance_invoices.DefaultRoute", nil),
					}),
			},
			{
				Key: "finance_invoices.DraftInvoicePostView",
				Value: lamu.GetPageView("finance_invoices.DraftInvoicePostForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_invoices.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_invoices.draft_invoice_detail", views.LayerDetail[DraftInvoice]{
						Key:          getters.Static("draft_invoice"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[DraftInvoice]{
							registry.Pair[string, views.QueryPatcher[DraftInvoice]]{Key: "finance_invoices.preload_draft_detail", Value: views.QueryPatcherPreload[DraftInvoice]{Fields: []string{"Lines"}}},
						},
					}).
					WithLayer("finance_invoices.draft_invoice_post", layerPostDraftInvoice{}),
			},

			{
				Key:   "finance_invoices.PostedInvoiceListView",
				Value: invoiceHubRedirectView("posted"),
			},
			{
				Key:   "finance_invoices.CancelledInvoiceListView",
				Value: invoiceHubRedirectView("cancelled"),
			},
			{
				Key: "finance_invoices.PostedInvoiceDetailView",
				Value: lamu.GetPageView("finance_invoices.PostedInvoiceDetail").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_invoices.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_invoices.posted_invoice_detail", views.LayerDetail[PostedInvoice]{
						Key:          getters.Static("posted_invoice"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[PostedInvoice]{
							registry.Pair[string, views.QueryPatcher[PostedInvoice]]{Key: "finance_invoices.preload_posted_detail", Value: views.QueryPatcherPreload[PostedInvoice]{Fields: []string{"Customer", "PaymentTerm", "Taxes", "Lines", "Lines.Product", "Lines.Taxes", "JournalEntry"}}},
						},
					}),
			},
			{
				Key: "finance_invoices.PostedInvoicePdfView",
				Value: lamu.GetPageView("finance_invoices.PostedInvoiceDetail").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_invoices.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_invoices.posted_invoice_pdf_detail", views.LayerDetail[PostedInvoice]{
						Key:          getters.Static("posted_invoice"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[PostedInvoice]{
							registry.Pair[string, views.QueryPatcher[PostedInvoice]]{Key: "finance_invoices.preload_posted_pdf", Value: views.QueryPatcherPreload[PostedInvoice]{Fields: postedInvoicePdfPreload}},
						},
					}).
					WithLayer("finance_invoices.posted_invoice_pdf", views.MethodLayer{
						Method:  http.MethodGet,
						Handler: postedInvoicePdfHandler,
					}),
			},
			{
				Key: "finance_invoices.PostedInvoiceCancelView",
				Value: lamu.GetPageView("finance_invoices.PostedInvoiceCancelForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_invoices.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_invoices.posted_invoice_detail", views.LayerDetail[PostedInvoice]{
						Key:          getters.Static("posted_invoice"),
						PathParamKey: getters.Static("id"),
					}).
					WithLayer("finance_invoices.posted_invoice_cancel", layerCancelPostedInvoice{}),
			},

			{
				Key: "finance_invoices.CancelledInvoiceDetailView",
				Value: lamu.GetPageView("finance_invoices.CancelledInvoiceDetail").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_invoices.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_invoices.cancelled_invoice_detail", views.LayerDetail[CancelledInvoice]{
						Key:          getters.Static("cancelled_invoice"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[CancelledInvoice]{
							registry.Pair[string, views.QueryPatcher[CancelledInvoice]]{Key: "finance_invoices.preload_cancelled_detail", Value: views.QueryPatcherPreload[CancelledInvoice]{Fields: []string{"Customer", "CreditNote", "Taxes", "Lines", "Lines.Product", "Lines.Taxes"}}},
						},
					}),
			},
			{
				Key: "finance_invoices.CancelledInvoicePdfView",
				Value: lamu.GetPageView("finance_invoices.CancelledInvoiceDetail").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_invoices.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_invoices.cancelled_invoice_pdf_detail", views.LayerDetail[CancelledInvoice]{
						Key:          getters.Static("cancelled_invoice"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[CancelledInvoice]{
							registry.Pair[string, views.QueryPatcher[CancelledInvoice]]{Key: "finance_invoices.preload_cancelled_detail", Value: views.QueryPatcherPreload[CancelledInvoice]{Fields: []string{"Customer", "CreditNote", "Taxes", "Lines", "Lines.Product", "Lines.Taxes"}}},
						},
					}).
					WithLayer("finance_invoices.cancelled_invoice_pdf", views.MethodLayer{
						Method:  http.MethodGet,
						Handler: cancelledInvoicePdfHandler,
					}),
			},
			{
				Key: "finance_invoices.CancelledInvoiceNewDraftView",
				Value: lamu.GetPageView("finance_invoices.CancelledInvoiceNewDraftForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_invoices.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_invoices.cancelled_invoice_detail", views.LayerDetail[CancelledInvoice]{
						Key:          getters.Static("cancelled_invoice"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[CancelledInvoice]{
							registry.Pair[string, views.QueryPatcher[CancelledInvoice]]{Key: "finance_invoices.preload_cancelled_newdraft", Value: views.QueryPatcherPreload[CancelledInvoice]{Fields: []string{"Lines", "Lines.Taxes", "Taxes"}}},
						},
					}).
					WithLayer("finance_invoices.cancelled_new_draft", layerNewDraftFromCancelled{}),
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
							registry.Pair[string, views.FormPatcher]{Key: "finance_invoices.payment_term_create_backing", Value: paymentTermCreateFormPatcher{}},
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

			{
				Key: "finance_invoices.PostedInvoiceFkSelectView",
				Value: lamu.GetPageView("finance_invoices.PostedInvoiceFkSelectionTable").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_invoices.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_invoices.posted_invoice_fk_list", views.LayerList[PostedInvoice]{
						Key: getters.Static("posted_invoices"),
						QueryPatchers: views.QueryPatchers[PostedInvoice]{
							registry.Pair[string, views.QueryPatcher[PostedInvoice]]{Key: "finance_invoices.posted_invoice_pick_eligible", Value: postedInvoicePickEligible{}},
							registry.Pair[string, views.QueryPatcher[PostedInvoice]]{Key: "finance_invoices.preload_posted_fk_pick", Value: views.QueryPatcherPreload[PostedInvoice]{Fields: []string{"Customer"}}},
						},
					}),
			},
			{
				Key: "finance_invoices.PaymentListView",
				Value: lamu.GetPageView("finance_invoices.PaymentTable").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_invoices.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_invoices.payment_list", views.LayerList[Payment]{
						Key: getters.Static("payments"),
						QueryPatchers: views.QueryPatchers[Payment]{
							registry.Pair[string, views.QueryPatcher[Payment]]{Key: "finance_invoices.preload_payment_list", Value: views.QueryPatcherPreload[Payment]{Fields: []string{"PostedInvoice", "Account", "Taxes"}}},
						},
					}),
			},
			{
				Key: "finance_invoices.PaymentCreateView",
				Value: lamu.GetPageView("finance_invoices.PaymentCreateForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_invoices.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_invoices.payment_create_query_defaults", paymentCreateQueryDefaultsLayer{}).
					WithLayer("finance_invoices.payment_create", views.LayerCreate[Payment]{
						FormPatchers: views.FormPatchers{
							registry.Pair[string, views.FormPatcher]{Key: "finance_invoices.payment_create_taxes", Value: paymentCreateFormPatcher{}},
						},
						SuccessURL: lamu.RoutePath("finance_invoices.PaymentDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("$id")),
						}),
					}),
			},
			{
				Key: "finance_invoices.PaymentDetailView",
				Value: lamu.GetPageView("finance_invoices.PaymentDetail").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_invoices.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_invoices.payment_detail", views.LayerDetail[Payment]{
						Key:          getters.Static("payment"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[Payment]{
							registry.Pair[string, views.QueryPatcher[Payment]]{Key: "finance_invoices.preload_payment_detail", Value: views.QueryPatcherPreload[Payment]{Fields: []string{"PostedInvoice", "PostedInvoice.Customer", "Account", "JournalEntry", "Taxes"}}},
						},
					}),
			},

			{
				Key:   "finance_invoices.PaidInvoiceListView",
				Value: invoiceHubRedirectView("paid"),
			},
			{
				Key: "finance_invoices.PaidInvoiceDetailView",
				Value: lamu.GetPageView("finance_invoices.PaidInvoiceDetail").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_invoices.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_invoices.paid_invoice_detail", views.LayerDetail[PaidInvoice]{
						Key:          getters.Static("paid_invoice"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[PaidInvoice]{
							registry.Pair[string, views.QueryPatcher[PaidInvoice]]{Key: "finance_invoices.preload_paid_invoice_detail", Value: views.QueryPatcherPreload[PaidInvoice]{Fields: settlementPostedInvoiceDetailPreload}},
						},
					}),
			},
			{
				Key: "finance_invoices.PaidInvoicePdfView",
				Value: lamu.GetPageView("finance_invoices.PaidInvoiceDetail").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_invoices.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_invoices.paid_invoice_pdf_detail", views.LayerDetail[PaidInvoice]{
						Key:          getters.Static("paid_invoice"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[PaidInvoice]{
							registry.Pair[string, views.QueryPatcher[PaidInvoice]]{Key: "finance_invoices.preload_paid_invoice_pdf", Value: views.QueryPatcherPreload[PaidInvoice]{Fields: settlementPostedInvoiceDetailPreload}},
						},
					}).
					WithLayer("finance_invoices.paid_invoice_pdf", views.MethodLayer{
						Method:  http.MethodGet,
						Handler: paidInvoicePdfHandler,
					}),
			},
			{
				Key:   "finance_invoices.PartiallyPaidInvoiceListView",
				Value: invoiceHubRedirectView("partial"),
			},
			{
				Key: "finance_invoices.PartiallyPaidInvoiceDetailView",
				Value: lamu.GetPageView("finance_invoices.PartiallyPaidInvoiceDetail").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_invoices.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_invoices.partially_paid_invoice_detail", views.LayerDetail[PartiallyPaidInvoice]{
						Key:          getters.Static("partially_paid_invoice"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[PartiallyPaidInvoice]{
							registry.Pair[string, views.QueryPatcher[PartiallyPaidInvoice]]{Key: "finance_invoices.preload_partially_paid_invoice_detail", Value: views.QueryPatcherPreload[PartiallyPaidInvoice]{Fields: settlementPostedInvoiceDetailPreload}},
						},
					}),
			},
			{
				Key: "finance_invoices.PartiallyPaidInvoicePdfView",
				Value: lamu.GetPageView("finance_invoices.PartiallyPaidInvoiceDetail").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("finance_invoices.superuser", SuperuserOnlyLayer{}).
					WithLayer("finance_invoices.partially_paid_invoice_pdf_detail", views.LayerDetail[PartiallyPaidInvoice]{
						Key:          getters.Static("partially_paid_invoice"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[PartiallyPaidInvoice]{
							registry.Pair[string, views.QueryPatcher[PartiallyPaidInvoice]]{Key: "finance_invoices.preload_partially_paid_invoice_pdf", Value: views.QueryPatcherPreload[PartiallyPaidInvoice]{Fields: settlementPostedInvoiceDetailPreload}},
						},
					}).
					WithLayer("finance_invoices.partially_paid_invoice_pdf", views.MethodLayer{
						Method:  http.MethodGet,
						Handler: partiallyPaidInvoicePdfHandler,
					}),
			},
		},
		Patches: []registry.Pair[string, func(*views.View) *views.View]{
			{Key: "finance_accounts.AccountingPreferencesView", Value: patchAccountingPreferencesView},
		},
	}
}

func invoiceHubRedirectView(tab string) *views.View {
	return lamu.RedirectView(invoiceHubURLWithTabGetter(tab)).
		WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
		WithLayer("finance_invoices.superuser", SuperuserOnlyLayer{})
}
