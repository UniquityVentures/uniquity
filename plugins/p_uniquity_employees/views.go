package p_uniquity_employees

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/UniquityVentures/lamu/fields"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/plugins/p_users"
	"github.com/UniquityVentures/lamu/registry"
	"github.com/UniquityVentures/lamu/views"
	"gorm.io/gorm"
)


// employeePointsTotalContextKey is where [employeeDetailPointsTotalLayer] stores the
// display string for SUM(points) (must not contain '.' — [getters.Key] path rules).
const employeePointsTotalContextKey = "employeePointsTotal"

// employeeDetailPointsTotalLayer attaches the employee’s lifetime points total for the detail page.
type employeeDetailPointsTotalLayer struct{}

func (employeeDetailPointsTotalLayer) Next(_ views.View, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		emp, err := getters.Key[Employee]("employee")(ctx)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		db, err := getters.DBFromContext(ctx)
		if err != nil {
			slog.Error("employees.employee_points_total: db", "error", err)
			ctx = context.WithValue(ctx, employeePointsTotalContextKey, "—")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		var row struct {
			Sum fields.DecimalSix `gorm:"column:sum"`
		}
		q := db.Model(&PointsTransaction{}).Where("to_employee_id = ?", emp.ID)
		if err := q.Select("COALESCE(SUM(points), 0) AS sum").Scan(&row).Error; err != nil {
			slog.Error("employees.employee_points_total: query", "error", err, "employeeID", emp.ID)
			ctx = context.WithValue(ctx, employeePointsTotalContextKey, "—")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		ctx = context.WithValue(ctx, employeePointsTotalContextKey, row.Sum.String())
		next.ServeHTTP(w, r.WithContext(ctx))
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

// PointsFormFromUserPatcher sets FromUserID from the signed-in superuser when creating a points transaction.
type PointsFormFromUserPatcher struct{}

func (PointsFormFromUserPatcher) Patch(_ views.View, r *http.Request, formData map[string]any, formErrors map[string]error) (map[string]any, map[string]error) {
	user := p_users.UserFromContext(r.Context(), "employees.PointsFormFromUserPatcher")
	if !user.IsSuperuser {
		formErrors["_form"] = errors.New("only superusers can create points transactions")
		return formData, formErrors
	}
	formData["FromUserID"] = user.ID
	return formData, formErrors
}

func pluginViews() lamu.PluginFeatures[*views.View] {
	return lamu.PluginFeatures[*views.View]{
		Entries: []registry.Pair[string, *views.View]{
			{
				Key: "employees.EmployeeListView",
				Value: lamu.GetPageView("employees.EmployeeTable").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("employees.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("employees.employee_list", views.LayerList[Employee]{
						Key: getters.Static("employees"),
						QueryPatchers: views.QueryPatchers[Employee]{
							{Key: "employees.preload_user", Value: employeeListPreload{}},
						},
					}),
			},
			{
				Key: "employees.EmployeeDetailView",
				Value: lamu.GetPageView("employees.EmployeeDetail").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("employees.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("employees.employee_detail", views.LayerDetail[Employee]{
						Key:          getters.Static("employee"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[Employee]{
							{Key: "employees.preload_user", Value: employeeListPreload{}},
						},
					}).
					WithLayer("employees.employee_points_total", employeeDetailPointsTotalLayer{}),
			},
			{
				Key: "employees.EmployeeCreateView",
				Value: lamu.GetPageView("employees.EmployeeCreateForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("employees.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("employees.employee_create", views.LayerCreate[Employee]{
						SuccessURL: lamu.RoutePath("employees.EmployeeDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("$id")),
						}),
					}),
			},
			{
				Key: "employees.EmployeeUpdateView",
				Value: lamu.GetPageView("employees.EmployeeUpdateForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("employees.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("employees.employee_detail", views.LayerDetail[Employee]{
						Key:          getters.Static("employee"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[Employee]{
							{Key: "employees.preload_user", Value: employeeListPreload{}},
						},
					}).
					WithLayer("employees.employee_update", views.LayerUpdate[Employee]{
						Key:        getters.Static("employee"),
						SuccessURL: lamu.RoutePath("employees.DefaultRoute", nil),
					}),
			},
			{
				Key: "employees.EmployeeDeleteView",
				Value: lamu.GetPageView("employees.EmployeeDeleteForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("employees.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("employees.employee_detail", views.LayerDetail[Employee]{
						Key:          getters.Static("employee"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[Employee]{
							{Key: "employees.preload_user", Value: employeeListPreload{}},
						},
					}).
					WithLayer("employees.employee_delete", views.LayerDelete[Employee]{
						Key:        getters.Static("employee"),
						SuccessURL: lamu.RoutePath("employees.DefaultRoute", nil),
					}),
			},
			{
				Key: "employees.EmployeeSelectView",
				Value: lamu.GetPageView("employees.EmployeeSelectionTable").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("employees.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("employees.employee_select_list", views.LayerList[Employee]{
						Key: getters.Static("employees"),
						QueryPatchers: views.QueryPatchers[Employee]{
							{Key: "employees.preload_user", Value: employeeListPreload{}},
						},
					}),
			},
			{
				Key: "employees.PointsListView",
				Value: lamu.GetPageView("employees.PointsTransactionTable").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("employees.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("employees.points_list", views.LayerList[PointsTransaction]{
						Key: getters.Static("pointsTransactions"),
						QueryPatchers: views.QueryPatchers[PointsTransaction]{
							{Key: "employees.points_preload", Value: pointsListPreload{}},
						},
					}),
			},
			{
				Key: "employees.PointsDetailView",
				Value: lamu.GetPageView("employees.PointsTransactionDetail").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("employees.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("employees.points_detail", views.LayerDetail[PointsTransaction]{
						Key:          getters.Static("pointsTransaction"),
						PathParamKey: getters.Static("id"),
						QueryPatchers: views.QueryPatchers[PointsTransaction]{
							{Key: "employees.points_preload", Value: pointsDetailPreload{}},
						},
					}),
			},
			{
				Key: "employees.PointsCreateView",
				Value: lamu.GetPageView("employees.PointsTransactionCreateForm").
					WithLayer("p_users.auth", p_users.AuthenticationLayer{}).
					WithLayer("employees.superuser", p_users.SuperuserOnlyLayer{}).
					WithLayer("employees.points_create", views.LayerCreate[PointsTransaction]{
						SuccessURL: lamu.RoutePath("employees.PointsDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("$id")),
						}),
						FormPatchers: views.FormPatchers{
							{Key: "employees.points_from_user", Value: PointsFormFromUserPatcher{}},
						},
					}),
			},
		},
	}
}
