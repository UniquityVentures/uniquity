package p_uniquity_finance_accounts

import (
	"database/sql/driver"
	"fmt"
)

// JournalType mirrors the PostgreSQL enum "JournalType".
type JournalType string

const (
	JournalTypeGeneral JournalType = "General"
)

func (t JournalType) Value() (driver.Value, error) {
	switch t {
	case JournalTypeGeneral:
		return string(t), nil
	default:
		return nil, fmt.Errorf("invalid JournalType: %q", t)
	}
}

func (t *JournalType) Scan(src any) error {
	if src == nil {
		return fmt.Errorf("JournalType: NULL")
	}
	var s string
	switch v := src.(type) {
	case []byte:
		s = string(v)
	case string:
		s = v
	default:
		return fmt.Errorf("JournalType: cannot scan %T", src)
	}
	switch JournalType(s) {
	case JournalTypeGeneral:
		*t = JournalType(s)
		return nil
	default:
		return fmt.Errorf("JournalType: unknown value %q", s)
	}
}
