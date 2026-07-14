package p_uniquity_finance_products

import "github.com/lariv-in/lago/views"

// productPreloadTaxes loads M2M taxes for forms and detail views.
var productPreloadTaxes = views.QueryPatcherPreload[Product]{Fields: []string{
	"Taxes",
}}
