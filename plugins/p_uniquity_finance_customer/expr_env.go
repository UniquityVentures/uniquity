package p_uniquity_finance_customer

import (
	"context"
	"fmt"
	"strings"

	"github.com/UniquityVentures/lamu/plugins/p_llm_assistant"
	"gorm.io/gorm"
)

type customerResult struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func init() {
	p_llm_assistant.ExprEnvRegistry.Register("search_customers", p_llm_assistant.ContextualFunc(func(ctx context.Context, db *gorm.DB) any {
		return func(query string) ([]customerResult, error) {
			var customers []Customer
			q := db.WithContext(ctx).Model(&Customer{})

			query = strings.TrimSpace(query)
			if query != "" {
				q = q.Where("name LIKE ?", "%"+query+"%")
			}

			if err := q.Order("name ASC").Find(&customers).Error; err != nil {
				return nil, err
			}

			results := make([]customerResult, 0, len(customers))
			for _, c := range customers {
				results = append(results, customerResult{
					ID:   c.ID,
					Name: c.Name,
				})
			}
			return results, nil
		}
	}))

	p_llm_assistant.ExprEnvRegistry.Register("create_customer", p_llm_assistant.ContextualFunc(func(ctx context.Context, db *gorm.DB) any {
		return func(name string, detailsVal any) (uint, error) {
			name = strings.TrimSpace(name)
			if name == "" {
				return 0, fmt.Errorf("customer name is required")
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

			customer := &Customer{
				Name:    name,
				Address: getStr("address", "Address"),
				GSTIN:   getStr("gstin", "GSTIN", "Gstin"),
				PAN:     getStr("pan", "PAN", "Pan"),
				Phone:   getStr("phone", "Phone"),
				Email:   getStr("email", "Email"),
				Website: getStr("website", "Website"),
			}

			if err := db.WithContext(ctx).Create(customer).Error; err != nil {
				return 0, err
			}
			return customer.ID, nil
		}
	}))
}
