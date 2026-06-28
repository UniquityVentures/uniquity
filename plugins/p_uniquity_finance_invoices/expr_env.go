package p_uniquity_finance_invoices

import (
	"context"
	"fmt"
	"time"

	"github.com/UniquityVentures/lamu/plugins/p_llm_assistant"
	"gorm.io/gorm"
)

func init() {
	p_llm_assistant.ExprEnvRegistry.Register("create_draft_invoice", p_llm_assistant.ContextualFunc(func(ctx context.Context, db *gorm.DB) any {
		return func(customerIDVal any, paymentTermIDVal any, dateStr string, linesVal []any) (uint, error) {
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
					return 0, fmt.Errorf("invalid numeric value: %v", v)
				}
			}

			customerID, err := toUint(customerIDVal)
			if err != nil {
				return 0, fmt.Errorf("customer_id: %w", err)
			}

			paymentTermID, err := toUint(paymentTermIDVal)
			if err != nil {
				return 0, fmt.Errorf("payment_term_id: %w", err)
			}

			var dt time.Time
			if dateStr == "" {
				dt = time.Now()
			} else {
				var err error
				dt, err = time.Parse("2006-01-02", dateStr)
				if err != nil {
					dt, err = time.Parse(time.RFC3339, dateStr)
					if err != nil {
						return 0, fmt.Errorf("invalid date format: %w", err)
					}
				}
			}

			var pendingLines []DraftLinePending
			for idx, item := range linesVal {
				m, ok := item.(map[string]any)
				if !ok {
					if ma, ok := item.(map[any]any); ok {
						m = make(map[string]any)
						for k, v := range ma {
							m[fmt.Sprint(k)] = v
						}
					} else {
						return 0, fmt.Errorf("line %d is not a map", idx+1)
					}
				}

				getVal := func(keys ...string) any {
					for _, k := range keys {
						if v, ok := m[k]; ok {
							return v
						}
					}
					return nil
				}

				prodIDVal := getVal("product_id", "productId", "ProductID")
				if prodIDVal == nil {
					return 0, fmt.Errorf("line %d: missing product_id", idx+1)
				}
				prodID, err := toUint(prodIDVal)
				if err != nil {
					return 0, fmt.Errorf("line %d product_id: %w", idx+1, err)
				}

				qtyVal := getVal("quantity", "Quantity")
				if qtyVal == nil {
					return 0, fmt.Errorf("line %d: missing quantity", idx+1)
				}

				rateVal := getVal("rate", "Rate")
				var rateStr string
				if rateVal != nil {
					rateStr = fmt.Sprint(rateVal)
				}

				var taxIDs []uint
				if tIDsVal := getVal("tax_ids", "taxIds", "TaxIDs"); tIDsVal != nil {
					if slice, ok := tIDsVal.([]any); ok {
						for _, sVal := range slice {
							tid, err := toUint(sVal)
							if err != nil {
								return 0, fmt.Errorf("line %d tax_id: %w", idx+1, err)
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

				pendingLines = append(pendingLines, DraftLinePending{
					ProductID: prodID,
					Quantity:  fmt.Sprint(qtyVal),
					Rate:      rateStr,
					TaxIDs:    taxIDs,
				})
			}

			draft := &DraftInvoice{
				CustomerID:    customerID,
				PaymentTermID: paymentTermID,
				Datetime:      dt,
				PendingLines:  pendingLines,
			}

			if err := db.WithContext(ctx).Create(draft).Error; err != nil {
				return 0, err
			}

			return draft.ID, nil
		}
	}))

	p_llm_assistant.ExprEnvRegistry.Register("list_payment_terms", p_llm_assistant.ContextualFunc(func(ctx context.Context, db *gorm.DB) any {
		return func() ([]map[string]any, error) {
			var terms []PaymentTerm
			if err := db.WithContext(ctx).Find(&terms).Error; err != nil {
				return nil, err
			}
			var results []map[string]any
			for _, pt := range terms {
				inst, err := ResolvePaymentTermInstanceFromTerm(ctx, &pt)
				summary := "Unknown"
				if err == nil {
					summary = inst.Summary()
				}
				results = append(results, map[string]any{
					"id":      pt.ID,
					"type":    pt.Type,
					"summary": summary,
				})
			}
			return results, nil
		}
	}))

	p_llm_assistant.ExprEnvRegistry.Register("add_relative_payment_term", p_llm_assistant.ContextualFunc(func(ctx context.Context, db *gorm.DB) any {
		return func(durationVal any) (uint, error) {
			var dur time.Duration
			switch v := durationVal.(type) {
			case time.Duration:
				dur = v
			case string:
				var err error
				dur, err = time.ParseDuration(v)
				if err != nil {
					return 0, fmt.Errorf("invalid duration format: %w", err)
				}
			case int:
				dur = time.Duration(v)
			case int64:
				dur = time.Duration(v)
			case float64:
				dur = time.Duration(v)
			default:
				return 0, fmt.Errorf("invalid duration kind: %T", durationVal)
			}
			if dur <= 0 {
				return 0, fmt.Errorf("duration must be positive")
			}

			var pt uint
			err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
				backing := PaymentTermRelative{Duration: dur}
				if err := tx.Create(&backing).Error; err != nil {
					return err
				}
				umbrella := PaymentTerm{
					Type:      PaymentTermTypeRelative,
					BackingID: backing.ID,
				}
				if err := tx.Create(&umbrella).Error; err != nil {
					return err
				}
				pt = umbrella.ID
				return nil
			})
			if err != nil {
				return 0, err
			}
			return pt, nil
		}
	}))

	p_llm_assistant.ExprEnvRegistry.Register("add_due_date_payment_term", p_llm_assistant.ContextualFunc(func(ctx context.Context, db *gorm.DB) any {
		return func(dateVal any) (uint, error) {
			var dt time.Time
			switch v := dateVal.(type) {
			case time.Time:
				dt = v
			case string:
				var err error
				dt, err = time.Parse("2006-01-02", v)
				if err != nil {
					dt, err = time.Parse(time.RFC3339, v)
					if err != nil {
						return 0, fmt.Errorf("invalid date format: %w", err)
					}
				}
			default:
				return 0, fmt.Errorf("invalid date kind: %T", dateVal)
			}
			if dt.IsZero() {
				return 0, fmt.Errorf("datetime is required")
			}

			var pt uint
			err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
				backing := PaymentTermDueDate{Datetime: dt}
				if err := tx.Create(&backing).Error; err != nil {
					return err
				}
				umbrella := PaymentTerm{
					Type:      PaymentTermTypeDueDate,
					BackingID: backing.ID,
				}
				if err := tx.Create(&umbrella).Error; err != nil {
					return err
				}
				pt = umbrella.ID
				return nil
			})
			if err != nil {
				return 0, err
			}
			return pt, nil
		}
	}))

	p_llm_assistant.ExprEnvRegistry.Register("find_due_date_payment_term", p_llm_assistant.ContextualFunc(func(ctx context.Context, db *gorm.DB) any {
		return func(dateVal any) (uint, error) {
			var dt time.Time
			switch v := dateVal.(type) {
			case time.Time:
				dt = v
			case string:
				var err error
				dt, err = time.Parse("2006-01-02", v)
				if err != nil {
					dt, err = time.Parse(time.RFC3339, v)
					if err != nil {
						return 0, fmt.Errorf("invalid date format: %w", err)
					}
				}
			default:
				return 0, fmt.Errorf("invalid date kind: %T", dateVal)
			}
			if dt.IsZero() {
				return 0, fmt.Errorf("datetime is required")
			}

			var backing PaymentTermDueDate
			err := db.WithContext(ctx).Where("datetime = ?", dt).First(&backing).Error
			if err == nil {
				var umbrella PaymentTerm
				err = db.WithContext(ctx).Where("type = ? AND backing_id = ?", PaymentTermTypeDueDate, backing.ID).First(&umbrella).Error
				if err == nil {
					return umbrella.ID, nil
				}
			}

			var pt uint
			err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
				backing = PaymentTermDueDate{Datetime: dt}
				if err := tx.Create(&backing).Error; err != nil {
					return err
				}
				umbrella := PaymentTerm{
					Type:      PaymentTermTypeDueDate,
					BackingID: backing.ID,
				}
				if err := tx.Create(&umbrella).Error; err != nil {
					return err
				}
				pt = umbrella.ID
				return nil
			})
			if err != nil {
				return 0, err
			}
			return pt, nil
		}
	}))
}
