package p_uniquity_finance_taxes

import (
	"context"
	"strings"

	"github.com/UniquityVentures/lamu/plugins/p_llm_assistant"
	"gorm.io/gorm"
)

type taxResult struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func init() {
	p_llm_assistant.ExprEnvRegistry.Register("search_taxes", p_llm_assistant.ContextualFunc(func(ctx context.Context, db *gorm.DB) any {
		return func(query string) ([]taxResult, error) {
			var taxes []Tax
			q := db.WithContext(ctx).Model(&Tax{})

			query = strings.TrimSpace(query)
			if query != "" {
				q = q.Where("name LIKE ?", "%"+query+"%")
			}

			if err := q.Order("name ASC").Find(&taxes).Error; err != nil {
				return nil, err
			}

			results := make([]taxResult, 0, len(taxes))
			for _, t := range taxes {
				results = append(results, taxResult{
					ID:   t.ID,
					Name: t.Name,
				})
			}
			return results, nil
		}
	}))
}
