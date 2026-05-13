module github.com/UniquityVentures/uniquity/plugins/p_uniquity_entities

go 1.26.1

require (
	github.com/UniquityVentures/lamu v0.4.7
	github.com/UniquityVentures/lamu/plugins/p_users v0.4.7
	github.com/UniquityVentures/uniquity/plugins/p_uniquity_currencies v0.0.0
	gorm.io/gorm v1.31.1
)

replace github.com/UniquityVentures/uniquity/plugins/p_uniquity_currencies => ../p_uniquity_currencies
