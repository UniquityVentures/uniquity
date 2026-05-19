package p_uniquity_finance_invoices

import (
	"context"

	finance_taxes "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_taxes"
)

type paymentCreateTaxesContextKey struct{}

// paymentTaxesFromContext returns withholding taxes parsed on the payment create form (before M2M replace).
func paymentTaxesFromContext(ctx context.Context) []finance_taxes.Tax {
	if ctx == nil {
		return nil
	}
	taxes, _ := ctx.Value(paymentCreateTaxesContextKey{}).([]finance_taxes.Tax)
	return taxes
}

func contextWithPaymentCreateTaxes(ctx context.Context, taxes []finance_taxes.Tax) context.Context {
	return context.WithValue(ctx, paymentCreateTaxesContextKey{}, taxes)
}
