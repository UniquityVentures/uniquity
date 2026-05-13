package p_uniquity_entities

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pluginRoutes() lamu.PluginFeatures[lamu.Route] {
	base := AppUrl + "e/"
	return lamu.PluginFeatures[lamu.Route]{
		Entries: []registry.Pair[string, lamu.Route]{
			{Key: "entities.EntityListRoute", Value: lamu.Route{Path: AppUrl, Handler: lamu.NewDynamicView("entities.EntityListView")}},
			{Key: "entities.EntitySelectRoute", Value: lamu.Route{Path: AppUrl + "select/", Handler: lamu.NewDynamicView("entities.EntitySelectView")}},
			{Key: "entities.EntityCreateRoute", Value: lamu.Route{Path: AppUrl + "create/", Handler: lamu.NewDynamicView("entities.EntityCreateView")}},
			{Key: "entities.EntityDetailRoute", Value: lamu.Route{Path: base + "{id}/", Handler: lamu.NewDynamicView("entities.EntityDetailView")}},
			{Key: "entities.EntityUpdateRoute", Value: lamu.Route{Path: base + "{id}/edit/", Handler: lamu.NewDynamicView("entities.EntityUpdateView")}},
			{Key: "entities.EntityDeleteRoute", Value: lamu.Route{Path: base + "{id}/delete/", Handler: lamu.NewDynamicView("entities.EntityDeleteView")}},
		},
	}
}
