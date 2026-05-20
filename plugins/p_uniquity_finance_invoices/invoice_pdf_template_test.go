package p_uniquity_finance_invoices

import (
	"bytes"
	"os"
	"strings"
	"text/template"
	"testing"
	"time"

	"github.com/UniquityVentures/lamu/fields"
	"github.com/UniquityVentures/lamu/getters"
	finance_customer "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_customer"
	finance_products "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_products"
	finance_taxes "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_taxes"
	"github.com/francescoalemanno/gotypst"
	"gorm.io/gorm"
)

func TestExampleInvoicePDFTemplateCompiles(t *testing.T) {
	b, err := os.ReadFile("../p_uniquity_finance_accounts/example_invoice_pdf_template.typ.tmpl")
	if err != nil {
		t.Fatal(err)
	}
	num := "INV/25-26/0051"
	inv := DraftInvoice{
		Model:    gorm.Model{ID: 1},
		Number:   &num,
		Datetime: time.Date(2026, 2, 8, 0, 0, 0, 0, time.UTC),
		Customer: finance_customer.Customer{
			Name:    "Wipro PARI Robotics Private Limited",
			Address: "Gat no 463A2-8 to 463A2-11 partial, 463A2-15 and 463A2-16,\nDhangarwadi Village, Dist. Satara, Taluka Khandala,\nShirwal 412801\nMaharashtra MH\nIndia",
			GSTIN:   "27AADCW2907L1Z2",
		},
		Lines: []DraftInvoiceLine{{
			Product: finance_products.Product{
				Name:      "Gandola rental charges for one month",
				HSNCode:   9973,
				Reference: "Site: Shree Jii",
			},
			Rate:     mustDec("27000"),
			Quantity: mustDec("2"),
			Taxes: []finance_taxes.Tax{{
				Name:       "GST 18%",
				Percentage: mustDec("18"),
				TaxType:    finance_taxes.TaxKindLevied,
			}},
		}},
		Taxes: []finance_taxes.Tax{
			{Name: "SGST", Percentage: mustDec("9"), TaxType: finance_taxes.TaxKindLevied},
			{Name: "CGST", Percentage: mustDec("9"), TaxType: finance_taxes.TaxKindLevied},
		},
	}
	tmpl, err := template.New("invoice_pdf").Funcs(invoicePDFTemplateFuncs()).Parse(string(b))
	if err != nil {
		t.Fatal(err)
	}
	var typstBuf bytes.Buffer
	if err := tmpl.Execute(&typstBuf, getters.MapFromStruct(inv)); err != nil {
		t.Fatal(err)
	}
	out := typstBuf.String()
	if !strings.Contains(out, "Sixty-Three Thousand Seven Hundred And Twenty Rupees") {
		t.Fatalf("expected amount in words from num2words, got:\n%s", out)
	}
	if _, err := gotypst.PDF(typstBuf.Bytes()); err != nil {
		t.Fatalf("gotypst: %v\n--- typst ---\n%s", err, out)
	}
}

func mustDec(s string) fields.DecimalSix {
	var d fields.DecimalSix
	if err := d.UnmarshalText([]byte(s)); err != nil {
		panic(err)
	}
	return d
}
