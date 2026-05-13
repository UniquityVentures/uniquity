package p_uniquity_accounting

import (
	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func patchAccountingMainMenuAccounts(p components.PageInterface) components.PageInterface {
	menu, ok := p.(components.MutableParentInterface)
	if !ok {
		panic("menu is not a MutableParentInterface, for some reason")
	}
	children := menu.GetChildren()
	children = append(children, &components.SidebarMenuItem{
		Title: getters.Static("Accounts"),
		Url:   lamu.RoutePath("accounting.AccountListRoute", nil),
		Icon:  "arrows-right-left",
	})
	menu.SetChildren(children)
	return p
}

func pageEntriesAccountPages() []registry.Pair[string, components.PageInterface] {
	createName := getters.Static("accounting.AccountCreateForm")

	return []registry.Pair[string, components.PageInterface]{
		{Key: "accounting.AccountDetailMenu", Value: &components.SidebarMenu{
			Title: getters.Format("Account #%d", getters.Any(getters.Key[uint]("account.ID"))),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("All accounts"),
				Url:   lamu.RoutePath("accounting.AccountListRoute", nil),
			},
			Children: []components.PageInterface{
				&components.SidebarMenuItem{
					Title: getters.Static("Detail"),
					Url: lamu.RoutePath("accounting.AccountDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("account.ID")),
					}),
				},
			},
		}},
		{Key: "accounting.AccountTable", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "accounting.MainMenu"}},
			Children: []components.PageInterface{
				&components.DataTable[Account]{
					UID:     "accounting-account-table",
					Classes: "w-full",
					Data:    getters.Key[components.ObjectList[Account]]("accounts"),
					Actions: []components.PageInterface{
						&components.TableButtonCreate{
							Link: lamu.RoutePath("accounting.AccountCreateRoute", nil),
						},
					},
					RowAttr: getters.RowAttrNavigate(lamu.RoutePath("accounting.AccountDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("$row.ID")),
					})),
					Columns: []components.TableColumn{
						{Label: "Name", Name: "Name", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Name")},
						}},
						{Label: "Code", Name: "Code", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Format("%d", getters.Any(getters.Key[int]("$row.Code")))},
						}},
						{Label: "Type", Name: "AccountType", Children: []components.PageInterface{
							&components.FieldCheckbox{Getter: getters.Key[bool]("$row.IsAsset")},
						}},
					},
				},
			},
		}},
		{Key: "accounting.AccountCreateForm", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "accounting.MainMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name:      createName,
					ActionURL: lamu.RoutePath("accounting.AccountCreateRoute", nil),
					Children: []components.PageInterface{
						&components.FormComponent[Account]{
							Attr:     getters.FormBubbling(createName),
							Title:    "Create account",
							Subtitle: "Add a new account",
							ChildrenInput: []components.PageInterface{
								&components.InputText{
									Label:    "Code",
									Name:     "Code",
									Required: true,
									Getter:   getters.Key[string]("$in.Code"),
								},
								&components.InputText{
									Label:    "Name",
									Name:     "Name",
									Required: true,
									Getter:   getters.Key[string]("$in.Name"),
								},
								&components.InputCheckbox{
									Label:    "Is Asset",
									Name:     "IsAsset",
									Required: true,
									Getter:   getters.Key[bool]("$in.IsAsset"),
								},
							},
							ChildrenAction: []components.PageInterface{
								&components.ButtonSubmit{Label: "Save"},
							},
						},
					},
				},
			},
		}},
		{Key: "accounting.AccountDetail", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "accounting.AccountDetailMenu"}},
			Children: []components.PageInterface{
				&components.Detail[Account]{
					Getter: getters.Key[Account]("account"),
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
									Title: "Code",
									Children: []components.PageInterface{
										&components.FieldText{Getter: getters.Format("%d", getters.Any(getters.Key[int]("$in.Code")))},
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
