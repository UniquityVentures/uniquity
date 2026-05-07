package p_uniquity_employees

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/UniquityVentures/lago/lago"
	"github.com/UniquityVentures/lago/plugins/p_users"
	"gorm.io/gorm"
)

// Employee is a 1:1 extension of User (at most one employee row per user).
type Employee struct {
	gorm.Model

	UserID uint        `gorm:"unique;not null"`
	User   p_users.User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// PointsTransaction is an append-only ledger row (updates blocked in GORM and UI).
type PointsTransaction struct {
	gorm.Model

	Points PointsDecimal `gorm:"type:numeric(19,2);not null"`

	FromUserID uint         `gorm:"not null"`
	FromUser   p_users.User `gorm:"foreignKey:FromUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	ToEmployeeID uint     `gorm:"not null"`
	ToEmployee   Employee `gorm:"foreignKey:ToEmployeeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (p *PointsTransaction) BeforeCreate(tx *gorm.DB) error {
	if p.Points.R != nil {
		p.Points.R = roundRatTo2Decimals(p.Points.R)
	} else {
		p.Points.R = big.NewRat(0, 1)
	}
	var from p_users.User
	if err := tx.First(&from, p.FromUserID).Error; err != nil {
		return fmt.Errorf("from user: %w", err)
	}
	if !from.IsSuperuser {
		return errors.New("from user must be a superuser")
	}
	return nil
}

func (*PointsTransaction) BeforeUpdate(_ *gorm.DB) error {
	return errors.New("points transactions cannot be updated")
}

func init() {
	lago.OnDBInit("p_uniquity_employees.models", func(d *gorm.DB) *gorm.DB {
		lago.RegisterModel[Employee](d)
		lago.RegisterModel[PointsTransaction](d)
		installPointsTransactionSuperuserTrigger(d)
		return d
	})

	lago.RegistryAdmin.Register("p_uniquity_employees_staff", lago.AdminPanel[Employee]{
		SearchField: "User.Name",
		ListFields:  []string{"User.Name", "User.Email", "UpdatedAt"},
		Preload:     []string{"User"},
	})

	lago.RegistryAdmin.Register("p_uniquity_employees_points", lago.AdminPanel[PointsTransaction]{
		SearchField: "FromUser.Name",
		ListFields: []string{
			"Points",
			"FromUser.Name",
			"ToEmployee.User.Name",
			"CreatedAt",
		},
		Preload: []string{"FromUser", "ToEmployee", "ToEmployee.User"},
	})
}
