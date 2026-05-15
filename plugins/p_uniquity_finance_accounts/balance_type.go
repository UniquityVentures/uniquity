package p_uniquity_finance_accounts

import (
	"database/sql/driver"
	"fmt"
)

// BalanceType mirrors the PostgreSQL enum "BalanceType" (labels Credit, Debit).
type BalanceType string

const (
	BalanceTypeCredit BalanceType = "Credit"
	BalanceTypeDebit  BalanceType = "Debit"
)

func (b BalanceType) Value() (driver.Value, error) {
	switch b {
	case BalanceTypeCredit, BalanceTypeDebit:
		return string(b), nil
	default:
		return nil, fmt.Errorf("invalid BalanceType: %q", b)
	}
}

func (b *BalanceType) Scan(src any) error {
	if src == nil {
		return fmt.Errorf("BalanceType: NULL")
	}
	var s string
	switch v := src.(type) {
	case []byte:
		s = string(v)
	case string:
		s = v
	default:
		return fmt.Errorf("BalanceType: cannot scan %T", src)
	}
	switch BalanceType(s) {
	case BalanceTypeCredit, BalanceTypeDebit:
		*b = BalanceType(s)
		return nil
	default:
		return fmt.Errorf("BalanceType: unknown value %q", s)
	}
}
