package p_uniquity_finance_accounts

import (
	"fmt"

	"gorm.io/gorm"
)

// BeforeCreate validates parent/child balance_type (mirrors accounts_enforce_parent_balance_type trigger).
func (a *Account) BeforeCreate(tx *gorm.DB) error {
	return a.validateParentBalanceTypeOnSave(tx)
}

// BeforeUpdate validates parent/child balance_type and blocks changing BalanceType when children disagree.
func (a *Account) BeforeUpdate(tx *gorm.DB) error {
	if err := a.validateParentBalanceTypeOnSave(tx); err != nil {
		return err
	}
	if a.ID == 0 {
		return nil
	}
	var old Account
	if err := tx.Select("balance_type").First(&old, a.ID).Error; err != nil {
		return fmt.Errorf("load account for update: %w", err)
	}
	if old.BalanceType == a.BalanceType {
		return nil
	}
	var n int64
	if err := tx.Model(&Account{}).
		Where("parent_id = ? AND balance_type <> ?", a.ID, a.BalanceType).
		Count(&n).Error; err != nil {
		return err
	}
	if n > 0 {
		return fmt.Errorf("cannot change balance_type while child accounts have a different balance_type")
	}
	return nil
}

func (a *Account) validateParentBalanceTypeOnSave(tx *gorm.DB) error {
	if a.ParentID == nil || *a.ParentID == 0 {
		return nil
	}
	var parent Account
	if err := tx.Select("id", "balance_type").First(&parent, *a.ParentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("parent account not found")
		}
		return fmt.Errorf("load parent account: %w", err)
	}
	if parent.BalanceType != a.BalanceType {
		return fmt.Errorf("balance_type must match the parent account balance_type")
	}
	return nil
}
