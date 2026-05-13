package p_uniquity_accounting

import (
	"github.com/UniquityVentures/lamu/fields"
	"gorm.io/gorm"
)

type Posting struct {
	gorm.Model

	Amount fields.DecimalSix `gorm:"type:numeric(19,6);not null"`

	AccountID uint `gorm:"not null"`
	Account   Account
}
