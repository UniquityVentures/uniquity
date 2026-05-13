package p_uniquity_accounting

import "gorm.io/gorm"

type Account struct {
	gorm.Model

	Code    int    `gorm:"autoIncrement:false"`
	Name    string `gorm:"size:100;not null"`
	IsAsset bool
}
