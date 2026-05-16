module github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_creditnotes

go 1.26.1

require (
	github.com/UniquityVentures/lamu v0.4.8
	github.com/UniquityVentures/lamu/plugins/p_users v0.4.8
	github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_accounts v0.0.0
	gorm.io/gorm v1.31.1
)

replace github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_accounts => ../p_uniquity_finance_accounts
