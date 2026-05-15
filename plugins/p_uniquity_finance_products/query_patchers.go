package p_uniquity_finance_products

import "github.com/UniquityVentures/lamu/views"

// productPreloadTaxes loads M2M taxes for forms and detail views.
var productPreloadTaxes = views.QueryPatcherPreload[Product]{Fields: []string{"Taxes"}}
