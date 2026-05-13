package p_uniquity_currencies

import (
	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pluginPages() lamu.PluginFeatures[components.PageInterface] {
	return lamu.PluginFeatures[components.PageInterface]{
		Entries: pageEntries(),
	}
}

func currencyFormInputs() []components.PageInterface {
	return []components.PageInterface{
		&components.ContainerError{
			Error: getters.Key[error]("$error.Code"),
			Children: []components.PageInterface{
				&components.InputText{
					Label:    "Code (ISO 4217)",
					Name:     "Code",
					Required: true,
					Getter:   getters.Key[string]("$in.Code"),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Name"),
			Children: []components.PageInterface{
				&components.InputText{
					Label:    "Name",
					Name:     "Name",
					Required: true,
					Getter:   getters.Key[string]("$in.Name"),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Symbol"),
			Children: []components.PageInterface{
				&components.InputText{
					Label:  "Symbol",
					Name:   "Symbol",
					Getter: getters.Key[string]("$in.Symbol"),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.IsActive"),
			Children: []components.PageInterface{
				&components.InputCheckbox{
					Label:    "Active",
					Name:     "IsActive",
					Required: true,
					Getter:   getters.Key[bool]("$in.IsActive"),
				},
			},
		},
	}
}

func pageEntries() []registry.Pair[string, components.PageInterface] {
	createName := getters.Static("currencies.CurrencyCreateForm")
	updateName := getters.Static("currencies.CurrencyUpdateForm")
	deleteName := getters.Static("currencies.CurrencyDeleteForm")

	return []registry.Pair[string, components.PageInterface]{
		{Key: "currencies.MainMenu", Value: &components.SidebarMenu{
			Title: getters.Static("Currencies"),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("Back to Home"),
				Url:   lamu.RoutePath("dashboard.AppsPage", nil),
			},
			Children: []components.PageInterface{
				&components.SidebarMenuItem{
					Page:  components.Page{Roles: []string{"superuser"}},
					Title: getters.Static("All currencies"),
					Url:   lamu.RoutePath("currencies.CurrencyListRoute", nil),
					Icon:  "banknotes",
				},
				&components.SidebarMenuItem{
					Page:  components.Page{Roles: []string{"superuser"}},
					Title: getters.Static("New currency"),
					Url:   lamu.RoutePath("currencies.CurrencyCreateRoute", nil),
					Icon:  "plus",
				},
			},
		}},
		{Key: "currencies.CurrencyDetailMenu", Value: &components.SidebarMenu{
			Title: getters.Format("%s — %s", getters.Any(getters.Key[string]("currency.Code")), getters.Any(getters.Key[string]("currency.Name"))),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("All currencies"),
				Url:   lamu.RoutePath("currencies.CurrencyListRoute", nil),
			},
			Children: []components.PageInterface{
				&components.SidebarMenuItem{
					Title: getters.Static("Detail"),
					Url: lamu.RoutePath("currencies.CurrencyDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("currency.ID")),
					}),
				},
				&components.SidebarMenuItem{
					Page:  components.Page{Roles: []string{"superuser"}},
					Title: getters.Static("Edit"),
					Url: lamu.RoutePath("currencies.CurrencyUpdateRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("currency.ID")),
					}),
				},
			},
		}},
		{Key: "currencies.CurrencyTable", Value: &components.ShellScaffold{
			Page:    components.Page{Roles: []string{"superuser"}},
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "currencies.MainMenu"}},
			Children: []components.PageInterface{
				&components.DataTable[Currency]{
					UID:     "currencies-table",
					Classes: "w-full",
					Data:    getters.Key[components.ObjectList[Currency]]("currencies"),
					Actions: []components.PageInterface{
						&components.TableButtonCreate{
							Link: lamu.RoutePath("currencies.CurrencyCreateRoute", nil),
							Page: components.Page{Roles: []string{"superuser"}},
						},
					},
					RowAttr: getters.RowAttrNavigate(lamu.RoutePath("currencies.CurrencyDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("$row.ID")),
					})),
					Columns: []components.TableColumn{
						{Label: "Code", Name: "Code", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Code")},
						}},
						{Label: "Name", Name: "Name", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Name")},
						}},
						{Label: "Symbol", Name: "Symbol", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Symbol")},
						}},
						{Label: "Active", Name: "IsActive", Children: []components.PageInterface{
							&components.FieldCheckbox{Getter: getters.Key[bool]("$row.IsActive")},
						}},
					},
				},
			},
		}},
		{Key: "currencies.CurrencySelectionTable", Value: &components.ShellScaffold{
			Page:    components.Page{Roles: []string{"superuser"}},
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "currencies.MainMenu"}},
			Children: []components.PageInterface{
				&components.DataTable[Currency]{
					UID:     "currencies-select-table",
					Classes: "w-full",
					Data:    getters.Key[components.ObjectList[Currency]]("currencies"),
					RowAttr: getters.RowAttrSelect("CurrencyID",
						getters.Key[uint]("$row.ID"),
						getters.Format("%s — %s", getters.Any(getters.Key[string]("$row.Code")), getters.Any(getters.Key[string]("$row.Name"))),
					),
					Columns: []components.TableColumn{
						{Label: "Code", Name: "Code", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Code")},
						}},
						{Label: "Name", Name: "Name", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Name")},
						}},
					},
				},
			},
		}},
		{Key: "currencies.CurrencyCreateForm", Value: &components.ShellScaffold{
			Page:    components.Page{Roles: []string{"superuser"}},
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "currencies.MainMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name:      createName,
					ActionURL: lamu.RoutePath("currencies.CurrencyCreateRoute", nil),
					Children: []components.PageInterface{
						&components.FormComponent[Currency]{
							Attr:          getters.FormBubbling(createName),
							Title:         "Create currency",
							Subtitle:      "Add an ISO 4217 currency",
							ChildrenInput: currencyFormInputs(),
							ChildrenAction: []components.PageInterface{
								&components.ButtonSubmit{Label: "Save"},
							},
						},
					},
				},
			},
		}},
		{Key: "currencies.CurrencyUpdateForm", Value: &components.ShellScaffold{
			Page:    components.Page{Roles: []string{"superuser"}},
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "currencies.CurrencyDetailMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name: updateName,
					ActionURL: lamu.RoutePath("currencies.CurrencyUpdateRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("currency.ID")),
					}),
					Children: []components.PageInterface{
						&components.FormComponent[Currency]{
							Getter:        getters.Key[Currency]("currency"),
							Attr:          getters.FormBubbling(updateName),
							Title:         "Edit currency",
							Subtitle:      "Update ISO currency details",
							ChildrenInput: currencyFormInputs(),
							ChildrenAction: []components.PageInterface{
								&components.ContainerRow{
									Classes: "flex flex-wrap justify-between gap-2 mt-2 items-center",
									Children: []components.PageInterface{
										&components.ContainerRow{
											Classes: "flex justify-end gap-2",
											Children: []components.PageInterface{
												&components.ButtonSubmit{Label: "Save"},
												&components.ButtonModalForm{
													Page:        components.Page{Roles: []string{"superuser"}},
													Label:       "Delete",
													Icon:        "trash",
													Name:        deleteName,
													Url: lamu.RoutePath("currencies.CurrencyDeleteRoute", map[string]getters.Getter[any]{
														"id": getters.Any(getters.Key[uint]("currency.ID")),
													}),
													FormPostURL: lamu.RoutePath("currencies.CurrencyDeleteRoute", map[string]getters.Getter[any]{
														"id": getters.Any(getters.Key[uint]("currency.ID")),
													}),
													ModalUID: "currency-delete-modal",
													Classes:  "btn-error",
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
		{Key: "currencies.CurrencyDeleteForm", Value: &components.Modal{
			Page: components.Page{Roles: []string{"superuser"}},
			UID:  "currency-delete-modal",
			Children: []components.PageInterface{
				&components.DeleteConfirmation{
					Title:   "Delete currency?",
					Message: "This permanently removes this currency. References from entities or accounts may block deletion.",
					Attr:    getters.FormBubbling(deleteName),
				},
			},
		}},
		{Key: "currencies.CurrencyDetail", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "currencies.CurrencyDetailMenu"}},
			Children: []components.PageInterface{
				&components.Detail[Currency]{
					Getter: getters.Key[Currency]("currency"),
					Children: []components.PageInterface{
						&components.ContainerColumn{
							Classes: "p-4",
							Children: []components.PageInterface{
								&components.LabelInline{
									Title: "Code",
									Children: []components.PageInterface{
										&components.FieldText{Getter: getters.Key[string]("$in.Code")},
									},
								},
								&components.LabelInline{
									Title: "Name",
									Children: []components.PageInterface{
										&components.FieldText{Getter: getters.Key[string]("$in.Name")},
									},
								},
								&components.LabelInline{
									Title: "Symbol",
									Children: []components.PageInterface{
										&components.FieldText{Getter: getters.Key[string]("$in.Symbol")},
									},
								},
								&components.LabelInline{
									Title: "Active",
									Children: []components.PageInterface{
										&components.FieldCheckbox{Getter: getters.Key[bool]("$in.IsActive")},
									},
								},
							},
						},
					},
				},
			},
		}},
	}
}
