package p_uniquity_finance_invoices

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pluginRoutes() lamu.PluginFeatures[lamu.Route] {
	pt := AppUrl + "pt/"
	inv := AppUrl + "i/"
	return lamu.PluginFeatures[lamu.Route]{
		Entries: []registry.Pair[string, lamu.Route]{
			{Key: "finance_invoices.DefaultRoute", Value: lamu.Route{Path: AppUrl, Handler: lamu.NewDynamicView("finance_invoices.InvoiceListView")}},
			{Key: "finance_invoices.InvoiceCreateRoute", Value: lamu.Route{Path: AppUrl + "create/", Handler: lamu.NewDynamicView("finance_invoices.InvoiceCreateView")}},
			{Key: "finance_invoices.InvoiceDetailRoute", Value: lamu.Route{Path: inv + "{id}/", Handler: lamu.NewDynamicView("finance_invoices.InvoiceDetailView")}},
			{Key: "finance_invoices.PaymentTermListRoute", Value: lamu.Route{Path: AppUrl + "payment-terms/", Handler: lamu.NewDynamicView("finance_invoices.PaymentTermListView")}},
			{Key: "finance_invoices.PaymentTermCreateRoute", Value: lamu.Route{Path: AppUrl + "payment-terms/create/", Handler: lamu.NewDynamicView("finance_invoices.PaymentTermCreateView")}},
			{Key: "finance_invoices.PaymentTermDetailRoute", Value: lamu.Route{Path: pt + "{id}/", Handler: lamu.NewDynamicView("finance_invoices.PaymentTermDetailView")}},
			{Key: "finance_invoices.PaymentTermDeleteRoute", Value: lamu.Route{Path: pt + "{id}/delete/", Handler: lamu.NewDynamicView("finance_invoices.PaymentTermDeleteView")}},
			{Key: "finance_invoices.PaymentTermFkSelectRoute", Value: lamu.Route{Path: AppUrl + "payment-terms/pick/", Handler: lamu.NewDynamicView("finance_invoices.PaymentTermFkSelectView")}},
		},
	}
}
