package p_uniquity_finance_invoices

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pluginModels() lamu.PluginFeatures[any] {
	return lamu.PluginFeatures[any]{
		Entries: []registry.Pair[string, any]{
			{Key: "p_uniquity_finance_invoices.PaymentTermDueDate", Value: PaymentTermDueDate{}},
			{Key: "p_uniquity_finance_invoices.PaymentTermRelative", Value: PaymentTermRelative{}},
			{Key: "p_uniquity_finance_invoices.PaymentTerm", Value: PaymentTerm{}},
			{Key: "p_uniquity_finance_invoices.Invoice", Value: Invoice{}},
			{Key: "p_uniquity_finance_invoices.InvoiceLine", Value: InvoiceLine{}},
		},
	}
}
