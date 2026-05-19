package p_uniquity_finance_products

import (
	"database/sql/driver"
	"fmt"
)

// ProductType mirrors the PostgreSQL enum "ProductType".
type ProductType string

const (
	ProductTypeGoods    ProductType = "Goods"
	ProductTypeServices ProductType = "Services"
	ProductTypeBoth     ProductType = "Both"
)

func (p ProductType) Value() (driver.Value, error) {
	switch p {
	case ProductTypeGoods, ProductTypeServices, ProductTypeBoth:
		return string(p), nil
	default:
		return nil, fmt.Errorf("invalid ProductType: %q", p)
	}
}

func (p *ProductType) Scan(src any) error {
	if src == nil {
		return fmt.Errorf("ProductType: NULL")
	}
	var s string
	switch v := src.(type) {
	case []byte:
		s = string(v)
	case string:
		s = v
	default:
		return fmt.Errorf("ProductType: cannot scan %T", src)
	}
	switch ProductType(s) {
	case ProductTypeGoods, ProductTypeServices, ProductTypeBoth:
		*p = ProductType(s)
		return nil
	default:
		return fmt.Errorf("ProductType: unknown value %q", s)
	}
}
