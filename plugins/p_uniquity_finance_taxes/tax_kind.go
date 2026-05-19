package p_uniquity_finance_taxes

import (
	"database/sql/driver"
	"fmt"
)

// TaxKind mirrors the PostgreSQL enum "TaxKind".
type TaxKind string

const (
	TaxKindLevied      TaxKind = "levied"
	TaxKindWithholding TaxKind = "withholding"
)

func (k TaxKind) Value() (driver.Value, error) {
	switch k {
	case TaxKindLevied, TaxKindWithholding:
		return string(k), nil
	default:
		return nil, fmt.Errorf("invalid TaxKind: %q", k)
	}
}

func (k *TaxKind) Scan(src any) error {
	if src == nil {
		return fmt.Errorf("TaxKind: NULL")
	}
	var s string
	switch v := src.(type) {
	case []byte:
		s = string(v)
	case string:
		s = v
	default:
		return fmt.Errorf("TaxKind: cannot scan %T", src)
	}
	switch TaxKind(s) {
	case TaxKindLevied, TaxKindWithholding:
		*k = TaxKind(s)
		return nil
	default:
		return fmt.Errorf("TaxKind: unknown value %q", s)
	}
}
