package p_uniquity_employees

import (
	"database/sql/driver"
	"encoding"
	"fmt"
	"math/big"
	"strings"
)

// PointsDecimal holds a monetary-style amount with exactly two decimal places
// persisted as NUMERIC and represented in Go with *big.Rat.
type PointsDecimal struct {
	R *big.Rat
}

var (
	_ encoding.TextMarshaler   = PointsDecimal{}
	_ encoding.TextUnmarshaler = (*PointsDecimal)(nil)
	_ driver.Valuer            = PointsDecimal{}
)

func (p PointsDecimal) rat() *big.Rat {
	if p.R == nil {
		return big.NewRat(0, 1)
	}
	return p.R
}

func roundRatTo2Decimals(r *big.Rat) *big.Rat {
	if r == nil {
		return big.NewRat(0, 1)
	}
	bf := new(big.Float).SetPrec(256).SetRat(r)
	s := bf.Text('f', 2)
	out := new(big.Rat)
	if _, ok := out.SetString(s); !ok {
		return big.NewRat(0, 1)
	}
	return out
}

// MarshalText implements encoding.TextMarshaler.
func (p PointsDecimal) MarshalText() ([]byte, error) {
	r := roundRatTo2Decimals(p.rat())
	return []byte(r.FloatString(2)), nil
}

// UnmarshalText implements encoding.TextUnmarshaler (mapstructure form binding).
func (p *PointsDecimal) UnmarshalText(text []byte) error {
	s := strings.TrimSpace(string(text))
	if s == "" {
		p.R = big.NewRat(0, 1)
		return nil
	}
	r := new(big.Rat)
	if _, ok := r.SetString(s); !ok {
		return fmt.Errorf("invalid points value %q", s)
	}
	p.R = roundRatTo2Decimals(r)
	return nil
}

// Value implements driver.Valuer for GORM / SQL.
func (p PointsDecimal) Value() (driver.Value, error) {
	r := roundRatTo2Decimals(p.rat())
	return r.FloatString(2), nil
}

// Scan implements sql.Scanner.
func (p *PointsDecimal) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		p.R = big.NewRat(0, 1)
		return nil
	case []byte:
		return p.UnmarshalText(v)
	case string:
		return p.UnmarshalText([]byte(v))
	case int64:
		p.R = big.NewRat(v, 1)
		p.R = roundRatTo2Decimals(p.R)
		return nil
	default:
		return fmt.Errorf("cannot scan %T into PointsDecimal", src)
	}
}

// String returns a fixed two-decimal string for UI.
func (p PointsDecimal) String() string {
	b, err := p.MarshalText()
	if err != nil {
		return "0.00"
	}
	return string(b)
}
