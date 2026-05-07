package p_uniquity_employees

import (
	"context"

	"github.com/UniquityVentures/lago/components"
	"github.com/UniquityVentures/lago/getters"
	"github.com/UniquityVentures/lago/lago"
	"github.com/UniquityVentures/lago/plugins/p_users"
)

func init() {
	registerMenuPages()
	registerEmployeePages()
	registerPointsPages()
}

func pointsDecimalStringGetter(ctxKey string) getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		pd, err := getters.Key[PointsDecimal](ctxKey)(ctx)
		if err != nil {
			return "", err
		}
		return pd.String(), nil
	}
}

func registerMenuPages() {
	lago.RegistryPage.Register("employees.MainMenu", &components.SidebarMenu{
		Title: getters.Static("Employees & points"),
		Back: &components.SidebarMenuItem{
			Title: getters.Static("Back to Home"),
			Url:   lago.RoutePath("dashboard.AppsPage", nil),
		},
		Children: []components.PageInterface{
			&components.SidebarMenuItem{
				Page:  components.Page{Roles: []string{"superuser"}},
				Title: getters.Static("Employees"),
				Url:   lago.RoutePath("employees.DefaultRoute", nil),
				Icon:  "users",
			},
			&components.SidebarMenuItem{
				Page:  components.Page{Roles: []string{"superuser"}},
				Title: getters.Static("Points"),
				Url:   lago.RoutePath("employees.PointsListRoute", nil),
				Icon:  "currency-dollar",
			},
			&components.SidebarMenuItem{
				Page:  components.Page{Roles: []string{"superuser"}},
				Title: getters.Static("New employee"),
				Url:   lago.RoutePath("employees.EmployeeCreateRoute", nil),
				Icon:  "plus",
			},
			&components.SidebarMenuItem{
				Page:  components.Page{Roles: []string{"superuser"}},
				Title: getters.Static("New points entry"),
				Url:   lago.RoutePath("employees.PointsCreateRoute", nil),
				Icon:  "plus",
			},
		},
	})

	lago.RegistryPage.Register("employees.EmployeeDetailMenu", &components.SidebarMenu{
		Title: getters.Format("Employee #%d", getters.Any(getters.Key[uint]("employee.ID"))),
		Back: &components.SidebarMenuItem{
			Title: getters.Static("All employees"),
			Url:   lago.RoutePath("employees.DefaultRoute", nil),
		},
		Children: []components.PageInterface{
			&components.SidebarMenuItem{
				Title: getters.Static("Detail"),
				Url: lago.RoutePath("employees.EmployeeDetailRoute", map[string]getters.Getter[any]{
					"id": getters.Any(getters.Key[uint]("employee.ID")),
				}),
			},
			&components.SidebarMenuItem{
				Page:  components.Page{Roles: []string{"superuser"}},
				Title: getters.Static("Edit"),
				Url: lago.RoutePath("employees.EmployeeUpdateRoute", map[string]getters.Getter[any]{
					"id": getters.Any(getters.Key[uint]("employee.ID")),
				}),
			},
		},
	})

	lago.RegistryPage.Register("employees.PointsDetailMenu", &components.SidebarMenu{
		Title: getters.Format("Points #%d", getters.Any(getters.Key[uint]("pointsTransaction.ID"))),
		Back: &components.SidebarMenuItem{
			Title: getters.Static("All points"),
			Url:   lago.RoutePath("employees.PointsListRoute", nil),
		},
		Children: []components.PageInterface{
			&components.SidebarMenuItem{
				Title: getters.Static("Detail"),
				Url: lago.RoutePath("employees.PointsDetailRoute", map[string]getters.Getter[any]{
					"id": getters.Any(getters.Key[uint]("pointsTransaction.ID")),
				}),
			},
		},
	})
}

func registerEmployeePages() {
	createName := getters.Static("employees.EmployeeCreateForm")
	updateName := getters.Static("employees.EmployeeUpdateForm")
	deleteName := getters.Static("employees.EmployeeDeleteForm")

	userPicker := &components.InputForeignKey[p_users.User]{
		Name:        "UserID",
		Label:       "User",
		Url:         lago.RoutePath("users.SelectRoute", nil),
		Display:     getters.Key[string]("$in.Name"),
		Placeholder: "Select user…",
		Required:    true,
		Getter:      getters.Association[p_users.User, uint](getters.Key[uint]("$in.UserID")),
	}

	lago.RegistryPage.Register("employees.EmployeeTable", &components.ShellScaffold{
		Sidebar: []components.PageInterface{lago.DynamicPage{Name: "employees.MainMenu"}},
		Children: []components.PageInterface{
			&components.DataTable[Employee]{
				UID:     "employee-table",
				Classes: "w-full",
				Data:    getters.Key[components.ObjectList[Employee]]("employees"),
				Actions: []components.PageInterface{
					&components.TableButtonCreate{
						Link: lago.RoutePath("employees.EmployeeCreateRoute", nil),
						Page: components.Page{Roles: []string{"superuser"}},
					},
				},
				RowAttr: getters.RowAttrNavigate(lago.RoutePath("employees.EmployeeDetailRoute", map[string]getters.Getter[any]{
					"id": getters.Any(getters.Key[uint]("$row.ID")),
				})),
				Columns: []components.TableColumn{
					{Label: "User", Name: "User.Name", Children: []components.PageInterface{
						&components.FieldText{Getter: getters.Key[string]("$row.User.Name")},
					}},
					{Label: "Email", Name: "User.Email", Children: []components.PageInterface{
						&components.FieldText{Getter: getters.Key[string]("$row.User.Email")},
					}},
				},
			},
		},
	})

	lago.RegistryPage.Register("employees.EmployeeCreateForm", &components.ShellScaffold{
		Page:    components.Page{Roles: []string{"superuser"}},
		Sidebar: []components.PageInterface{lago.DynamicPage{Name: "employees.MainMenu"}},
		Children: []components.PageInterface{
			&components.FormListenBoostedPost{
				Name:      createName,
				ActionURL: lago.RoutePath("employees.EmployeeCreateRoute", nil),
				Children: []components.PageInterface{
					&components.FormComponent[Employee]{
						Attr:          getters.FormBubbling(createName),
						Title:         "Create employee",
						Subtitle:      "Link a user to an employee record",
						ChildrenInput: []components.PageInterface{userPicker},
						ChildrenAction: []components.PageInterface{
							&components.ButtonSubmit{Label: "Save"},
						},
					},
				},
			},
		},
	})

	lago.RegistryPage.Register("employees.EmployeeUpdateForm", &components.ShellScaffold{
		Page:    components.Page{Roles: []string{"superuser"}},
		Sidebar: []components.PageInterface{lago.DynamicPage{Name: "employees.EmployeeDetailMenu"}},
		Children: []components.PageInterface{
			&components.FormListenBoostedPost{
				Name: updateName,
				ActionURL: lago.RoutePath("employees.EmployeeUpdateRoute", map[string]getters.Getter[any]{
					"id": getters.Any(getters.Key[uint]("employee.ID")),
				}),
				Children: []components.PageInterface{
					&components.FormComponent[Employee]{
						Getter:        getters.Key[Employee]("employee"),
						Attr:          getters.FormBubbling(updateName),
						Title:         "Edit employee",
						Subtitle:      "Change linked user",
						ChildrenInput: []components.PageInterface{userPicker},
						ChildrenAction: []components.PageInterface{
							&components.ContainerRow{
								Classes: "flex flex-wrap justify-between gap-2 mt-2 items-center",
								Children: []components.PageInterface{
									&components.ContainerRow{
										Classes: "flex justify-end gap-2",
										Children: []components.PageInterface{
											&components.ButtonSubmit{Label: "Update"},
											&components.ButtonModalForm{
												Page:        components.Page{Roles: []string{"superuser"}},
												Label:       "Delete",
												Icon:        "trash",
												Name:        deleteName,
												Url:         lago.RoutePath("employees.EmployeeDeleteRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("employee.ID"))}),
												FormPostURL: lago.RoutePath("employees.EmployeeDeleteRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("employee.ID"))}),
												ModalUID:    "employee-delete-modal",
												Classes:     "btn-error",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})

	lago.RegistryPage.Register("employees.EmployeeDeleteForm", &components.Modal{
		Page: components.Page{Roles: []string{"superuser"}},
		UID:  "employee-delete-modal",
		Children: []components.PageInterface{
			&components.DeleteConfirmation{
				Title:   "Delete employee?",
				Message: "This removes the employee record. The user account is not deleted.",
				Attr:    getters.FormBubbling(getters.Key[string]("$get.name")),
			},
		},
	})

	lago.RegistryPage.Register("employees.EmployeeDetail", &components.ShellScaffold{
		Sidebar: []components.PageInterface{lago.DynamicPage{Name: "employees.EmployeeDetailMenu"}},
		Children: []components.PageInterface{
			&components.Detail[Employee]{
				Getter: getters.Key[Employee]("employee"),
				Children: []components.PageInterface{
					&components.ContainerColumn{
						Classes: "p-4",
						Children: []components.PageInterface{
							&components.LabelInline{
								Title: "User",
								Children: []components.PageInterface{
									&components.FieldText{Getter: getters.Key[string]("$in.User.Name")},
								},
							},
							&components.LabelInline{
								Title: "Email",
								Children: []components.PageInterface{
									&components.FieldText{Getter: getters.Key[string]("$in.User.Email")},
								},
							},
						},
					},
				},
			},
		},
	})

	lago.RegistryPage.Register("employees.EmployeeSelectionTable", &components.Modal{
		UID: "employee-selection-modal",
		Children: []components.PageInterface{
			&components.DataTable[Employee]{
				UID:     "employee-selection-table",
				Title:   "Select employee",
				Data:    getters.Key[components.ObjectList[Employee]]("employees"),
				RowAttr: getters.RowAttrSelect("ToEmployeeID", getters.Key[uint]("$row.ID"), getters.Key[string]("$row.User.Name")),
				Columns: []components.TableColumn{
					{Label: "User", Name: "User.Name", Children: []components.PageInterface{
						&components.FieldText{Getter: getters.Key[string]("$row.User.Name")},
					}},
					{Label: "Email", Name: "User.Email", Children: []components.PageInterface{
						&components.FieldText{Getter: getters.Key[string]("$row.User.Email")},
					}},
				},
			},
		},
	})
}

func registerPointsPages() {
	createName := getters.Static("employees.PointsTransactionCreateForm")

	toEmployeePicker := &components.InputForeignKey[Employee]{
		Name:        "ToEmployeeID",
		Label:       "Employee",
		Url:         lago.RoutePath("employees.EmployeeSelectRoute", nil),
		Display:     getters.Key[string]("$in.User.Name"),
		Placeholder: "Select employee…",
		Required:    true,
		Getter:      getters.Association[Employee, uint](getters.Key[uint]("$in.ToEmployeeID")),
	}

	pointsInput := &components.InputText{
		Label:    "Points",
		Name:     "Points",
		Required: true,
		Getter:   pointsDecimalStringGetter("$in.Points"),
	}

	lago.RegistryPage.Register("employees.PointsTransactionTable", &components.ShellScaffold{
		Sidebar: []components.PageInterface{lago.DynamicPage{Name: "employees.MainMenu"}},
		Children: []components.PageInterface{
			&components.DataTable[PointsTransaction]{
				UID:     "points-transaction-table",
				Classes: "w-full",
				Data:    getters.Key[components.ObjectList[PointsTransaction]]("pointsTransactions"),
				Actions: []components.PageInterface{
					&components.TableButtonCreate{
						Link: lago.RoutePath("employees.PointsCreateRoute", nil),
						Page: components.Page{Roles: []string{"superuser"}},
					},
				},
				RowAttr: getters.RowAttrNavigate(lago.RoutePath("employees.PointsDetailRoute", map[string]getters.Getter[any]{
					"id": getters.Any(getters.Key[uint]("$row.ID")),
				})),
				Columns: []components.TableColumn{
					{Label: "Points", Name: "Points", Children: []components.PageInterface{
						&components.FieldText{Getter: pointsDecimalStringGetter("$row.Points")},
					}},
					{Label: "From (superuser)", Name: "FromUser.Name", Children: []components.PageInterface{
						&components.FieldText{Getter: getters.Key[string]("$row.FromUser.Name")},
					}},
					{Label: "To employee", Name: "ToEmployee.User.Name", Children: []components.PageInterface{
						&components.FieldText{Getter: getters.Key[string]("$row.ToEmployee.User.Name")},
					}},
				},
			},
		},
	})

	lago.RegistryPage.Register("employees.PointsTransactionCreateForm", &components.ShellScaffold{
		Page:    components.Page{Roles: []string{"superuser"}},
		Sidebar: []components.PageInterface{lago.DynamicPage{Name: "employees.MainMenu"}},
		Children: []components.PageInterface{
			&components.FormListenBoostedPost{
				Name:      createName,
				ActionURL: lago.RoutePath("employees.PointsCreateRoute", nil),
				Children: []components.PageInterface{
					&components.FormComponent[PointsTransaction]{
						Attr:          getters.FormBubbling(createName),
						Title:         "Create points transaction",
						Subtitle:      "From user is the signed-in superuser (set automatically).",
						ChildrenInput: []components.PageInterface{toEmployeePicker, pointsInput},
						ChildrenAction: []components.PageInterface{
							&components.ButtonSubmit{Label: "Save"},
						},
					},
				},
			},
		},
	})

	lago.RegistryPage.Register("employees.PointsTransactionDetail", &components.ShellScaffold{
		Sidebar: []components.PageInterface{lago.DynamicPage{Name: "employees.PointsDetailMenu"}},
		Children: []components.PageInterface{
			&components.Detail[PointsTransaction]{
				Getter: getters.Key[PointsTransaction]("pointsTransaction"),
				Children: []components.PageInterface{
					&components.ContainerColumn{
						Classes: "p-4",
						Children: []components.PageInterface{
							&components.LabelInline{
								Title: "Points",
								Children: []components.PageInterface{
									&components.FieldText{Getter: pointsDecimalStringGetter("$in.Points")},
								},
							},
							&components.LabelInline{
								Title: "From",
								Children: []components.PageInterface{
									&components.FieldText{Getter: getters.Key[string]("$in.FromUser.Name")},
								},
							},
							&components.LabelInline{
								Title: "To employee",
								Children: []components.PageInterface{
									&components.FieldText{Getter: getters.Key[string]("$in.ToEmployee.User.Name")},
								},
							},
						},
					},
				},
			},
		},
	})
}
