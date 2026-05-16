package p_uniquity_finance_fiscal_year

import (
	"time"

	"github.com/UniquityVentures/lamu/lamu"
	"gorm.io/gorm"
)

// FiscalYear is an accounting period with inclusive datetime bounds (similar to semesters).
type FiscalYear struct {
	gorm.Model

	Code     string `gorm:"uniqueIndex"`
	Name     string `gorm:"not null"`
	Start    time.Time `gorm:"column:starts_at;not null"`
	End      time.Time `gorm:"column:ends_at;not null"`
	IsActive bool   `gorm:"not null;default:true"`
}

func init() {
	lamu.RegistryAdmin.Register("p_uniquity_finance_fiscal_year.FiscalYear", lamu.AdminPanel[FiscalYear]{
		SearchField: "Name",
		ListFields:  []string{"Code", "Name", "Start", "End", "IsActive"},
	})
}
