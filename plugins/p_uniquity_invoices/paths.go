package p_uniquity_invoices

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pluginRoutes() lamu.PluginFeatures[lamu.Route] {
	inv := AppUrl + "i/"
	line := AppUrl + "l/"
	return lamu.PluginFeatures[lamu.Route]{
		Entries: []registry.Pair[string, lamu.Route]{
			{Key: "invoices.DefaultRoute", Value: lamu.Route{Path: AppUrl, Handler: lamu.NewDynamicView("invoices.InvoiceListView")}},
			{Key: "invoices.InvoiceCreateRoute", Value: lamu.Route{Path: AppUrl + "create/", Handler: lamu.NewDynamicView("invoices.InvoiceCreateView")}},
			{Key: "invoices.InvoiceDetailRoute", Value: lamu.Route{Path: inv + "{id}/", Handler: lamu.NewDynamicView("invoices.InvoiceDetailView")}},
			{Key: "invoices.InvoiceUpdateRoute", Value: lamu.Route{Path: inv + "{id}/edit/", Handler: lamu.NewDynamicView("invoices.InvoiceUpdateView")}},
			{Key: "invoices.InvoiceDeleteRoute", Value: lamu.Route{Path: inv + "{id}/delete/", Handler: lamu.NewDynamicView("invoices.InvoiceDeleteView")}},

			{Key: "invoices.InvoiceLineCreateRoute", Value: lamu.Route{Path: inv + "{invoiceId}/lines/create/", Handler: lamu.NewDynamicView("invoices.InvoiceLineCreateView")}},
			{Key: "invoices.InvoiceLineDetailRoute", Value: lamu.Route{Path: line + "{id}/", Handler: lamu.NewDynamicView("invoices.InvoiceLineDetailView")}},
			{Key: "invoices.InvoiceLineUpdateRoute", Value: lamu.Route{Path: line + "{id}/edit/", Handler: lamu.NewDynamicView("invoices.InvoiceLineUpdateView")}},
			{Key: "invoices.InvoiceLineDeleteRoute", Value: lamu.Route{Path: line + "{id}/delete/", Handler: lamu.NewDynamicView("invoices.InvoiceLineDeleteView")}},

			{Key: "invoices.ContactSelectRoute", Value: lamu.Route{Path: AppUrl + "contacts/select/", Handler: lamu.NewDynamicView("invoices.ContactSelectView")}},
			{Key: "invoices.PaymentTermSelectRoute", Value: lamu.Route{Path: AppUrl + "payment-terms/select/", Handler: lamu.NewDynamicView("invoices.PaymentTermSelectView")}},
		},
	}
}
