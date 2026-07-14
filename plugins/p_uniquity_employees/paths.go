package p_uniquity_employees

import (
	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/registry"
)

func pluginRoutes() lago.PluginFeatures[lago.Route] {
	emp := AppUrl + "emp/"
	basePts := AppUrl + "points/"
	return lago.PluginFeatures[lago.Route]{
		Entries: []registry.Pair[string, lago.Route]{
			{Key: "employees.DefaultRoute", Value: lago.Route{Path: AppUrl, Handler: lago.NewDynamicView("employees.EmployeeListView")}},
			{Key: "employees.EmployeeCreateRoute", Value: lago.Route{Path: AppUrl + "create/", Handler: lago.NewDynamicView("employees.EmployeeCreateView")}},
			{Key: "employees.EmployeeSelectRoute", Value: lago.Route{Path: AppUrl + "select/", Handler: lago.NewDynamicView("employees.EmployeeSelectView")}},
			{Key: "employees.EmployeeDetailRoute", Value: lago.Route{Path: emp + "{id}/", Handler: lago.NewDynamicView("employees.EmployeeDetailView")}},
			{Key: "employees.EmployeeUpdateRoute", Value: lago.Route{Path: emp + "{id}/edit/", Handler: lago.NewDynamicView("employees.EmployeeUpdateView")}},
			{Key: "employees.EmployeeDeleteRoute", Value: lago.Route{Path: emp + "{id}/delete/", Handler: lago.NewDynamicView("employees.EmployeeDeleteView")}},
			{Key: "employees.PointsListRoute", Value: lago.Route{Path: basePts, Handler: lago.NewDynamicView("employees.PointsListView")}},
			{Key: "employees.PointsCreateRoute", Value: lago.Route{Path: basePts + "create/", Handler: lago.NewDynamicView("employees.PointsCreateView")}},
			{Key: "employees.PointsDetailRoute", Value: lago.Route{Path: basePts + "{id}/", Handler: lago.NewDynamicView("employees.PointsDetailView")}},
		},
	}
}
