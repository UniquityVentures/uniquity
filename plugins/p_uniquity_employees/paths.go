package p_uniquity_employees

import (
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pluginRoutes() lamu.PluginFeatures[lamu.Route] {
	emp := AppUrl + "emp/"
	basePts := AppUrl + "points/"
	return lamu.PluginFeatures[lamu.Route]{
		Entries: []registry.Pair[string, lamu.Route]{
			{Key: "employees.DefaultRoute", Value: lamu.Route{Path: AppUrl, Handler: lamu.NewDynamicView("employees.EmployeeListView")}},
			{Key: "employees.EmployeeCreateRoute", Value: lamu.Route{Path: AppUrl + "create/", Handler: lamu.NewDynamicView("employees.EmployeeCreateView")}},
			{Key: "employees.EmployeeSelectRoute", Value: lamu.Route{Path: AppUrl + "select/", Handler: lamu.NewDynamicView("employees.EmployeeSelectView")}},
			{Key: "employees.EmployeeDetailRoute", Value: lamu.Route{Path: emp + "{id}/", Handler: lamu.NewDynamicView("employees.EmployeeDetailView")}},
			{Key: "employees.EmployeeUpdateRoute", Value: lamu.Route{Path: emp + "{id}/edit/", Handler: lamu.NewDynamicView("employees.EmployeeUpdateView")}},
			{Key: "employees.EmployeeDeleteRoute", Value: lamu.Route{Path: emp + "{id}/delete/", Handler: lamu.NewDynamicView("employees.EmployeeDeleteView")}},
			{Key: "employees.PointsListRoute", Value: lamu.Route{Path: basePts, Handler: lamu.NewDynamicView("employees.PointsListView")}},
			{Key: "employees.PointsCreateRoute", Value: lamu.Route{Path: basePts + "create/", Handler: lamu.NewDynamicView("employees.PointsCreateView")}},
			{Key: "employees.PointsDetailRoute", Value: lamu.Route{Path: basePts + "{id}/", Handler: lamu.NewDynamicView("employees.PointsDetailView")}},
		},
	}
}
