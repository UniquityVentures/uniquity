package p_uniquity_accounting

import (
	"context"
	"slices"

	currencies "github.com/UniquityVentures/uniquity/plugins/p_uniquity_currencies"
	ent "github.com/UniquityVentures/uniquity/plugins/p_uniquity_entities"
	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

const accountingMainMenuAccountsKey = "accounting.MainMenu.item.accounts"

func accountingSidebarHasChildKey(children []components.PageInterface, key string) bool {
	for _, c := range children {
		if c.GetKey() == key {
			return true
		}
	}
	return false
}

func patchAccountingMainMenuAccounts(p components.PageInterface) components.PageInterface {
	m, ok := p.(*components.SidebarMenu)
	if !ok {
		panic("accounting.MainMenu patch expected *components.SidebarMenu")
	}
	children := slices.Clone(m.Children)
	if !accountingSidebarHasChildKey(children, accountingMainMenuAccountsKey) {
		children = append(children, &components.SidebarMenuItem{
			Page:  components.Page{Key: accountingMainMenuAccountsKey},
			Title: getters.Static("Accounts"),
			Url:   lamu.RoutePath("accounting.AccountListRoute", nil),
			Icon:  "arrows-right-left",
		})
	}
	return &components.SidebarMenu{
		Page:     m.Page,
		Title:    m.Title,
		Back:     m.Back,
		Children: children,
	}
}

func accountTypePairGetter(g getters.Getter[string]) getters.Getter[registry.Pair[string, string]] {
	labels := make(map[string]string)
	for _, p := range AccountTypeChoices() {
		labels[p.Key] = p.Value
	}
	return func(ctx context.Context) (registry.Pair[string, string], error) {
		k, err := g(ctx)
		if err != nil {
			return registry.Pair[string, string]{}, err
		}
		lab := labels[k]
		if lab == "" {
			lab = k
		}
		return registry.Pair[string, string]{Key: k, Value: lab}, nil
	}
}

func accountFormInputs() []components.PageInterface {
	return []components.PageInterface{
		&components.ContainerError{
			Error: getters.Key[error]("$error.EntityID"),
			Children: []components.PageInterface{
				&components.InputForeignKey[ent.Entity]{
					Name:        "EntityID",
					Label:       "Entity",
					Url:         lamu.RoutePath("entities.EntitySelectRoute", nil),
					Display:     getters.Key[string]("$in.Entity.Name"),
					Placeholder: "Select entity…",
					Required:    true,
					Getter:      getters.Association[ent.Entity, uint](getters.Key[uint]("$in.EntityID")),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Code"),
			Children: []components.PageInterface{
				&components.InputText{
					Label:  "Code",
					Name:   "Code",
					Getter: getters.Key[string]("$in.Code"),
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
			Error: getters.Key[error]("$error.AccountType"),
			Children: []components.PageInterface{
				&components.InputSelect[string]{
					Label:    "Account type",
					Name:     "AccountType",
					Required: true,
					Choices:  getters.Static(AccountTypeChoices()),
					Getter:   accountTypePairGetter(getters.Key[string]("$in.AccountType")),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.CurrencyID"),
			Children: []components.PageInterface{
				&components.InputForeignKey[currencies.Currency]{
					Name:        "CurrencyID",
					Label:       "Account currency (optional)",
					Url:         lamu.RoutePath("currencies.CurrencySelectRoute", nil),
					Display:     getters.Format("%s — %s", getters.Any(getters.Key[string]("$in.Currency.Code")), getters.Any(getters.Key[string]("$in.Currency.Name"))),
					Placeholder: "Optional…",
					Required:    false,
					Getter:      getters.Association[currencies.Currency, *uint](getters.Key[*uint]("$in.CurrencyID")),
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
		&components.ContainerError{
			Error: getters.Key[error]("$error.IsReconcilable"),
			Children: []components.PageInterface{
				&components.InputCheckbox{
					Label:    "Allow reconciliation",
					Name:     "IsReconcilable",
					Required: true,
					Getter:   getters.Key[bool]("$in.IsReconcilable"),
				},
			},
		},
	}
}

func pageEntriesAccountPages() []registry.Pair[string, components.PageInterface] {
	createName := getters.Static("accounting.AccountCreateForm")
	updateName := getters.Static("accounting.AccountUpdateForm")

	return []registry.Pair[string, components.PageInterface]{
		{Key: "accounting.AccountDetailMenu", Value: &components.SidebarMenu{
			Title: getters.Format("%s — %s", getters.Any(getters.Key[string]("account.Code")), getters.Any(getters.Key[string]("account.Name"))),
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
				&components.SidebarMenuItem{
					Title: getters.Static("Edit"),
					Url: lamu.RoutePath("accounting.AccountUpdateRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("account.ID")),
					}),
				},
				&components.SidebarMenuItem{
					Title: getters.Static("Transfer"),
					Url: lamu.RoutePath("accounting.AccountTransferRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("account.ID")),
					}),
					Icon: "arrows-right-left",
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
						{Label: "Entity", Name: "Entity", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Entity.Name")},
						}},
						{Label: "Code", Name: "Code", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Code")},
						}},
						{Label: "Name", Name: "Name", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Name")},
						}},
						{Label: "Type", Name: "AccountType", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.AccountType")},
						}},
						{Label: "Active", Name: "IsActive", Children: []components.PageInterface{
							&components.FieldCheckbox{Getter: getters.Key[bool]("$row.IsActive")},
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
							Attr:          getters.FormBubbling(createName),
							Title:         "Create account",
							Subtitle:      "Add a chart-of-accounts entry",
							ChildrenInput: accountFormInputs(),
							ChildrenAction: []components.PageInterface{
								&components.ButtonSubmit{Label: "Save"},
							},
						},
					},
				},
			},
		}},
		{Key: "accounting.AccountUpdateForm", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "accounting.AccountDetailMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name: updateName,
					ActionURL: lamu.RoutePath("accounting.AccountUpdateRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("account.ID")),
					}),
					Children: []components.PageInterface{
						&components.FormComponent[Account]{
							Getter:        getters.Key[Account]("account"),
							Attr:          getters.FormBubbling(updateName),
							Title:         "Edit account",
							Subtitle:      "Update COA details",
							ChildrenInput: accountFormInputs(),
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
									Title: "Entity",
									Children: []components.PageInterface{
										&components.FieldText{Getter: getters.Key[string]("$in.Entity.Name")},
									},
								},
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
									Title: "Account type",
									Children: []components.PageInterface{
										&components.FieldText{Getter: getters.Key[string]("$in.AccountType")},
									},
								},
								&components.LabelInline{
									Title: "Currency",
									Children: []components.PageInterface{
										&components.FieldText{Getter: getters.Format("%s — %s", getters.Any(getters.Key[string]("$in.Currency.Code")), getters.Any(getters.Key[string]("$in.Currency.Name")))},
									},
								},
								&components.LabelInline{
									Title: "Active",
									Children: []components.PageInterface{
										&components.FieldCheckbox{Getter: getters.Key[bool]("$in.IsActive")},
									},
								},
								&components.LabelInline{
									Title: "Reconcilable",
									Children: []components.PageInterface{
										&components.FieldCheckbox{Getter: getters.Key[bool]("$in.IsReconcilable")},
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
