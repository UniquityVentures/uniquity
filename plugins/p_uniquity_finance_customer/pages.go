package p_uniquity_finance_customer

import (
	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

const financeAccountsMainMenuCustomersLinkKey = "finance_customers.FinanceAccountsMainMenuLink"

func patchFinanceAccountsMainMenuForCustomers(page components.PageInterface) components.PageInterface {
	menu, ok := page.(*components.SidebarMenu)
	if !ok {
		panic("p_uniquity_finance_customer: finance_accounts.MainMenu must be *components.SidebarMenu")
	}
	for _, ch := range menu.Children {
		if item, ok := ch.(*components.SidebarMenuItem); ok && item.GetKey() == financeAccountsMainMenuCustomersLinkKey {
			return menu
		}
	}
	newChildren := append([]components.PageInterface{}, menu.Children...)
	newChildren = append(newChildren, &components.SidebarMenuItem{
		Page:  components.Page{Key: financeAccountsMainMenuCustomersLinkKey, Roles: []string{"superuser"}},
		Title: getters.Static("Customers"),
		Url:   lamu.RoutePath("finance_customers.DefaultRoute", nil),
		Icon:  "building-storefront",
	})
	cloned := *menu
	cloned.Children = newChildren
	return &cloned
}

func pluginPages() lamu.PluginFeatures[components.PageInterface] {
	e := pageEntriesCustomerMenus()
	e = append(e, pageEntriesCustomerPages()...)
	e = append(e, pageEntriesCustomerFkSelectPages()...)
	return lamu.PluginFeatures[components.PageInterface]{
		Entries: e,
		Patches: []registry.Pair[string, func(components.PageInterface) components.PageInterface]{
			{Key: "finance_accounts.MainMenu", Value: patchFinanceAccountsMainMenuForCustomers},
		},
	}
}

func pageEntriesCustomerMenus() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "finance_customers.CustomerDetailMenu", Value: &components.SidebarMenu{
			Title: getters.Format("%s", getters.Any(getters.Key[string]("customer.Name"))),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("All customers"),
				Url:   lamu.RoutePath("finance_customers.DefaultRoute", nil),
			},
			Children: []components.PageInterface{
				&components.SidebarMenuItem{
					Title: getters.Static("Detail"),
					Url: lamu.RoutePath("finance_customers.CustomerDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("customer.ID")),
					}),
				},
				&components.SidebarMenuItem{
					Page:  components.Page{Roles: []string{"superuser"}},
					Title: getters.Static("Edit"),
					Url: lamu.RoutePath("finance_customers.CustomerUpdateRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("customer.ID")),
					}),
				},
			},
		}},
	}
}

func customerCreateFormInputs() []components.PageInterface {
	return []components.PageInterface{
		&components.ContainerError{
			Error: getters.Key[error]("$error.Name"),
			Children: []components.PageInterface{
				&components.InputText{Name: "Name", Label: "Name", Required: true},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Address"),
			Children: []components.PageInterface{
				&components.InputTextarea{Name: "Address", Label: "Address", Rows: 4},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.GSTIN"),
			Children: []components.PageInterface{
				&components.InputText{Name: "GSTIN", Label: "GSTIN"},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.PAN"),
			Children: []components.PageInterface{
				&components.InputText{Name: "PAN", Label: "PAN"},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Phone"),
			Children: []components.PageInterface{
				&components.InputPhone{Name: "Phone", Label: "Phone"},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Email"),
			Children: []components.PageInterface{
				&components.InputEmail{Name: "Email", Label: "Email"},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Website"),
			Children: []components.PageInterface{
				&components.InputText{Name: "Website", Label: "Website"},
			},
		},
	}
}

func customerUpdateFormInputs() []components.PageInterface {
	return []components.PageInterface{
		&components.ContainerError{
			Error: getters.Key[error]("$error.Name"),
			Children: []components.PageInterface{
				&components.InputText{Name: "Name", Label: "Name", Required: true, Getter: getters.Key[string]("$in.Name")},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Address"),
			Children: []components.PageInterface{
				&components.InputTextarea{Name: "Address", Label: "Address", Getter: getters.Key[string]("$in.Address"), Rows: 4},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.GSTIN"),
			Children: []components.PageInterface{
				&components.InputText{Name: "GSTIN", Label: "GSTIN", Getter: getters.Key[string]("$in.GSTIN")},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.PAN"),
			Children: []components.PageInterface{
				&components.InputText{Name: "PAN", Label: "PAN", Getter: getters.Key[string]("$in.PAN")},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Phone"),
			Children: []components.PageInterface{
				&components.InputPhone{Name: "Phone", Label: "Phone", Getter: getters.Key[string]("$in.Phone")},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Email"),
			Children: []components.PageInterface{
				&components.InputEmail{Name: "Email", Label: "Email", Getter: getters.Key[string]("$in.Email")},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Website"),
			Children: []components.PageInterface{
				&components.InputText{Name: "Website", Label: "Website", Getter: getters.Key[string]("$in.Website")},
			},
		},
	}
}

func pageEntriesCustomerPages() []registry.Pair[string, components.PageInterface] {
	createName := getters.Static("finance_customers.CustomerCreateForm")
	updateName := getters.Static("finance_customers.CustomerUpdateForm")
	deleteName := getters.Static("finance_customers.CustomerDeleteForm")

	return []registry.Pair[string, components.PageInterface]{
		{Key: "finance_customers.CustomerTable", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_accounts.MainMenu"}},
			Children: []components.PageInterface{
				&components.DataTable[Customer]{
					UID:     "finance-customer-table",
					Classes: "w-full",
					Data:    getters.Key[components.ObjectList[Customer]]("customers"),
					Actions: []components.PageInterface{
						&components.TableButtonCreate{
							Link: lamu.RoutePath("finance_customers.CustomerCreateRoute", nil),
							Page: components.Page{Roles: []string{"superuser"}},
						},
					},
					RowAttr: getters.RowAttrNavigate(lamu.RoutePath("finance_customers.CustomerDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("$row.ID")),
					})),
					Columns: []components.TableColumn{
						{Label: "Name", Name: "Name", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Name")},
						}},
						{Label: "Email", Name: "Email", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Email")},
						}},
						{Label: "Phone", Name: "Phone", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Phone")},
						}},
						{Label: "GSTIN", Name: "GSTIN", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.GSTIN")},
						}},
					},
				},
			},
		}},
		{Key: "finance_customers.CustomerCreateForm", Value: &components.ShellScaffold{
			Page:    components.Page{Roles: []string{"superuser"}},
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_accounts.MainMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name:      createName,
					ActionURL: lamu.RoutePath("finance_customers.CustomerCreateRoute", nil),
					Children: []components.PageInterface{
						&components.FormComponent[Customer]{
							Attr:          getters.FormBubbling(createName),
							Title:         "Create customer",
							Subtitle:      "Contact and tax details",
							ChildrenInput: customerCreateFormInputs(),
							ChildrenAction: []components.PageInterface{
								&components.ButtonSubmit{Label: "Save"},
							},
						},
					},
				},
			},
		}},
		{Key: "finance_customers.CustomerUpdateForm", Value: &components.ShellScaffold{
			Page:    components.Page{Roles: []string{"superuser"}},
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_customers.CustomerDetailMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name: updateName,
					ActionURL: lamu.RoutePath("finance_customers.CustomerUpdateRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("customer.ID")),
					}),
					Children: []components.PageInterface{
						&components.FormComponent[Customer]{
							Getter:        getters.Key[Customer]("customer"),
							Attr:          getters.FormBubbling(updateName),
							Title:         "Edit customer",
							Subtitle:      "Update contact and tax details",
							ChildrenInput: customerUpdateFormInputs(),
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
													Url:         lamu.RoutePath("finance_customers.CustomerDeleteRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("customer.ID"))}),
													FormPostURL: lamu.RoutePath("finance_customers.CustomerDeleteRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("customer.ID"))}),
													ModalUID:    "finance-customer-delete-modal",
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
		{Key: "finance_customers.CustomerDeleteForm", Value: &components.Modal{
			Page: components.Page{Roles: []string{"superuser"}},
			UID:  "finance-customer-delete-modal",
			Children: []components.PageInterface{
				&components.DeleteConfirmation{
					Title:   "Delete customer?",
					Message: "This permanently removes the customer record.",
					Attr:    getters.FormBubbling(getters.Key[string]("$get.name")),
				},
			},
		}},
		{Key: "finance_customers.CustomerDetail", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_customers.CustomerDetailMenu"}},
			Children: []components.PageInterface{
				&components.Detail[Customer]{
					Getter: getters.Key[Customer]("customer"),
					Children: []components.PageInterface{
						&components.ContainerColumn{
							Classes: "p-4",
							Children: []components.PageInterface{
								&components.LabelInline{Title: "Name", Children: []components.PageInterface{
									&components.FieldText{Getter: getters.Key[string]("$in.Name")},
								}},
								&components.LabelInline{Title: "Address", Children: []components.PageInterface{
									&components.FieldText{Getter: getters.Key[string]("$in.Address")},
								}},
								&components.LabelInline{Title: "GSTIN", Children: []components.PageInterface{
									&components.FieldText{Getter: getters.Key[string]("$in.GSTIN")},
								}},
								&components.LabelInline{Title: "PAN", Children: []components.PageInterface{
									&components.FieldText{Getter: getters.Key[string]("$in.PAN")},
								}},
								&components.LabelInline{Title: "Phone", Children: []components.PageInterface{
									&components.FieldText{Getter: getters.Key[string]("$in.Phone")},
								}},
								&components.LabelInline{Title: "Email", Children: []components.PageInterface{
									&components.FieldText{Getter: getters.Key[string]("$in.Email")},
								}},
								&components.LabelInline{Title: "Website", Children: []components.PageInterface{
									&components.FieldText{Getter: getters.Key[string]("$in.Website")},
								}},
							},
						},
					},
				},
			},
		}},
	}
}
