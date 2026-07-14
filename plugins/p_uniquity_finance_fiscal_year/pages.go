package p_uniquity_finance_fiscal_year

import (
	"time"

	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/components"
	"github.com/lariv-in/lago/getters"
	"github.com/lariv-in/lago/registry"
)

const financeAccountsMainMenuFiscalYearsLinkKey = "finance_fiscal_years.FinanceAccountsMainMenuLink"

func patchFinanceAccountsMainMenuForFiscalYears(page components.PageInterface) components.PageInterface {
	menu, ok := page.(*components.SidebarMenu)
	if !ok {
		panic("p_uniquity_finance_fiscal_year: finance_accounts.MainMenu must be *components.SidebarMenu")
	}
	for _, ch := range menu.Children {
		if item, ok := ch.(*components.SidebarMenuItem); ok && item.GetKey() == financeAccountsMainMenuFiscalYearsLinkKey {
			return menu
		}
	}
	newChildren := append([]components.PageInterface{}, menu.Children...)
	newChildren = append(newChildren, &components.SidebarMenuItem{
		Page:  components.Page{Key: financeAccountsMainMenuFiscalYearsLinkKey, Roles: []string{"superuser"}},
		Title: getters.Static("Fiscal years"),
		Url:   lago.RoutePath("finance_fiscal_years.DefaultRoute", nil),
		Icon:  "calendar-days",
	})
	cloned := *menu
	cloned.Children = newChildren
	return &cloned
}

func pluginPages() lago.PluginFeatures[components.PageInterface] {
	e := pageEntriesFiscalYearMenus()
	e = append(e, pageEntriesFiscalYearPages()...)
	e = append(e, pageEntriesFiscalYearSelectionPages()...)
	return lago.PluginFeatures[components.PageInterface]{
		Entries: e,
		Patches: []registry.Pair[string, func(components.PageInterface) components.PageInterface]{
			{Key: "finance_accounts.MainMenu", Value: patchFinanceAccountsMainMenuForFiscalYears},
		},
	}
}

func pageEntriesFiscalYearMenus() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "finance_fiscal_years.FiscalYearDetailMenu", Value: &components.SidebarMenu{
			Title: getters.Format("%s", getters.Any(getters.Key[string]("fiscal_year.Name"))),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("All fiscal years"),
				Url:   lago.RoutePath("finance_fiscal_years.DefaultRoute", nil),
			},
			Children: []components.PageInterface{
				&components.SidebarMenuItem{
					Title: getters.Static("Detail"),
					Url: lago.RoutePath("finance_fiscal_years.FiscalYearDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("fiscal_year.ID")),
					}),
				},
				&components.SidebarMenuItem{
					Page:  components.Page{Roles: []string{"superuser"}},
					Title: getters.Static("Edit"),
					Url: lago.RoutePath("finance_fiscal_years.FiscalYearUpdateRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("fiscal_year.ID")),
					}),
				},
			},
		}},
	}
}

func fiscalYearCreateFormInputs() []components.PageInterface {
	return []components.PageInterface{
		&components.ContainerError{
			Error: getters.Key[error]("$error.Code"),
			Children: []components.PageInterface{
				&components.InputText{Name: "Code", Label: "Code", Required: true},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Name"),
			Children: []components.PageInterface{
				&components.InputText{Name: "Name", Label: "Name", Required: true},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Start"),
			Children: []components.PageInterface{
				&components.InputDatetime{Name: "Start", Label: "Start", Required: true},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.End"),
			Children: []components.PageInterface{
				&components.InputDatetime{Name: "End", Label: "End", Required: true},
			},
		},
		&components.InputCheckbox{Name: "IsActive", Label: "Active", Getter: getters.Key[bool]("$in.IsActive")},
	}
}

func fiscalYearUpdateFormInputs() []components.PageInterface {
	return []components.PageInterface{
		&components.ContainerError{
			Error: getters.Key[error]("$error.Code"),
			Children: []components.PageInterface{
				&components.InputText{Name: "Code", Label: "Code", Required: true, Getter: getters.Key[string]("$in.Code")},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Name"),
			Children: []components.PageInterface{
				&components.InputText{Name: "Name", Label: "Name", Required: true, Getter: getters.Key[string]("$in.Name")},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Start"),
			Children: []components.PageInterface{
				&components.InputDatetime{Name: "Start", Label: "Start", Required: true, Getter: getters.Key[time.Time]("$in.Start")},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.End"),
			Children: []components.PageInterface{
				&components.InputDatetime{Name: "End", Label: "End", Required: true, Getter: getters.Key[time.Time]("$in.End")},
			},
		},
		&components.InputCheckbox{Name: "IsActive", Label: "Active", Getter: getters.Key[bool]("$in.IsActive")},
	}
}

func pageEntriesFiscalYearPages() []registry.Pair[string, components.PageInterface] {
	createName := getters.Static("finance_fiscal_years.FiscalYearCreateForm")
	updateName := getters.Static("finance_fiscal_years.FiscalYearUpdateForm")
	deleteName := getters.Static("finance_fiscal_years.FiscalYearDeleteForm")

	return []registry.Pair[string, components.PageInterface]{
		{Key: "finance_fiscal_years.FiscalYearTable", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lago.DynamicPage{Name: "finance_accounts.MainMenu"}},
			Children: []components.PageInterface{
				&components.DataTable[FiscalYear]{
					UID:     "finance-fiscal-year-table",
					Classes: "w-full",
					Data:    getters.Key[components.ObjectList[FiscalYear]]("fiscal_years"),
					Actions: []components.PageInterface{
						&components.TableButtonCreate{
							Link: lago.RoutePath("finance_fiscal_years.FiscalYearCreateRoute", nil),
							Page: components.Page{Roles: []string{"superuser"}},
						},
					},
					RowAttr: getters.RowAttrNavigate(lago.RoutePath("finance_fiscal_years.FiscalYearDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("$row.ID")),
					})),
					Columns: []components.TableColumn{
						{Label: "Code", Name: "Code", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Code")},
						}},
						{Label: "Name", Name: "Name", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Name")},
						}},
						{Label: "Start", Name: "Start", Children: []components.PageInterface{
							&components.FieldDatetime{Getter: getters.Key[time.Time]("$row.Start")},
						}},
						{Label: "End", Name: "End", Children: []components.PageInterface{
							&components.FieldDatetime{Getter: getters.Key[time.Time]("$row.End")},
						}},
						{Label: "Active", Name: "IsActive", Children: []components.PageInterface{
							&components.FieldCheckbox{Getter: getters.Key[bool]("$row.IsActive")},
						}},
					},
				},
			},
		}},
		{Key: "finance_fiscal_years.FiscalYearCreateForm", Value: &components.ShellScaffold{
			Page:    components.Page{Roles: []string{"superuser"}},
			Sidebar: []components.PageInterface{lago.DynamicPage{Name: "finance_accounts.MainMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name:      createName,
					ActionURL: lago.RoutePath("finance_fiscal_years.FiscalYearCreateRoute", nil),
					Children: []components.PageInterface{
						&components.FormComponent[FiscalYear]{
							Attr:          getters.FormBubbling(createName),
							Title:         "Create fiscal year",
							Subtitle:      "Code, name, and period bounds",
							ChildrenInput: fiscalYearCreateFormInputs(),
							ChildrenAction: []components.PageInterface{
								&components.ButtonSubmit{Label: "Save"},
							},
						},
					},
				},
			},
		}},
		{Key: "finance_fiscal_years.FiscalYearUpdateForm", Value: &components.ShellScaffold{
			Page:    components.Page{Roles: []string{"superuser"}},
			Sidebar: []components.PageInterface{lago.DynamicPage{Name: "finance_fiscal_years.FiscalYearDetailMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name: updateName,
					ActionURL: lago.RoutePath("finance_fiscal_years.FiscalYearUpdateRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("fiscal_year.ID")),
					}),
					Children: []components.PageInterface{
						&components.FormComponent[FiscalYear]{
							Getter:        getters.Key[FiscalYear]("fiscal_year"),
							Attr:          getters.FormBubbling(updateName),
							Title:         "Edit fiscal year",
							Subtitle:      "Update code, name, and period bounds",
							ChildrenInput: fiscalYearUpdateFormInputs(),
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
													Url:         lago.RoutePath("finance_fiscal_years.FiscalYearDeleteRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("fiscal_year.ID"))}),
													FormPostURL: lago.RoutePath("finance_fiscal_years.FiscalYearDeleteRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("fiscal_year.ID"))}),
													ModalUID:    "finance-fiscal-year-delete-modal",
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
		}},
		{Key: "finance_fiscal_years.FiscalYearDeleteForm", Value: &components.Modal{
			Page: components.Page{Roles: []string{"superuser"}},
			UID:  "finance-fiscal-year-delete-modal",
			Children: []components.PageInterface{
				&components.DeleteConfirmation{
					Title:   "Delete fiscal year?",
					Message: "This permanently removes the fiscal year record.",
					Attr:    getters.FormBubbling(getters.Key[string]("$get.name")),
				},
			},
		}},
		{Key: "finance_fiscal_years.FiscalYearDetail", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lago.DynamicPage{Name: "finance_fiscal_years.FiscalYearDetailMenu"}},
			Children: []components.PageInterface{
				&components.Detail[FiscalYear]{
					Getter: getters.Key[FiscalYear]("fiscal_year"),
					Children: []components.PageInterface{
						&components.ContainerColumn{
							Classes: "p-4",
							Children: []components.PageInterface{
								&components.LabelInline{Title: "Code", Children: []components.PageInterface{
									&components.FieldText{Getter: getters.Key[string]("$in.Code")},
								}},
								&components.LabelInline{Title: "Name", Children: []components.PageInterface{
									&components.FieldText{Getter: getters.Key[string]("$in.Name")},
								}},
								&components.LabelInline{Title: "Start", Children: []components.PageInterface{
									&components.FieldDatetime{Getter: getters.Key[time.Time]("$in.Start")},
								}},
								&components.LabelInline{Title: "End", Children: []components.PageInterface{
									&components.FieldDatetime{Getter: getters.Key[time.Time]("$in.End")},
								}},
								&components.LabelInline{Title: "Active", Children: []components.PageInterface{
									&components.FieldCheckbox{Getter: getters.Key[bool]("$in.IsActive")},
								}},
							},
						},
					},
				},
			},
		}},
	}
}
