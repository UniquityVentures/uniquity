module github.com/UniquityVentures/uniquity/plugins/p_uniquity_products

go 1.26.1

require (
	github.com/UniquityVentures/lamu v0.4.7
	github.com/UniquityVentures/lamu/plugins/p_users v0.4.7
	github.com/UniquityVentures/uniquity/plugins/p_uniquity_entities v0.0.0
	gorm.io/gorm v1.31.1
	maragu.dev/gomponents v1.3.0
)

replace github.com/UniquityVentures/uniquity/plugins/p_uniquity_entities => ../p_uniquity_entities
