package p_uniquity_finance_accounts

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/getters"
	"github.com/lariv-in/lago/views"
	"gorm.io/gorm"
)

// balanceTypeScopeQueryParam scopes account pickers to one balance type (not a model field).
const balanceTypeScopeQueryParam = "balance_type_scope"

// AccountSelectRouteURL returns the account picker URL filtered to the given balance type.
func AccountSelectRouteURL(balanceType BalanceType) getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		base, err := lago.RoutePath("finance_accounts.AccountSelectRoute", nil)(ctx)
		if err != nil {
			return "", err
		}
		u, err := url.Parse(base)
		if err != nil {
			return "", fmt.Errorf("parse account select URL: %w", err)
		}
		q := u.Query()
		q.Set(balanceTypeScopeQueryParam, string(balanceType))
		u.RawQuery = q.Encode()
		return u.String(), nil
	}
}

// accountSelectBalanceTypeScope enforces balance_type_scope from the query string (survives filter forms via hidden field).
type accountSelectBalanceTypeScope struct{}

func (accountSelectBalanceTypeScope) Patch(_ views.View, r *http.Request, query gorm.ChainInterface[Account]) gorm.ChainInterface[Account] {
	bt := r.URL.Query().Get(balanceTypeScopeQueryParam)
	if bt == "" {
		return query
	}
	return query.Where("balance_type = ?", bt)
}

// ValidateLeafAccountBalanceType ensures accountID is a non-group account with the expected balance type.
func ValidateLeafAccountBalanceType(tx *gorm.DB, accountID uint, want BalanceType, label string) error {
	if accountID == 0 {
		return fmt.Errorf("%s is required", label)
	}
	var acct Account
	if err := tx.Select("id", "balance_type", "is_group").First(&acct, accountID).Error; err != nil {
		return fmt.Errorf("%s: %w", label, err)
	}
	if acct.IsGroup {
		return fmt.Errorf("%s: group accounts cannot be used for posting", label)
	}
	if acct.BalanceType != want {
		return fmt.Errorf("%s: account must have balance type %s", label, want)
	}
	return nil
}
