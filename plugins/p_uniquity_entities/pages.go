package p_uniquity_entities

import (
	currencies "github.com/UniquityVentures/uniquity/plugins/p_uniquity_currencies"
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

func entityFormInputs() []components.PageInterface {
	return []components.PageInterface{
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
			Error: getters.Key[error]("$error.Slug"),
			Children: []components.PageInterface{
				&components.InputText{
					Label:  "Slug",
					Name:   "Slug",
					Getter: getters.Key[string]("$in.Slug"),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.TaxID"),
			Children: []components.PageInterface{
				&components.InputText{
					Label:  "Tax ID",
					Name:   "TaxID",
					Getter: getters.Key[string]("$in.TaxID"),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.CurrencyID"),
			Children: []components.PageInterface{
				&components.InputForeignKey[currencies.Currency]{
					Name:        "CurrencyID",
					Label:       "Functional currency",
					Url:         lamu.RoutePath("currencies.CurrencySelectRoute", nil),
					Display:     getters.Format("%s — %s", getters.Any(getters.Key[string]("$in.Currency.Code")), getters.Any(getters.Key[string]("$in.Currency.Name"))),
					Placeholder: "Select currency…",
					Required:    true,
					Getter:      getters.Association[currencies.Currency, uint](getters.Key[uint]("$in.CurrencyID")),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Phone"),
			Children: []components.PageInterface{
				&components.InputPhone{
					Label:  "Phone",
					Name:   "Phone",
					Getter: getters.Key[string]("$in.Phone"),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Address"),
			Children: []components.PageInterface{
				&components.InputTextarea{
					Label:  "Address",
					Name:   "Address",
					Getter: getters.Key[string]("$in.Address"),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Mobile1"),
			Children: []components.PageInterface{
				&components.InputPhone{
					Label:  "Mobile 1",
					Name:   "Mobile1",
					Getter: getters.Key[string]("$in.Mobile1"),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Mobile2"),
			Children: []components.PageInterface{
				&components.InputPhone{
					Label:  "Mobile 2",
					Name:   "Mobile2",
					Getter: getters.Key[string]("$in.Mobile2"),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Email"),
			Children: []components.PageInterface{
				&components.InputEmail{
					Label:  "Email",
					Name:   "Email",
					Getter: getters.Key[string]("$in.Email"),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Website"),
			Children: []components.PageInterface{
				&components.InputText{
					Label:  "Website",
					Name:   "Website",
					Getter: getters.Key[string]("$in.Website"),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.LogoPath"),
			Children: []components.PageInterface{
				&components.InputText{
					Label:  "Logo path (storage key, optional)",
					Name:   "LogoPath",
					Getter: getters.Key[string]("$in.LogoPath"),
				},
			},
		},
	}
}

func pageEntries() []registry.Pair[string, components.PageInterface] {
	createName := getters.Static("entities.EntityCreateForm")
	updateName := getters.Static("entities.EntityUpdateForm")
	deleteName := getters.Static("entities.EntityDeleteForm")

	return []registry.Pair[string, components.PageInterface]{
		{Key: "entities.MainMenu", Value: &components.SidebarMenu{
			Title: getters.Static("Entities"),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("Back to Home"),
				Url:   lamu.RoutePath("dashboard.AppsPage", nil),
			},
			Children: []components.PageInterface{
				&components.SidebarMenuItem{
					Page:  components.Page{Roles: []string{"superuser"}},
					Title: getters.Static("All entities"),
					Url:   lamu.RoutePath("entities.EntityListRoute", nil),
					Icon:  "building-office",
				},
				&components.SidebarMenuItem{
					Page:  components.Page{Roles: []string{"superuser"}},
					Title: getters.Static("New entity"),
					Url:   lamu.RoutePath("entities.EntityCreateRoute", nil),
					Icon:  "plus",
				},
			},
		}},
		{Key: "entities.EntityDetailMenu", Value: &components.SidebarMenu{
			Title: getters.Format("Entity #%d", getters.Any(getters.Key[uint]("entity.ID"))),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("All entities"),
				Url:   lamu.RoutePath("entities.EntityListRoute", nil),
			},
			Children: []components.PageInterface{
				&components.SidebarMenuItem{
					Title: getters.Static("Detail"),
					Url: lamu.RoutePath("entities.EntityDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("entity.ID")),
					}),
				},
				&components.SidebarMenuItem{
					Page:  components.Page{Roles: []string{"superuser"}},
					Title: getters.Static("Edit"),
					Url: lamu.RoutePath("entities.EntityUpdateRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("entity.ID")),
					}),
				},
			},
		}},
			{Key: "entities.EntitySelectionTable", Value: &components.ShellScaffold{
				Page:    components.Page{Roles: []string{"superuser"}},
				Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "entities.MainMenu"}},
				Children: []components.PageInterface{
					&components.DataTable[Entity]{
						UID:     "entities-select-table",
						Classes: "w-full",
						Data:    getters.Key[components.ObjectList[Entity]]("entities"),
						RowAttr: getters.RowAttrSelect("EntityID",
							getters.Key[uint]("$row.ID"),
							getters.Format("%s (%s)", getters.Any(getters.Key[string]("$row.Name")), getters.Any(getters.Key[string]("$row.Currency.Code"))),
						),
						Columns: []components.TableColumn{
							{Label: "Name", Name: "Name", Children: []components.PageInterface{
								&components.FieldText{Getter: getters.Key[string]("$row.Name")},
							}},
							{Label: "Slug", Name: "Slug", Children: []components.PageInterface{
								&components.FieldText{Getter: getters.Key[string]("$row.Slug")},
							}},
						},
					},
				},
			}},
			{Key: "entities.EntityTable", Value: &components.ShellScaffold{
			Page:    components.Page{Roles: []string{"superuser"}},
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "entities.MainMenu"}},
			Children: []components.PageInterface{
				&components.DataTable[Entity]{
					UID:     "entities-table",
					Classes: "w-full",
					Data:    getters.Key[components.ObjectList[Entity]]("entities"),
					Actions: []components.PageInterface{
						&components.TableButtonCreate{
							Link: lamu.RoutePath("entities.EntityCreateRoute", nil),
							Page: components.Page{Roles: []string{"superuser"}},
						},
					},
					RowAttr: getters.RowAttrNavigate(lamu.RoutePath("entities.EntityDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("$row.ID")),
					})),
					Columns: []components.TableColumn{
						{Label: "Name", Name: "Name", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Name")},
						}},
						{Label: "Slug", Name: "Slug", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Slug")},
						}},
						{Label: "Currency", Name: "Currency", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Currency.Code")},
						}},
						{Label: "Email", Name: "Email", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Email")},
						}},
						{Label: "Phone", Name: "Phone", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Phone")},
						}},
						{Label: "Mobile 1", Name: "Mobile1", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Mobile1")},
						}},
					},
				},
			},
		}},
		{Key: "entities.EntityCreateForm", Value: &components.ShellScaffold{
			Page:    components.Page{Roles: []string{"superuser"}},
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "entities.MainMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name:      createName,
					ActionURL: lamu.RoutePath("entities.EntityCreateRoute", nil),
					Children: []components.PageInterface{
						&components.FormComponent[Entity]{
							Attr:          getters.FormBubbling(createName),
							Title:         "Create entity",
							Subtitle:      "Add a new entity",
							ChildrenInput: entityFormInputs(),
							ChildrenAction: []components.PageInterface{
								&components.ButtonSubmit{Label: "Save"},
							},
						},
					},
				},
			},
		}},
		{Key: "entities.EntityUpdateForm", Value: &components.ShellScaffold{
			Page:    components.Page{Roles: []string{"superuser"}},
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "entities.EntityDetailMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name: updateName,
					ActionURL: lamu.RoutePath("entities.EntityUpdateRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("entity.ID")),
					}),
					Children: []components.PageInterface{
						&components.FormComponent[Entity]{
							Getter:        getters.Key[Entity]("entity"),
							Attr:          getters.FormBubbling(updateName),
							Title:         "Edit entity",
							Subtitle:      "Update contact details",
							ChildrenInput: entityFormInputs(),
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
													Url: lamu.RoutePath("entities.EntityDeleteRoute", map[string]getters.Getter[any]{
														"id": getters.Any(getters.Key[uint]("entity.ID")),
													}),
													FormPostURL: lamu.RoutePath("entities.EntityDeleteRoute", map[string]getters.Getter[any]{
														"id": getters.Any(getters.Key[uint]("entity.ID")),
													}),
													ModalUID: "entity-delete-modal",
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
		{Key: "entities.EntityDeleteForm", Value: &components.Modal{
			Page: components.Page{Roles: []string{"superuser"}},
			UID:  "entity-delete-modal",
			Children: []components.PageInterface{
				&components.DeleteConfirmation{
					Title:   "Delete entity?",
					Message: "This permanently removes this entity record.",
					Attr:    getters.FormBubbling(deleteName),
				},
			},
		}},
		{Key: "entities.EntityDetail", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "entities.EntityDetailMenu"}},
			Children: []components.PageInterface{
				&components.Detail[Entity]{
					Getter: getters.Key[Entity]("entity"),
					Children: []components.PageInterface{
						&components.ContainerColumn{
							Classes: "p-4",
							Children: []components.PageInterface{
								&components.LabelInline{
									Title: "Name",
									Children: []components.PageInterface{
										&components.FieldText{Getter: getters.Key[string]("$in.Name")},
									},
								},
								&components.LabelInline{
									Title: "Slug",
									Children: []components.PageInterface{
										&components.FieldText{Getter: getters.Key[string]("$in.Slug")},
									},
								},
								&components.LabelInline{
									Title: "Tax ID",
									Children: []components.PageInterface{
										&components.FieldText{Getter: getters.Key[string]("$in.TaxID")},
									},
								},
								&components.LabelInline{
									Title: "Functional currency",
									Children: []components.PageInterface{
										&components.FieldText{Getter: getters.Format("%s — %s", getters.Any(getters.Key[string]("$in.Currency.Code")), getters.Any(getters.Key[string]("$in.Currency.Name")))},
									},
								},
								&components.LabelInline{
									Title: "Phone",
									Children: []components.PageInterface{
										&components.FieldText{Getter: getters.Key[string]("$in.Phone")},
									},
								},
								&components.LabelInline{
									Title: "Address",
									Children: []components.PageInterface{
										&components.FieldText{Getter: getters.Key[string]("$in.Address")},
									},
								},
								&components.LabelInline{
									Title: "Mobile 1",
									Children: []components.PageInterface{
										&components.FieldText{Getter: getters.Key[string]("$in.Mobile1")},
									},
								},
								&components.LabelInline{
									Title: "Mobile 2",
									Children: []components.PageInterface{
										&components.FieldText{Getter: getters.Key[string]("$in.Mobile2")},
									},
								},
								&components.LabelInline{
									Title: "Email",
									Children: []components.PageInterface{
										&components.FieldText{Getter: getters.Key[string]("$in.Email")},
									},
								},
								&components.LabelInline{
									Title: "Website",
									Children: []components.PageInterface{
										&components.FieldText{Getter: getters.Key[string]("$in.Website")},
									},
								},
								&components.LabelInline{
									Title: "Logo path",
									Children: []components.PageInterface{
										&components.FieldText{Getter: getters.Key[string]("$in.LogoPath")},
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
