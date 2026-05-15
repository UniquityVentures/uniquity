package p_uniquity_finance_invoices

import (
	"database/sql/driver"
	"fmt"
)

// InvoiceStatus mirrors the PostgreSQL enum "InvoiceStatus".
type InvoiceStatus string

const (
	InvoiceStatusDraft     InvoiceStatus = "draft"
	InvoiceStatusPosted    InvoiceStatus = "posted"
	InvoiceStatusCancelled InvoiceStatus = "cancelled"
)

func (s InvoiceStatus) Value() (driver.Value, error) {
	switch s {
	case InvoiceStatusDraft, InvoiceStatusPosted, InvoiceStatusCancelled:
		return string(s), nil
	default:
		return nil, fmt.Errorf("invalid InvoiceStatus: %q", s)
	}
}

func (s *InvoiceStatus) Scan(src any) error {
	if src == nil {
		return fmt.Errorf("InvoiceStatus: NULL")
	}
	var str string
	switch v := src.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return fmt.Errorf("InvoiceStatus: cannot scan %T", src)
	}
	switch InvoiceStatus(str) {
	case InvoiceStatusDraft, InvoiceStatusPosted, InvoiceStatusCancelled:
		*s = InvoiceStatus(str)
		return nil
	default:
		return fmt.Errorf("InvoiceStatus: unknown value %q", str)
	}
}
