package p_uniquity_finance_products

import (
	"log/slog"

	finance_accounts "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_accounts"
	"github.com/lariv-in/lago"
	"gorm.io/gorm"
)

// Product preference form field names shared with the accounting preferences page patch.
const (
	ProductPrefInventoryAccountIDField = "InventoryAccountID"
	ProductPrefCostOfSalesAcctIDField  = "CostOfSalesAcctID"
)

// ProductPreferences is the singleton row for product-wide GL account settings (one row, typically id = 1).
type ProductPreferences struct {
	gorm.Model

	InventoryAccountID *uint                     `gorm:"column:inventory_account_id"`
	InventoryAccount   *finance_accounts.Account `gorm:"foreignKey:InventoryAccountID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CostOfSalesAcctID  *uint                     `gorm:"column:cost_of_sales_account_id"`
	CostOfSalesAccount *finance_accounts.Account `gorm:"foreignKey:CostOfSalesAcctID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// LoadProductPreferences returns the singleton preferences row, creating id 1 if missing.
func LoadProductPreferences(db *gorm.DB) ProductPreferences {
	var prefs ProductPreferences
	if err := db.FirstOrCreate(&prefs, ProductPreferences{Model: gorm.Model{ID: 1}}).Error; err != nil {
		slog.Warn("LoadProductPreferences", "error", err)
	}
	return prefs
}

func init() {
	lago.RegistryAdmin.Register("p_uniquity_finance_products.ProductPreferences", lago.AdminPanel[ProductPreferences]{
		SearchField: "",
	})
}
