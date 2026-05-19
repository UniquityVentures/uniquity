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
			{Key: "p_uniquity_finance_invoices.DraftInvoice", Value: DraftInvoice{}},
			{Key: "p_uniquity_finance_invoices.DraftInvoiceLine", Value: DraftInvoiceLine{}},
			{Key: "p_uniquity_finance_invoices.PostedInvoice", Value: PostedInvoice{}},
			{Key: "p_uniquity_finance_invoices.PostedInvoiceLine", Value: PostedInvoiceLine{}},
			{Key: "p_uniquity_finance_invoices.CancelledInvoice", Value: CancelledInvoice{}},
			{Key: "p_uniquity_finance_invoices.CancelledInvoiceLine", Value: CancelledInvoiceLine{}},
			{Key: "p_uniquity_finance_invoices.Payment", Value: Payment{}},
			{Key: "p_uniquity_finance_invoices.PartiallyPaidInvoice", Value: PartiallyPaidInvoice{}},
			{Key: "p_uniquity_finance_invoices.PaidInvoice", Value: PaidInvoice{}},
		},
	}
}
