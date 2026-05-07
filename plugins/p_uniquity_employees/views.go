package p_uniquity_employees

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/UniquityVentures/lago/getters"
	"github.com/UniquityVentures/lago/lago"
	"github.com/UniquityVentures/lago/plugins/p_users"
	"github.com/UniquityVentures/lago/registry"
	"github.com/UniquityVentures/lago/views"
	"gorm.io/gorm"
)

// superuserOnlyLayer returns 401 unless the authenticated user is a superuser.
type superuserOnlyLayer struct{}

func (superuserOnlyLayer) Next(_ views.View, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := p_users.UserFromContext(r.Context(), "employees.superuserOnlyLayer")
		if !user.IsSuperuser {
			slog.Error("employees.superuserOnlyLayer: forbidden", "user_id", user.ID)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type employeeListPreload struct{}

func (employeeListPreload) Patch(_ views.View, _ *http.Request, query gorm.ChainInterface[Employee]) gorm.ChainInterface[Employee] {
	return query.Preload("User", nil)
}

type pointsListPreload struct{}

func (pointsListPreload) Patch(_ views.View, _ *http.Request, query gorm.ChainInterface[PointsTransaction]) gorm.ChainInterface[PointsTransaction] {
	return query.Preload("FromUser", nil).Preload("ToEmployee", nil).Preload("ToEmployee.User", nil)
}

type pointsDetailPreload struct{}

func (pointsDetailPreload) Patch(_ views.View, _ *http.Request, query gorm.ChainInterface[PointsTransaction]) gorm.ChainInterface[PointsTransaction] {
	return query.Preload("FromUser", nil).Preload("ToEmployee", nil).Preload("ToEmployee.User", nil)
}

type pointsFormFromUserPatcher struct{}

func (pointsFormFromUserPatcher) Patch(_ views.View, r *http.Request, formData map[string]any, formErrors map[string]error) (map[string]any, map[string]error) {
	user := p_users.UserFromContext(r.Context(), "pointsFormFromUserPatcher")
	if !user.IsSuperuser {
		formErrors["_form"] = errors.New("only superusers can create points transactions")
		return formData, formErrors
	}
	formData["FromUserID"] = user.ID
	return formData, formErrors
}

func init() {
	// --- Employee ---
	lago.RegistryView.Register("employees.EmployeeListView",
		lago.GetPageView("employees.EmployeeTable").
			WithLayer("users.auth", p_users.AuthenticationLayer{}).
			WithLayer("employees.superuser", superuserOnlyLayer{}).
			WithLayer("employees.employee_list", views.LayerList[Employee]{
				Key: getters.Static("employees"),
				QueryPatchers: views.QueryPatchers[Employee]{
					registry.Pair[string, views.QueryPatcher[Employee]]{Key: "employees.preload_user", Value: employeeListPreload{}},
				},
			}))

	lago.RegistryView.Register("employees.EmployeeDetailView",
		lago.GetPageView("employees.EmployeeDetail").
			WithLayer("users.auth", p_users.AuthenticationLayer{}).
			WithLayer("employees.superuser", superuserOnlyLayer{}).
			WithLayer("employees.employee_detail", views.LayerDetail[Employee]{
				Key:          getters.Static("employee"),
				PathParamKey: getters.Static("id"),
				QueryPatchers: views.QueryPatchers[Employee]{
					registry.Pair[string, views.QueryPatcher[Employee]]{Key: "employees.preload_user", Value: employeeListPreload{}},
				},
			}))

	lago.RegistryView.Register("employees.EmployeeCreateView",
		lago.GetPageView("employees.EmployeeCreateForm").
			WithLayer("users.auth", p_users.AuthenticationLayer{}).
			WithLayer("employees.superuser", superuserOnlyLayer{}).
			WithLayer("employees.employee_create", views.LayerCreate[Employee]{
				SuccessURL: lago.RoutePath("employees.EmployeeDetailRoute", map[string]getters.Getter[any]{
					"id": getters.Any(getters.Key[uint]("$id")),
				}),
			}))

	lago.RegistryView.Register("employees.EmployeeUpdateView",
		lago.GetPageView("employees.EmployeeUpdateForm").
			WithLayer("users.auth", p_users.AuthenticationLayer{}).
			WithLayer("employees.superuser", superuserOnlyLayer{}).
			WithLayer("employees.employee_detail", views.LayerDetail[Employee]{
				Key:          getters.Static("employee"),
				PathParamKey: getters.Static("id"),
				QueryPatchers: views.QueryPatchers[Employee]{
					registry.Pair[string, views.QueryPatcher[Employee]]{Key: "employees.preload_user", Value: employeeListPreload{}},
				},
			}).
			WithLayer("employees.employee_update", views.LayerUpdate[Employee]{
				Key:        getters.Static("employee"),
				SuccessURL: lago.RoutePath("employees.DefaultRoute", nil),
			}))

	lago.RegistryView.Register("employees.EmployeeDeleteView",
		lago.GetPageView("employees.EmployeeDeleteForm").
			WithLayer("users.auth", p_users.AuthenticationLayer{}).
			WithLayer("employees.superuser", superuserOnlyLayer{}).
			WithLayer("employees.employee_detail", views.LayerDetail[Employee]{
				Key:          getters.Static("employee"),
				PathParamKey: getters.Static("id"),
				QueryPatchers: views.QueryPatchers[Employee]{
					registry.Pair[string, views.QueryPatcher[Employee]]{Key: "employees.preload_user", Value: employeeListPreload{}},
				},
			}).
			WithLayer("employees.employee_delete", views.LayerDelete[Employee]{
				Key:        getters.Static("employee"),
				SuccessURL: lago.RoutePath("employees.DefaultRoute", nil),
			}))

	lago.RegistryView.Register("employees.EmployeeSelectView",
		lago.GetPageView("employees.EmployeeSelectionTable").
			WithLayer("users.auth", p_users.AuthenticationLayer{}).
			WithLayer("employees.superuser", superuserOnlyLayer{}).
			WithLayer("employees.employee_select_list", views.LayerList[Employee]{
				Key: getters.Static("employees"),
				QueryPatchers: views.QueryPatchers[Employee]{
					registry.Pair[string, views.QueryPatcher[Employee]]{Key: "employees.preload_user", Value: employeeListPreload{}},
				},
			}))

	// --- Points (no update view) ---
	lago.RegistryView.Register("employees.PointsListView",
		lago.GetPageView("employees.PointsTransactionTable").
			WithLayer("users.auth", p_users.AuthenticationLayer{}).
			WithLayer("employees.superuser", superuserOnlyLayer{}).
			WithLayer("employees.points_list", views.LayerList[PointsTransaction]{
				Key: getters.Static("pointsTransactions"),
				QueryPatchers: views.QueryPatchers[PointsTransaction]{
					registry.Pair[string, views.QueryPatcher[PointsTransaction]]{Key: "employees.points_preload", Value: pointsListPreload{}},
				},
			}))

	lago.RegistryView.Register("employees.PointsDetailView",
		lago.GetPageView("employees.PointsTransactionDetail").
			WithLayer("users.auth", p_users.AuthenticationLayer{}).
			WithLayer("employees.superuser", superuserOnlyLayer{}).
			WithLayer("employees.points_detail", views.LayerDetail[PointsTransaction]{
				Key:          getters.Static("pointsTransaction"),
				PathParamKey: getters.Static("id"),
				QueryPatchers: views.QueryPatchers[PointsTransaction]{
					registry.Pair[string, views.QueryPatcher[PointsTransaction]]{Key: "employees.points_preload", Value: pointsDetailPreload{}},
				},
			}))

	lago.RegistryView.Register("employees.PointsCreateView",
		lago.GetPageView("employees.PointsTransactionCreateForm").
			WithLayer("users.auth", p_users.AuthenticationLayer{}).
			WithLayer("employees.superuser", superuserOnlyLayer{}).
			WithLayer("employees.points_create", views.LayerCreate[PointsTransaction]{
				SuccessURL: lago.RoutePath("employees.PointsDetailRoute", map[string]getters.Getter[any]{
					"id": getters.Any(getters.Key[uint]("$id")),
				}),
				FormPatchers: views.FormPatchers{
					registry.Pair[string, views.FormPatcher]{Key: "employees.points_from_user", Value: pointsFormFromUserPatcher{}},
				},
			}))
}
