package p_uniquity_accounting

import (
	"github.com/UniquityVentures/lamu/registry"
	currencies "github.com/UniquityVentures/uniquity/plugins/p_uniquity_currencies"
	ent "github.com/UniquityVentures/uniquity/plugins/p_uniquity_entities"
	"gorm.io/gorm"
)

// Chart of account type values (Django-style groups).
const (
	AccountTypeAssetReceivable     = "asset_receivable"
	AccountTypeAssetCash           = "asset_cash"
	AccountTypeAssetCurrent        = "asset_current"
	AccountTypeAssetNonCurrent     = "asset_non_current"
	AccountTypeAssetPrepayments    = "asset_prepayments"
	AccountTypeAssetFixed          = "asset_fixed"
	AccountTypeLiabilityPayable    = "liability_payable"
	AccountTypeLiabilityCreditCard = "liability_credit_card"
	AccountTypeLiabilityCurrent    = "liability_current"
	AccountTypeLiabilityNonCurrent = "liability_non_current"
	AccountTypeEquity              = "equity"
	AccountTypeEquityRetained      = "equity_retained"
	AccountTypeIncome              = "income"
	AccountTypeIncomeOther         = "income_other"
	AccountTypeExpense             = "expense"
	AccountTypeExpenseDepreciation = "expense_depreciation"
	AccountTypeExpenseDirectCost   = "expense_direct_cost"
	AccountTypeOffBalance          = "off_balance"
)

// CodeRangeByAccountType suggests numeric code ranges by account type (UI / validation hints; not stored in DB).
var CodeRangeByAccountType = map[string][2]int{
	AccountTypeAssetCash:           {1000, 1099},
	AccountTypeAssetReceivable:     {1100, 1199},
	AccountTypeAssetCurrent:        {1200, 1299},
	AccountTypeAssetPrepayments:    {1300, 1399},
	AccountTypeAssetNonCurrent:     {1500, 1599},
	AccountTypeAssetFixed:          {1600, 1699},
	AccountTypeLiabilityPayable:    {2100, 2199},
	AccountTypeLiabilityCreditCard: {2200, 2299},
	AccountTypeLiabilityCurrent:    {2300, 2399},
	AccountTypeLiabilityNonCurrent: {2500, 2599},
	AccountTypeEquity:              {3000, 3099},
	AccountTypeEquityRetained:      {3100, 3199},
	AccountTypeIncome:              {4000, 4199},
	AccountTypeIncomeOther:         {4200, 4299},
	AccountTypeExpense:             {5000, 5099},
	AccountTypeExpenseDepreciation: {5100, 5199},
	AccountTypeExpenseDirectCost:   {5200, 5299},
	AccountTypeOffBalance:          {9000, 9999},
}

// AccountTypeChoices returns select options (key + label) for HTML forms.
func AccountTypeChoices() []registry.Pair[string, string] {
	return []registry.Pair[string, string]{
		{Key: AccountTypeAssetReceivable, Value: "Receivable"},
		{Key: AccountTypeAssetCash, Value: "Bank and Cash"},
		{Key: AccountTypeAssetCurrent, Value: "Current Assets"},
		{Key: AccountTypeAssetNonCurrent, Value: "Non-current Assets"},
		{Key: AccountTypeAssetPrepayments, Value: "Prepayments"},
		{Key: AccountTypeAssetFixed, Value: "Fixed Assets"},
		{Key: AccountTypeLiabilityPayable, Value: "Payable"},
		{Key: AccountTypeLiabilityCreditCard, Value: "Credit Card"},
		{Key: AccountTypeLiabilityCurrent, Value: "Current Liabilities"},
		{Key: AccountTypeLiabilityNonCurrent, Value: "Non-current Liabilities"},
		{Key: AccountTypeEquity, Value: "Equity"},
		{Key: AccountTypeEquityRetained, Value: "Current Year Earnings"},
		{Key: AccountTypeIncome, Value: "Income"},
		{Key: AccountTypeIncomeOther, Value: "Other Income"},
		{Key: AccountTypeExpense, Value: "Expenses"},
		{Key: AccountTypeExpenseDepreciation, Value: "Depreciation"},
		{Key: AccountTypeExpenseDirectCost, Value: "Cost of Revenue"},
		{Key: AccountTypeOffBalance, Value: "Off-Balance Sheet"},
	}
}

// Account is a chart-of-accounts entry scoped to a legal entity.
type Account struct {
	gorm.Model

	EntityID uint       `gorm:"not null;index"`
	Entity   ent.Entity `gorm:"constraint:OnDelete:CASCADE"`

	Code        string `gorm:"size:20"`
	Name        string `gorm:"size:100;not null"`
	AccountType string `gorm:"size:30;not null;index"`

	CurrencyID *uint                `gorm:"index"`
	Currency   *currencies.Currency `gorm:"constraint:OnDelete:SET NULL"`

	IsActive       bool `gorm:"not null;default:true"`
	IsReconcilable bool `gorm:"not null;default:false"`
}
