package p_uniquity_invoices

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pluginModels() lamu.PluginFeatures[any] {
	return lamu.PluginFeatures[any]{
		Entries: []registry.Pair[string, any]{
			{Key: "p_uniquity_invoices.Contact", Value: Contact{}},
			{Key: "p_uniquity_invoices.PaymentTerm", Value: PaymentTerm{}},
			{Key: "p_uniquity_invoices.Invoice", Value: Invoice{}},
			{Key: "p_uniquity_invoices.InvoiceLine", Value: InvoiceLine{}},
		},
	}
}
