package p_uniquity_accounting

import (
	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func patchAccountingMainMenuTransactions(p components.PageInterface) components.PageInterface {
	menu, ok := p.(components.MutableParentInterface)
	if !ok {
		panic("menu is not a MutableParentInterface, for some reason")
	}
	children := menu.GetChildren()
	children = append(children, &components.SidebarMenuItem{
		Title: getters.Static("Transactions"),
		Url:   lamu.RoutePath("accounting.TransactionListRoute", nil),
		Icon:  "arrows-right-left",
	})
	menu.SetChildren(children)
	return p
}

func pageEntriesTransactionPages() []registry.Pair[string, components.PageInterface] {
	createName := getters.Static("accounting.TransactionCreateForm")

	accountInput := &components.InputForeignKey[Account]{
		Name:        "AccountID",
		Label:       "Account",
		Url:         lamu.RoutePath("accounting.TransactionAccountSelectRoute", nil),
		Display:     getters.Format("%s - %s", getters.Any(getters.Key[string]("$in.Code")), getters.Any(getters.Key[string]("$in.Name"))),
		Placeholder: "Select account...",
		Required:    true,
		Getter:      getters.Association[Account, uint](getters.Key[uint]("$in.AccountID")),
	}

	return []registry.Pair[string, components.PageInterface]{
		{Key: "accounting.TransactionDetailMenu", Value: &components.SidebarMenu{
			Title: getters.Format("Transaction #%d", getters.Any(getters.Key[uint]("transaction.ID"))),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("All transactions"),
				Url:   lamu.RoutePath("accounting.TransactionListRoute", nil),
			},
			Children: []components.PageInterface{
				&components.SidebarMenuItem{
					Title: getters.Static("Detail"),
					Url: lamu.RoutePath("accounting.TransactionDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("transaction.ID")),
					}),
				},
			},
		}},
		{Key: "accounting.TransactionTable", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "accounting.MainMenu"}},
			Children: []components.PageInterface{
				&components.DataTable[Posting]{
					UID:     "accounting-transaction-table",
					Classes: "w-full",
					Data:    getters.Key[components.ObjectList[Posting]]("transactions"),
					Actions: []components.PageInterface{
						&components.TableButtonCreate{
							Link: lamu.RoutePath("accounting.TransactionCreateRoute", nil),
						},
					},
					RowAttr: getters.RowAttrNavigate(lamu.RoutePath("accounting.TransactionDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("$row.ID")),
					})),
					Columns: []components.TableColumn{
						{Label: "Amount", Name: "Amount", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Amount")},
						}},
						{Label: "Account", Name: "Account", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Account.Name")},
						}},
					},
				},
			},
		}},
		{Key: "accounting.TransactionCreateForm", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "accounting.MainMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name:      createName,
					ActionURL: lamu.RoutePath("accounting.TransactionCreateRoute", nil),
					Children: []components.PageInterface{
						&components.FormComponent[Posting]{
							Attr:     getters.FormBubbling(createName),
							Title:    "Create transaction",
							Subtitle: "Record a transaction against an account",
							ChildrenInput: []components.PageInterface{
								&components.InputText{
									Label:    "Amount",
									Name:     "Amount",
									Required: true,
									Getter:   getters.Key[string]("$in.Amount"),
								},
								accountInput,
							},
							ChildrenAction: []components.PageInterface{
								&components.ButtonSubmit{Label: "Save"},
							},
						},
					},
				},
			},
		}},
		{Key: "accounting.TransactionDetail", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "accounting.TransactionDetailMenu"}},
			Children: []components.PageInterface{
				&components.Detail[Posting]{
					Getter: getters.Key[Posting]("transaction"),
					Children: []components.PageInterface{
						&components.ContainerColumn{
							Classes: "p-4",
							Children: []components.PageInterface{
								&components.LabelInline{
									Title: "Amount",
									Children: []components.PageInterface{
										&components.FieldText{Getter: getters.Key[string]("$in.Amount")},
									},
								},
								&components.LabelInline{
									Title: "Account",
									Children: []components.PageInterface{
										&components.FieldText{Getter: getters.Format("%s - %s", getters.Any(getters.Key[string]("$in.Account.Code")), getters.Any(getters.Key[string]("$in.Account.Name")))},
									},
								},
							},
						},
					},
				},
			},
		}},
		{Key: "accounting.TransactionAccountSelectionTable", Value: &components.Modal{
			UID: "accounting-account-select-modal",
			Children: []components.PageInterface{
				&components.DataTable[Account]{
					UID:   "accounting-account-select-table",
					Title: "Select account",
					Data:  getters.Key[components.ObjectList[Account]]("accounts"),
					RowAttr: getters.RowAttrSelect("AccountID",
						getters.Key[uint]("$row.ID"),
						getters.Format("%s - %s", getters.Any(getters.Key[string]("$row.Code")), getters.Any(getters.Key[string]("$row.Name"))),
					),
					Columns: []components.TableColumn{
						{Label: "Code", Name: "Code", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Code")},
						}},
						{Label: "Name", Name: "Name", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Name")},
						}},
						{Label: "Type", Name: "AccountType", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.AccountType")},
						}},
					},
				},
			},
		}},
	}
}
