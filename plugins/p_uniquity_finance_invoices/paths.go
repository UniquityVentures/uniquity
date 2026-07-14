package p_uniquity_finance_invoices

import (
	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/registry"
)

func pluginRoutes() lago.PluginFeatures[lago.Route] {
	pt := AppUrl + "pt/"
	inv := AppUrl + "i/"
	posted := AppUrl + "posted/"
	cancelled := AppUrl + "cancelled/"
	payments := AppUrl + "payments/"
	paidInv := AppUrl + "paid/"
	partInv := AppUrl + "partial/"
	return lago.PluginFeatures[lago.Route]{
		Entries: []registry.Pair[string, lago.Route]{
			{Key: "finance_invoices.DefaultRoute", Value: lago.Route{Path: AppUrl, Handler: lago.NewDynamicView("finance_invoices.DraftInvoiceListView")}},
			{Key: "finance_invoices.DraftInvoiceCreateRoute", Value: lago.Route{Path: AppUrl + "create/", Handler: lago.NewDynamicView("finance_invoices.DraftInvoiceCreateView")}},
			{Key: "finance_invoices.DraftInvoiceDetailRoute", Value: lago.Route{Path: inv + "{id}/", Handler: lago.NewDynamicView("finance_invoices.DraftInvoiceDetailView")}},
			{Key: "finance_invoices.DraftInvoicePdfRoute", Value: lago.Route{Path: inv + "{id}/pdf/", Handler: lago.NewDynamicView("finance_invoices.DraftInvoicePdfView")}},
			{Key: "finance_invoices.DraftInvoiceUpdateRoute", Value: lago.Route{Path: inv + "{id}/edit/", Handler: lago.NewDynamicView("finance_invoices.DraftInvoiceUpdateView")}},
			{Key: "finance_invoices.DraftInvoiceDeleteRoute", Value: lago.Route{Path: inv + "{id}/delete/", Handler: lago.NewDynamicView("finance_invoices.DraftInvoiceDeleteView")}},
			{Key: "finance_invoices.DraftInvoicePostRoute", Value: lago.Route{Path: inv + "{id}/post/", Handler: lago.NewDynamicView("finance_invoices.DraftInvoicePostView")}},

			{Key: "finance_invoices.PostedInvoiceListRoute", Value: lago.Route{Path: posted, Handler: lago.NewDynamicView("finance_invoices.PostedInvoiceListView")}},
			{Key: "finance_invoices.PostedInvoiceDetailRoute", Value: lago.Route{Path: posted + "{id}/", Handler: lago.NewDynamicView("finance_invoices.PostedInvoiceDetailView")}},
			{Key: "finance_invoices.PostedInvoicePdfRoute", Value: lago.Route{Path: posted + "{id}/pdf/", Handler: lago.NewDynamicView("finance_invoices.PostedInvoicePdfView")}},
			{Key: "finance_invoices.PostedInvoiceCancelRoute", Value: lago.Route{Path: posted + "{id}/cancel/", Handler: lago.NewDynamicView("finance_invoices.PostedInvoiceCancelView")}},

			{Key: "finance_invoices.CancelledInvoiceListRoute", Value: lago.Route{Path: cancelled, Handler: lago.NewDynamicView("finance_invoices.CancelledInvoiceListView")}},
			{Key: "finance_invoices.CancelledInvoiceDetailRoute", Value: lago.Route{Path: cancelled + "{id}/", Handler: lago.NewDynamicView("finance_invoices.CancelledInvoiceDetailView")}},
			{Key: "finance_invoices.CancelledInvoicePdfRoute", Value: lago.Route{Path: cancelled + "{id}/pdf/", Handler: lago.NewDynamicView("finance_invoices.CancelledInvoicePdfView")}},
			{Key: "finance_invoices.CancelledInvoiceNewDraftRoute", Value: lago.Route{Path: cancelled + "{id}/new-draft/", Handler: lago.NewDynamicView("finance_invoices.CancelledInvoiceNewDraftView")}},

			{Key: "finance_invoices.PaymentTermListRoute", Value: lago.Route{Path: AppUrl + "payment-terms/", Handler: lago.NewDynamicView("finance_invoices.PaymentTermListView")}},
			{Key: "finance_invoices.PaymentTermCreateRoute", Value: lago.Route{Path: AppUrl + "payment-terms/create/", Handler: lago.NewDynamicView("finance_invoices.PaymentTermCreateView")}},
			{Key: "finance_invoices.PaymentTermDetailRoute", Value: lago.Route{Path: pt + "{id}/", Handler: lago.NewDynamicView("finance_invoices.PaymentTermDetailView")}},
			{Key: "finance_invoices.PaymentTermDeleteRoute", Value: lago.Route{Path: pt + "{id}/delete/", Handler: lago.NewDynamicView("finance_invoices.PaymentTermDeleteView")}},
			{Key: "finance_invoices.PaymentTermFkSelectRoute", Value: lago.Route{Path: AppUrl + "payment-terms/pick/", Handler: lago.NewDynamicView("finance_invoices.PaymentTermFkSelectView")}},
			{Key: "finance_invoices.PostedInvoiceFkSelectRoute", Value: lago.Route{Path: AppUrl + "posted/pick/", Handler: lago.NewDynamicView("finance_invoices.PostedInvoiceFkSelectView")}},

			{Key: "finance_invoices.PaymentListRoute", Value: lago.Route{Path: payments, Handler: lago.NewDynamicView("finance_invoices.PaymentListView")}},
			{Key: "finance_invoices.PaymentCreateRoute", Value: lago.Route{Path: payments + "create/", Handler: lago.NewDynamicView("finance_invoices.PaymentCreateView")}},
			{Key: "finance_invoices.PaymentDetailRoute", Value: lago.Route{Path: payments + "{id}/", Handler: lago.NewDynamicView("finance_invoices.PaymentDetailView")}},

			{Key: "finance_invoices.PaidInvoiceListRoute", Value: lago.Route{Path: paidInv, Handler: lago.NewDynamicView("finance_invoices.PaidInvoiceListView")}},
			{Key: "finance_invoices.PaidInvoiceDetailRoute", Value: lago.Route{Path: paidInv + "{id}/", Handler: lago.NewDynamicView("finance_invoices.PaidInvoiceDetailView")}},
			{Key: "finance_invoices.PaidInvoicePdfRoute", Value: lago.Route{Path: paidInv + "{id}/pdf/", Handler: lago.NewDynamicView("finance_invoices.PaidInvoicePdfView")}},
			{Key: "finance_invoices.PartiallyPaidInvoiceListRoute", Value: lago.Route{Path: partInv, Handler: lago.NewDynamicView("finance_invoices.PartiallyPaidInvoiceListView")}},
			{Key: "finance_invoices.PartiallyPaidInvoiceDetailRoute", Value: lago.Route{Path: partInv + "{id}/", Handler: lago.NewDynamicView("finance_invoices.PartiallyPaidInvoiceDetailView")}},
			{Key: "finance_invoices.PartiallyPaidInvoicePdfRoute", Value: lago.Route{Path: partInv + "{id}/pdf/", Handler: lago.NewDynamicView("finance_invoices.PartiallyPaidInvoicePdfView")}},
		},
	}
}
