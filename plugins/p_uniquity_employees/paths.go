package p_uniquity_employees

import "github.com/UniquityVentures/lago/lago"

func init() {
	registerEmployeeRoutes()
	registerPointsRoutes()
}

func registerEmployeeRoutes() {
	// List / create / select live directly under AppUrl. Per-record routes use
	// /employees/emp/{id}/... so they cannot collide with /employees/points/...
	// (e.g. /employees/points/delete/ vs /employees/{id}/delete/).
	emp := AppUrl + "emp/"
	_ = lago.RegistryRoute.Register("employees.DefaultRoute", lago.Route{
		Path:    AppUrl,
		Handler: lago.NewDynamicView("employees.EmployeeListView"),
	})
	_ = lago.RegistryRoute.Register("employees.EmployeeCreateRoute", lago.Route{
		Path:    AppUrl + "create/",
		Handler: lago.NewDynamicView("employees.EmployeeCreateView"),
	})
	_ = lago.RegistryRoute.Register("employees.EmployeeSelectRoute", lago.Route{
		Path:    AppUrl + "select/",
		Handler: lago.NewDynamicView("employees.EmployeeSelectView"),
	})
	_ = lago.RegistryRoute.Register("employees.EmployeeDetailRoute", lago.Route{
		Path:    emp + "{id}/",
		Handler: lago.NewDynamicView("employees.EmployeeDetailView"),
	})
	_ = lago.RegistryRoute.Register("employees.EmployeeUpdateRoute", lago.Route{
		Path:    emp + "{id}/edit/",
		Handler: lago.NewDynamicView("employees.EmployeeUpdateView"),
	})
	_ = lago.RegistryRoute.Register("employees.EmployeeDeleteRoute", lago.Route{
		Path:    emp + "{id}/delete/",
		Handler: lago.NewDynamicView("employees.EmployeeDeleteView"),
	})
}

func registerPointsRoutes() {
	base := AppUrl + "points/"
	_ = lago.RegistryRoute.Register("employees.PointsListRoute", lago.Route{
		Path:    base,
		Handler: lago.NewDynamicView("employees.PointsListView"),
	})
	_ = lago.RegistryRoute.Register("employees.PointsCreateRoute", lago.Route{
		Path:    base + "create/",
		Handler: lago.NewDynamicView("employees.PointsCreateView"),
	})
	_ = lago.RegistryRoute.Register("employees.PointsDetailRoute", lago.Route{
		Path:    base + "{id}/",
		Handler: lago.NewDynamicView("employees.PointsDetailView"),
	})
}
