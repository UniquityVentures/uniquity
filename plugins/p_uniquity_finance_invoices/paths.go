package p_uniquity_finance_invoices

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pluginRoutes() lamu.PluginFeatures[lamu.Route] {
	pt := AppUrl + "pt/"
	inv := AppUrl + "i/"
	posted := AppUrl + "posted/"
	cancelled := AppUrl + "cancelled/"
	return lamu.PluginFeatures[lamu.Route]{
		Entries: []registry.Pair[string, lamu.Route]{
			{Key: "finance_invoices.DefaultRoute", Value: lamu.Route{Path: AppUrl, Handler: lamu.NewDynamicView("finance_invoices.DraftInvoiceListView")}},
			{Key: "finance_invoices.DraftInvoiceCreateRoute", Value: lamu.Route{Path: AppUrl + "create/", Handler: lamu.NewDynamicView("finance_invoices.DraftInvoiceCreateView")}},
			{Key: "finance_invoices.DraftInvoiceDetailRoute", Value: lamu.Route{Path: inv + "{id}/", Handler: lamu.NewDynamicView("finance_invoices.DraftInvoiceDetailView")}},
			{Key: "finance_invoices.DraftInvoicePdfRoute", Value: lamu.Route{Path: inv + "{id}/pdf/", Handler: lamu.NewDynamicView("finance_invoices.DraftInvoicePdfView")}},
			{Key: "finance_invoices.DraftInvoiceUpdateRoute", Value: lamu.Route{Path: inv + "{id}/edit/", Handler: lamu.NewDynamicView("finance_invoices.DraftInvoiceUpdateView")}},
			{Key: "finance_invoices.DraftInvoiceDeleteRoute", Value: lamu.Route{Path: inv + "{id}/delete/", Handler: lamu.NewDynamicView("finance_invoices.DraftInvoiceDeleteView")}},
			{Key: "finance_invoices.DraftInvoicePostRoute", Value: lamu.Route{Path: inv + "{id}/post/", Handler: lamu.NewDynamicView("finance_invoices.DraftInvoicePostView")}},

			{Key: "finance_invoices.PostedInvoiceListRoute", Value: lamu.Route{Path: posted, Handler: lamu.NewDynamicView("finance_invoices.PostedInvoiceListView")}},
			{Key: "finance_invoices.PostedInvoiceDetailRoute", Value: lamu.Route{Path: posted + "{id}/", Handler: lamu.NewDynamicView("finance_invoices.PostedInvoiceDetailView")}},
			{Key: "finance_invoices.PostedInvoiceCancelRoute", Value: lamu.Route{Path: posted + "{id}/cancel/", Handler: lamu.NewDynamicView("finance_invoices.PostedInvoiceCancelView")}},

			{Key: "finance_invoices.CancelledInvoiceListRoute", Value: lamu.Route{Path: cancelled, Handler: lamu.NewDynamicView("finance_invoices.CancelledInvoiceListView")}},
			{Key: "finance_invoices.CancelledInvoiceDetailRoute", Value: lamu.Route{Path: cancelled + "{id}/", Handler: lamu.NewDynamicView("finance_invoices.CancelledInvoiceDetailView")}},
			{Key: "finance_invoices.CancelledInvoicePdfRoute", Value: lamu.Route{Path: cancelled + "{id}/pdf/", Handler: lamu.NewDynamicView("finance_invoices.CancelledInvoicePdfView")}},
			{Key: "finance_invoices.CancelledInvoiceNewDraftRoute", Value: lamu.Route{Path: cancelled + "{id}/new-draft/", Handler: lamu.NewDynamicView("finance_invoices.CancelledInvoiceNewDraftView")}},

			{Key: "finance_invoices.PaymentTermListRoute", Value: lamu.Route{Path: AppUrl + "payment-terms/", Handler: lamu.NewDynamicView("finance_invoices.PaymentTermListView")}},
			{Key: "finance_invoices.PaymentTermCreateRoute", Value: lamu.Route{Path: AppUrl + "payment-terms/create/", Handler: lamu.NewDynamicView("finance_invoices.PaymentTermCreateView")}},
			{Key: "finance_invoices.PaymentTermDetailRoute", Value: lamu.Route{Path: pt + "{id}/", Handler: lamu.NewDynamicView("finance_invoices.PaymentTermDetailView")}},
			{Key: "finance_invoices.PaymentTermDeleteRoute", Value: lamu.Route{Path: pt + "{id}/delete/", Handler: lamu.NewDynamicView("finance_invoices.PaymentTermDeleteView")}},
			{Key: "finance_invoices.PaymentTermFkSelectRoute", Value: lamu.Route{Path: AppUrl + "payment-terms/pick/", Handler: lamu.NewDynamicView("finance_invoices.PaymentTermFkSelectView")}},
		},
	}
}
