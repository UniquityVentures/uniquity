package p_uniquity_finance_products

import (
	"context"
	"fmt"
	"strings"

	"github.com/UniquityVentures/lamu/fields"
	"github.com/UniquityVentures/lamu/plugins/p_llm_assistant"
	finance_taxes "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_taxes"
	"gorm.io/gorm"
)

type productResult struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func init() {
	p_llm_assistant.ExprEnvRegistry.Register("search_products", p_llm_assistant.ContextualFunc(func(ctx context.Context, db *gorm.DB) any {
		return func(query string) ([]productResult, error) {
			var products []Product
			q := db.WithContext(ctx).Model(&Product{})

			query = strings.TrimSpace(query)
			if query != "" {
				q = q.Where("name LIKE ?", "%"+query+"%")
			}

			if err := q.Order("name ASC").Find(&products).Error; err != nil {
				return nil, err
			}

			results := make([]productResult, 0, len(products))
			for _, p := range products {
				results = append(results, productResult{
					ID:   p.ID,
					Name: p.Name,
				})
			}
			return results, nil
		}
	}))

	p_llm_assistant.ExprEnvRegistry.Register("create_product", p_llm_assistant.ContextualFunc(func(ctx context.Context, db *gorm.DB) any {
		return func(name string, reference string, detailsVal any) (uint, error) {
			name = strings.TrimSpace(name)
			if name == "" {
				return 0, fmt.Errorf("product name is required")
			}
			reference = strings.TrimSpace(reference)
			if reference == "" {
				return 0, fmt.Errorf("product reference code is required")
			}

			var details map[string]any
			if detailsVal != nil {
				if m, ok := detailsVal.(map[string]any); ok {
					details = m
				} else if ma, ok := detailsVal.(map[any]any); ok {
					details = make(map[string]any)
					for k, v := range ma {
						details[fmt.Sprint(k)] = v
					}
				}
			}

			getStr := func(keys ...string) string {
				for _, k := range keys {
					if v, ok := details[k]; ok && v != nil {
						return strings.TrimSpace(fmt.Sprint(v))
					}
				}
				return ""
			}

			getInt64 := func(keys ...string) int64 {
				for _, k := range keys {
					if v, ok := details[k]; ok && v != nil {
						switch val := v.(type) {
						case int:
							return int64(val)
						case int64:
							return val
						case float64:
							return int64(val)
						}
					}
				}
				return 0
			}

			pType := ProductTypeGoods
			if typeStr := getStr("type", "Type"); typeStr != "" {
				switch ProductType(typeStr) {
				case ProductTypeGoods, ProductTypeServices, ProductTypeBoth:
					pType = ProductType(typeStr)
				default:
					return 0, fmt.Errorf("invalid product type: %q", typeStr)
				}
			}

			var salesPrice fields.DecimalSix
			if spVal, ok := details["sales_price"]; ok && spVal != nil {
				_ = salesPrice.UnmarshalText([]byte(fmt.Sprint(spVal)))
			} else if spVal, ok := details["SalesPrice"]; ok && spVal != nil {
				_ = salesPrice.UnmarshalText([]byte(fmt.Sprint(spVal)))
			}

			var baseCost fields.DecimalSix
			if bcVal, ok := details["base_cost"]; ok && bcVal != nil {
				_ = baseCost.UnmarshalText([]byte(fmt.Sprint(bcVal)))
			} else if bcVal, ok := details["BaseCost"]; ok && bcVal != nil {
				_ = baseCost.UnmarshalText([]byte(fmt.Sprint(bcVal)))
			}

			product := &Product{
				Name:       name,
				Reference:  reference,
				Type:       pType,
				Remarks:    getStr("remarks", "Remarks"),
				BaseCost:   baseCost,
				SalesPrice: salesPrice,
				HSNCode:    getInt64("hsn_code", "hsn", "HSNCode"),
			}

			var taxIDs []uint
			if tIDsVal, ok := details["tax_ids"]; ok && tIDsVal != nil {
				toUint := func(v any) (uint, error) {
					switch val := v.(type) {
					case int:
						return uint(val), nil
					case int64:
						return uint(val), nil
					case float64:
						return uint(val), nil
					case uint:
						return val, nil
					default:
						return 0, fmt.Errorf("invalid tax ID: %v", v)
					}
				}
				if slice, ok := tIDsVal.([]any); ok {
					for _, sVal := range slice {
						tid, err := toUint(sVal)
						if err != nil {
							return 0, err
						}
						taxIDs = append(taxIDs, tid)
					}
				} else if slice, ok := tIDsVal.([]uint); ok {
					taxIDs = slice
				} else if slice, ok := tIDsVal.([]int); ok {
					for _, i := range slice {
						taxIDs = append(taxIDs, uint(i))
					}
				}
			}

			err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
				if err := tx.Create(product).Error; err != nil {
					return err
				}
				if len(taxIDs) > 0 {
					var taxes []finance_taxes.Tax
					if err := tx.Where("id IN ?", taxIDs).Find(&taxes).Error; err != nil {
						return err
					}
					if err := tx.Model(product).Association("Taxes").Append(taxes); err != nil {
						return err
					}
				}
				return nil
			})
			if err != nil {
				return 0, err
			}
			return product.ID, nil
		}
	}))
}
