package p_uniquity_finance_invoices

import (
	"fmt"
	"math/big"
	"strings"

	finance_taxes "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_taxes"
	"github.com/lariv-in/lago/fields"
)

func effectiveTaxKind(t finance_taxes.Tax) finance_taxes.TaxKind {
	if t.TaxType == finance_taxes.TaxKindWithholding {
		return finance_taxes.TaxKindWithholding
	}
	return finance_taxes.TaxKindLevied
}

func taxesOfKind(taxes []finance_taxes.Tax, kind finance_taxes.TaxKind) []finance_taxes.Tax {
	var out []finance_taxes.Tax
	for _, t := range taxes {
		if effectiveTaxKind(t) == kind {
			out = append(out, t)
		}
	}
	return out
}

func taxesLevied(taxes []finance_taxes.Tax) []finance_taxes.Tax {
	return taxesOfKind(taxes, finance_taxes.TaxKindLevied)
}

func taxesWithholding(taxes []finance_taxes.Tax) []finance_taxes.Tax {
	return taxesOfKind(taxes, finance_taxes.TaxKindWithholding)
}

func sumTaxPercents(taxes []finance_taxes.Tax) fields.DecimalSix {
	acc := big.NewRat(0, 1)
	for _, t := range taxes {
		if t.Percentage.R != nil {
			acc.Add(acc, t.Percentage.R)
		}
	}
	return fields.DecimalSix{R: acc}.NormalizeDecimals()
}

func taxAmountOnBase(base, pctSum fields.DecimalSix) fields.DecimalSix {
	if base.R == nil || base.R.Sign() == 0 {
		return fields.DecimalSix{R: big.NewRat(0, 1)}.NormalizeDecimals()
	}
	if pctSum.R == nil || pctSum.R.Sign() == 0 {
		return fields.DecimalSix{R: big.NewRat(0, 1)}.NormalizeDecimals()
	}
	hundred := big.NewRat(100, 1)
	r := new(big.Rat).Quo(pctSum.R, hundred)
	return decMul(base, fields.DecimalSix{R: r}).NormalizeDecimals()
}

func taxAmountForTax(base fields.DecimalSix, tax finance_taxes.Tax) fields.DecimalSix {
	return taxAmountOnBase(base, sumTaxPercents([]finance_taxes.Tax{tax}))
}

// invoiceLineAmountBreakdown returns untaxed base, levied tax, withholding tax, and net line total
// (untaxed + levied − withholding).
func invoiceLineAmountBreakdown(qty, rate fields.DecimalSix, taxes []finance_taxes.Tax) (untaxed, leviedAmt, withholdingAmt, netTotal fields.DecimalSix) {
	untaxed = decMul(qty, rate)
	leviedAmt = taxAmountOnBase(untaxed, sumTaxPercents(taxesLevied(taxes)))
	withholdingAmt = taxAmountOnBase(untaxed, sumTaxPercents(taxesWithholding(taxes)))
	netTotal = decSub(decSum(untaxed, leviedAmt), withholdingAmt).NormalizeDecimals()
	return
}

type invoiceLinesTotals struct {
	UntaxedSubtotal  fields.DecimalSix
	LinesLevied      fields.DecimalSix
	LinesWithholding fields.DecimalSix
}

func (t invoiceLinesTotals) linesGrossBeforeWithholding() fields.DecimalSix {
	return decSum(t.UntaxedSubtotal, t.LinesLevied).NormalizeDecimals()
}

func headerTaxSplit(untaxedSubtotal fields.DecimalSix, headerTaxes []finance_taxes.Tax, lineTaxIDs map[uint]struct{}) (levied, withholding fields.DecimalSix) {
	for _, tax := range documentLevelHeaderTaxes(headerTaxes, lineTaxIDs) {
		amt := taxAmountForTax(untaxedSubtotal, tax)
		if effectiveTaxKind(tax) == finance_taxes.TaxKindWithholding {
			withholding = decSum(withholding, amt)
		} else {
			levied = decSum(levied, amt)
		}
	}
	return levied.NormalizeDecimals(), withholding.NormalizeDecimals()
}

// invoiceReceivableGrandTotal is the net amount due (AR) after levied and withholding taxes.
func invoiceReceivableGrandTotal(totals invoiceLinesTotals, headerTaxes []finance_taxes.Tax, lineTaxIDs map[uint]struct{}) fields.DecimalSix {
	headerLevied, headerWithholding := headerTaxSplit(totals.UntaxedSubtotal, headerTaxes, lineTaxIDs)
	gross := decSum(totals.linesGrossBeforeWithholding(), headerLevied)
	withheld := decSum(totals.LinesWithholding, headerWithholding)
	return decSub(gross, withheld).NormalizeDecimals()
}

func withholdingTaxAccountID(t finance_taxes.Tax) (uint, error) {
	if t.AccountID == nil || *t.AccountID == 0 {
		name := strings.TrimSpace(t.Name)
		if name == "" {
			name = fmt.Sprintf("#%d", t.ID)
		}
		return 0, fmt.Errorf("withholding tax %q requires a ledger account", name)
	}
	return *t.AccountID, nil
}

func validateWithholdingTaxAccounts(taxes []finance_taxes.Tax) error {
	for _, t := range taxes {
		if effectiveTaxKind(t) == finance_taxes.TaxKindWithholding {
			if _, err := withholdingTaxAccountID(t); err != nil {
				return err
			}
		}
	}
	return nil
}

// paymentWithholdingTotal returns total withholding tax on a payment settlement amount.
func paymentWithholdingTotal(settlement fields.DecimalSix, taxes []finance_taxes.Tax) fields.DecimalSix {
	return taxAmountOnBase(settlement, sumTaxPercents(taxesWithholding(taxes))).NormalizeDecimals()
}

// paymentBankAmount is cash received (settlement minus payment-time withholding).
func paymentBankAmount(settlement fields.DecimalSix, taxes []finance_taxes.Tax) fields.DecimalSix {
	return decSub(settlement, paymentWithholdingTotal(settlement, taxes)).NormalizeDecimals()
}

func validatePaymentTaxes(taxes []finance_taxes.Tax) error {
	if err := validateWithholdingTaxAccounts(taxes); err != nil {
		return err
	}
	for _, t := range taxes {
		if effectiveTaxKind(t) == finance_taxes.TaxKindLevied {
			name := strings.TrimSpace(t.Name)
			if name == "" {
				name = fmt.Sprintf("#%d", t.ID)
			}
			return fmt.Errorf("levied tax %q cannot be applied on a payment; use withholding taxes only", name)
		}
	}
	return nil
}

func collectTaxesFromLines(lines []DraftInvoiceLine) []finance_taxes.Tax {
	var out []finance_taxes.Tax
	for _, ln := range lines {
		out = append(out, ln.Taxes...)
	}
	return out
}

func decimalSixDisplayWithholding(d fields.DecimalSix) string {
	if decimalIsZero(d) {
		return "—"
	}
	return "(" + decimalSixDisplay(d) + ")"
}
