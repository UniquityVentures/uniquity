package p_uniquity_invoices

import (
	"embed"

	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

//go:embed migrations
var migrationsFS embed.FS

func pluginMigrations() lamu.PluginFeatures[lamu.UsefulFilesystem] {
	return lamu.PluginFeatures[lamu.UsefulFilesystem]{
		Entries: []registry.Pair[string, lamu.UsefulFilesystem]{
			{Key: "p_uniquity_invoices.migrations", Value: migrationsFS},
		},
	}
}
