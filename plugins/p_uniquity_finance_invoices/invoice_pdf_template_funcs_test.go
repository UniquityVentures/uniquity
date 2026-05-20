package p_uniquity_finance_invoices

import (
	"testing"

	"github.com/UniquityVentures/lamu/getters"
	finance_taxes "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_taxes"
	"gorm.io/gorm"
)

func TestInvoiceGrandTotalWordsTemplate(t *testing.T) {
	inv := DraftInvoice{
		Lines: []DraftInvoiceLine{{
			Rate:     mustDec("27000"),
			Quantity: mustDec("2"),
		}},
		Taxes: []finance_taxes.Tax{
			{Model: gorm.Model{ID: 1}, Name: "SGST", Percentage: mustDec("9"), TaxType: finance_taxes.TaxKindLevied},
			{Model: gorm.Model{ID: 2}, Name: "CGST", Percentage: mustDec("9"), TaxType: finance_taxes.TaxKindLevied},
		},
	}
	got, err := invoiceGrandTotalWordsTemplate(getters.MapFromStruct(inv))
	if err != nil {
		t.Fatal(err)
	}
	want := "Sixty-Three Thousand Seven Hundred And Twenty Rupees"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestNum2wordsRupeesTemplate(t *testing.T) {
	got, err := num2wordsRupeesTemplate(63720)
	if err != nil {
		t.Fatal(err)
	}
	if got != "Sixty-Three Thousand Seven Hundred And Twenty Rupees" {
		t.Fatalf("got %q", got)
	}
}
