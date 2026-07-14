package p_uniquity_finance_invoices

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	finance_taxes "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_taxes"
	"github.com/divan/num2words"
	"github.com/lariv-in/lago/fields"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// invoicePDFTemplateFuncs returns Go text/template helpers for invoice PDF Typst sources.
// Available in templates: num2words, num2wordsAnd, num2wordsRupees, invoiceGrandTotalWords.
func invoicePDFTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"num2words":              num2wordsTemplate,
		"num2wordsAnd":           num2wordsAndTemplate,
		"num2wordsRupees":        num2wordsRupeesTemplate,
		"invoiceGrandTotalWords": invoiceGrandTotalWordsTemplate,
		"urlImage":               urlImageTemplate,
	}
}

func num2wordsTemplate(v any) (string, error) {
	n, err := coerceTemplateInt(v)
	if err != nil {
		return "", err
	}
	return num2words.Convert(n), nil
}

func num2wordsAndTemplate(v any) (string, error) {
	n, err := coerceTemplateInt(v)
	if err != nil {
		return "", err
	}
	return num2words.ConvertAnd(n), nil
}

func num2wordsRupeesTemplate(v any) (string, error) {
	n, err := coerceTemplateInt(v)
	if err != nil {
		return "", err
	}
	return invoiceAmountWords(n), nil
}

func invoiceGrandTotalWordsTemplate(root any) (string, error) {
	m, ok := root.(map[string]any)
	if !ok || m == nil {
		return "", fmt.Errorf("invoiceGrandTotalWords: expected map root")
	}
	grand, err := invoicePDFReceivableGrandTotal(m)
	if err != nil {
		return "", err
	}
	return invoiceAmountWordsFromDecimal(grand), nil
}

func invoiceAmountWords(amount int) string {
	if amount < 0 {
		return titleInvoiceWords(num2words.ConvertAnd(amount)) + " Rupees"
	}
	return titleInvoiceWords(num2words.ConvertAnd(amount)) + " Rupees"
}

func invoiceAmountWordsFromDecimal(d fields.DecimalSix) string {
	n, err := decimalSixRoundedInt(d)
	if err != nil {
		return ""
	}
	return invoiceAmountWords(n)
}

func titleInvoiceWords(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	title := cases.Title(language.English)
	parts := strings.Fields(s)
	for i, p := range parts {
		if strings.EqualFold(p, "and") {
			parts[i] = "And"
			continue
		}
		parts[i] = title.String(p)
	}
	return strings.Join(parts, " ")
}

func invoicePDFReceivableGrandTotal(root map[string]any) (fields.DecimalSix, error) {
	headerTaxes, _ := root["Taxes"].([]finance_taxes.Tax)
	switch lines := root["Lines"].(type) {
	case []DraftInvoiceLine:
		totals, lineTaxIDs := accumulateInvoiceLineTotals(lines)
		return invoiceReceivableGrandTotal(totals, headerTaxes, lineTaxIDs), nil
	case []PostedInvoiceLine:
		totals, lineTaxIDs := accumulatePostedInvoiceLineTotals(lines)
		return invoiceReceivableGrandTotal(totals, headerTaxes, lineTaxIDs), nil
	case []CancelledInvoiceLine:
		totals, lineTaxIDs := accumulateCancelledInvoiceLineTotals(lines)
		return invoiceReceivableGrandTotal(totals, headerTaxes, lineTaxIDs), nil
	default:
		return fields.DecimalSix{R: big.NewRat(0, 1)}, nil
	}
}

func coerceTemplateInt(v any) (int, error) {
	switch n := v.(type) {
	case int:
		return n, nil
	case int64:
		return int(n), nil
	case int32:
		return int(n), nil
	case uint:
		return int(n), nil
	case uint64:
		return int(n), nil
	case float64:
		return int(math.Round(n)), nil
	case float32:
		return int(math.Round(float64(n))), nil
	case fields.DecimalSix:
		return decimalSixRoundedInt(n)
	case string:
		var d fields.DecimalSix
		if err := d.UnmarshalText([]byte(n)); err != nil {
			return 0, fmt.Errorf("num2words: invalid number %q", n)
		}
		return decimalSixRoundedInt(d)
	default:
		return 0, fmt.Errorf("num2words: unsupported type %T", v)
	}
}

func decimalSixRoundedInt(d fields.DecimalSix) (int, error) {
	if d.R == nil {
		return 0, nil
	}
	f, _ := d.R.Float64()
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return 0, fmt.Errorf("num2words: non-finite decimal")
	}
	return int(math.Round(f)), nil
}

func urlImageTemplate(urlStr string) (string, error) {
	if urlStr == "" {
		return "", fmt.Errorf("urlImage: empty URL")
	}
	h := sha256.New()
	h.Write([]byte(urlStr))
	hashName := hex.EncodeToString(h.Sum(nil))

	ext := filepath.Ext(urlStr)
	if idx := strings.Index(ext, "?"); idx != -1 {
		ext = ext[:idx]
	}
	if ext == "" {
		ext = ".png"
	}

	filename := hashName + ext
	tmpPath := filepath.Join("/tmp", filename)

	if _, err := os.Stat(tmpPath); err == nil {
		return filename, nil
	}

	writeFallback := func() {
		img := image.NewRGBA(image.Rect(0, 0, 1, 1))
		img.Set(0, 0, color.Transparent)
		f, err := os.Create(tmpPath)
		if err == nil {
			_ = png.Encode(f, img)
			f.Close()
		}
	}

	resp, err := http.Get(urlStr)
	if err != nil || resp.StatusCode != http.StatusOK {
		writeFallback()
		return filename, nil
	}
	defer resp.Body.Close()

	out, err := os.Create(tmpPath)
	if err != nil {
		writeFallback()
		return filename, nil
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		writeFallback()
		return filename, nil
	}

	return filename, nil
}
