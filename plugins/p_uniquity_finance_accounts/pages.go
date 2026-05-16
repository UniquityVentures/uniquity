package p_uniquity_finance_accounts

import (
	"context"
	"fmt"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
	"gorm.io/gorm"
)

var balanceTypeChoices = getters.Static([]registry.Pair[BalanceType, string]{
	{Key: BalanceTypeCredit, Value: "Credit"},
	{Key: BalanceTypeDebit, Value: "Debit"},
})

func balanceTypeSelectGetter(ctxKey string) getters.Getter[registry.Pair[BalanceType, string]] {
	return func(ctx context.Context) (registry.Pair[BalanceType, string], error) {
		bt, err := getters.Key[BalanceType](ctxKey)(ctx)
		if err != nil {
			return registry.Pair[BalanceType, string]{}, err
		}
		if bt == "" {
			return registry.Pair[BalanceType, string]{}, nil
		}
		return registry.Pair[BalanceType, string]{Key: bt, Value: string(bt)}, nil
	}
}

// accountParentLabel returns "code — name" for the FK in list/detail rows, or "—" if none.
func accountParentLabel(rowPrefix string) getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		p, err := getters.Key[*uint](rowPrefix + ".ParentID")(ctx)
		if err != nil {
			return "", err
		}
		if p == nil || *p == 0 {
			return "—", nil
		}
		db, err := getters.DBFromContext(ctx)
		if err != nil {
			return "", err
		}
		a, err := gorm.G[Account](db).Where("id = ?", *p).First(ctx)
		if err != nil {
			return "—", nil
		}
		return fmt.Sprintf("%d — %s", a.Code, a.Name), nil
	}
}

// accountTableGroupCell renders stacked rectangles for a group/summary row or a document icon
// for a posting row (Heroicons via Icon).
func accountTableGroupCell(rowPrefix string) []components.PageInterface {
	return []components.PageInterface{
		&components.ShowIf{
			Getter: func(ctx context.Context) (any, error) {
				return getters.Key[bool](rowPrefix + ".IsGroup")(ctx)
			},
			Children: []components.PageInterface{
				&components.Icon{Name: "rectangle-stack", Classes: "heroicon-sm"},
			},
		},
		&components.ShowIf{
			Getter: func(ctx context.Context) (any, error) {
				isGroup, err := getters.Key[bool](rowPrefix + ".IsGroup")(ctx)
				if err != nil {
					return false, err
				}
				return !isGroup, nil
			},
			Children: []components.PageInterface{
				&components.Icon{Name: "document-text", Classes: "heroicon-sm opacity-70"},
			},
		},
	}
}

// accountTableNameCell renders the account name with a leading group/posting icon in a row.
func accountTableNameCell(rowPrefix string) []components.PageInterface {
	return []components.PageInterface{
		&components.ContainerRow{
			Classes: "items-center gap-2 min-w-0",
			Children: append(accountTableGroupCell(rowPrefix), &components.FieldText{
				Getter: getters.Key[string](rowPrefix + ".Name"),
			}),
		},
	}
}

// accountDetailParentHref is the detail URL for the parent account, or empty if none.
func accountDetailParentHref() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		p, err := getters.Key[*uint]("$in.ParentID")(ctx)
		if err != nil {
			return "", err
		}
		if p == nil || *p == 0 {
			return "", nil
		}
		pid := *p
		return lamu.RoutePath("finance_accounts.AccountDetailRoute", map[string]getters.Getter[any]{
			"id": func(context.Context) (any, error) { return pid, nil },
		})(ctx)
	}
}

func pluginPages() lamu.PluginFeatures[components.PageInterface] {
	var entries []registry.Pair[string, components.PageInterface]
	entries = append(entries, pageMenus()...)
	entries = append(entries, pageFilterPages()...)
	entries = append(entries, pageAccountCRUD()...)
	entries = append(entries, pageCurrencyCRUD()...)
	entries = append(entries, pageJournalCRUD()...)
	entries = append(entries, pageJournalFKSelectPages()...)
	entries = append(entries, pageJournalEntryCreatePages()...)
	entries = append(entries, pageJournalEntryDetailPages()...)
	entries = append(entries, pageAccountingPreferencesPages()...)
	return lamu.PluginFeatures[components.PageInterface]{Entries: entries}
}

func pageMenus() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "finance_accounts.MainMenu", Value: &components.SidebarMenu{
			Title: getters.Static("Finance accounts"),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("Back to Home"),
				Url:   lamu.RoutePath("dashboard.AppsPage", nil),
			},
			Children: []components.PageInterface{
				&components.SidebarMenuItem{
					Page:  components.Page{Roles: []string{"superuser"}},
					Title: getters.Static("Accounts"),
					Url:   lamu.RoutePath("finance_accounts.DefaultRoute", nil),
					Icon:  "building-library",
				},
				&components.SidebarMenuItem{
					Page:  components.Page{Roles: []string{"superuser"}},
					Title: getters.Static("Currencies"),
					Url:   lamu.RoutePath("finance_accounts.CurrencyListRoute", nil),
					Icon:  "currency-dollar",
				},
				&components.SidebarMenuItem{
					Page:  components.Page{Roles: []string{"superuser"}},
					Title: getters.Static("Journals"),
					Url:   lamu.RoutePath("finance_accounts.JournalListRoute", nil),
					Icon:  "book-open",
				},
				&components.SidebarMenuItem{
					Page:  components.Page{Roles: []string{"superuser"}},
					Title: getters.Static("Accounting preferences"),
					Url:   lamu.RoutePath("finance_accounts.AccountingPreferencesRoute", nil),
					Icon:  "adjustments-horizontal",
				},
			},
		}},
		{Key: "finance_accounts.AccountDetailMenu", Value: &components.SidebarMenu{
			Title: getters.Format("Account #%d", getters.Any(getters.Key[uint]("account.ID"))),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("All accounts"),
				Url:   lamu.RoutePath("finance_accounts.DefaultRoute", nil),
			},
			Children: []components.PageInterface{
				&components.SidebarMenuItem{
					Title: getters.Static("Detail"),
					Url: lamu.RoutePath("finance_accounts.AccountDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("account.ID")),
					}),
				},
				&components.SidebarMenuItem{
					Page:  components.Page{Roles: []string{"superuser"}},
					Title: getters.Static("Edit"),
					Url: lamu.RoutePath("finance_accounts.AccountUpdateRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("account.ID")),
					}),
				},
			},
		}},
		{Key: "finance_accounts.CurrencyDetailMenu", Value: &components.SidebarMenu{
			Title: getters.Format("Currency #%d", getters.Any(getters.Key[uint]("currency.ID"))),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("All currencies"),
				Url:   lamu.RoutePath("finance_accounts.CurrencyListRoute", nil),
			},
			Children: []components.PageInterface{
				&components.SidebarMenuItem{
					Title: getters.Static("Detail"),
					Url: lamu.RoutePath("finance_accounts.CurrencyDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("currency.ID")),
					}),
				},
				&components.SidebarMenuItem{
					Page:  components.Page{Roles: []string{"superuser"}},
					Title: getters.Static("Edit"),
					Url: lamu.RoutePath("finance_accounts.CurrencyUpdateRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("currency.ID")),
					}),
				},
			},
		}},
		{Key: "finance_accounts.JournalDetailMenu", Value: &components.SidebarMenu{
			Title: getters.Format("Journal #%d", getters.Any(getters.Key[uint]("journal.ID"))),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("All journals"),
				Url:   lamu.RoutePath("finance_accounts.JournalListRoute", nil),
			},
			Children: []components.PageInterface{
				&components.SidebarMenuItem{
					Title: getters.Static("Detail"),
					Url: lamu.RoutePath("finance_accounts.JournalDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("journal.ID")),
					}),
				},
				&components.SidebarMenuItem{
					Page:  components.Page{Roles: []string{"superuser"}},
					Title: getters.Static("Edit"),
					Url: lamu.RoutePath("finance_accounts.JournalUpdateRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("journal.ID")),
					}),
				},
			},
		}},
		{Key: "finance_accounts.JournalEntryDetailMenu", Value: &components.SidebarMenu{
			Title: getters.Format("Journal entry #%d", getters.Any(getters.Key[uint]("journalEntry.ID"))),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("Back to journal"),
				Url: lamu.RoutePath("finance_accounts.JournalDetailRoute", map[string]getters.Getter[any]{
					"id": getters.Any(getters.Key[uint]("journalEntry.JournalID")),
				}),
			},
			Children: []components.PageInterface{
				&components.SidebarMenuItem{
					Title: getters.Static("Detail"),
					Url: lamu.RoutePath("finance_accounts.JournalEntryDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("journalEntry.ID")),
					}),
				},
			},
		}},
	}
}

func pageAccountCRUD() []registry.Pair[string, components.PageInterface] {
	createName := getters.Static("finance_accounts.AccountCreateForm")
	updateName := getters.Static("finance_accounts.AccountUpdateForm")
	deleteName := getters.Static("finance_accounts.AccountDeleteForm")

	nameInput := &components.InputText{
		Name:     "Name",
		Label:    "Name",
		Required: true,
		Getter:   getters.Key[string]("$in.Name"),
	}
	codeInput := &components.InputNumber[int]{
		Name:     "Code",
		Label:    "Code",
		Required: true,
		Getter:   getters.Key[int]("$in.Code"),
	}
	groupInput := &components.InputCheckbox{
		Name:   "IsGroup",
		Label:  "Group account (summary)",
		Getter: getters.Key[bool]("$in.IsGroup"),
	}
	balanceTypeInput := &components.InputSelect[BalanceType]{
		Name:     "BalanceType",
		Label:    "Balance type",
		Required: true,
		Choices:  balanceTypeChoices,
		Getter:   balanceTypeSelectGetter("$in.BalanceType"),
	}
	parentPicker := &components.InputForeignKey[Account]{
		Name:        "ParentID",
		Label:       "Parent account",
		Url:         lamu.RoutePath("finance_accounts.AccountSelectRoute", nil),
		Display:     getters.Format("%d — %s", getters.Any(getters.Key[uint]("$in.ID")), getters.Any(getters.Key[string]("$in.Name"))),
		Placeholder: "Optional parent…",
		Required:    false,
		Getter:      getters.Association[Account, *uint](getters.Key[*uint]("$in.ParentID")),
	}

	return []registry.Pair[string, components.PageInterface]{
		{Key: "finance_accounts.AccountTable", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_accounts.MainMenu"}},
			Children: []components.PageInterface{
				&components.DataTable[Account]{
					UID:     "finance-accounts-table",
					Title:   "Top-level accounts",
					Classes: "w-full",
					Data:    getters.Key[components.ObjectList[Account]]("accounts"),
					Actions: []components.PageInterface{
						&components.TableButtonFilter{Child: lamu.DynamicPage{Name: "finance_accounts.AccountFilter"}},
						&components.TableButtonCreate{
							Link: lamu.RoutePath("finance_accounts.AccountCreateRoute", nil),
							Page: components.Page{Roles: []string{"superuser"}},
						},
					},
					RowAttr: getters.RowAttrNavigate(lamu.RoutePath("finance_accounts.AccountDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("$row.ID")),
					})),
					Columns: []components.TableColumn{
						{Label: "Name", Name: "Name", Children: accountTableNameCell("$row")},
						{Label: "Code", Name: "Code", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Format("%d", getters.Any(getters.Key[int]("$row.Code")))},
						}},
						{Label: "Balance type", Name: "BalanceType", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Format("%s", getters.Any(getters.Key[BalanceType]("$row.BalanceType")))},
						}},
					},
				},
			},
		}},
		{Key: "finance_accounts.AccountCreateForm", Value: &components.ShellScaffold{
			Page:    components.Page{Roles: []string{"superuser"}},
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_accounts.MainMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name:      createName,
					ActionURL: lamu.RoutePath("finance_accounts.AccountCreateRoute", nil),
					Children: []components.PageInterface{
						&components.FormComponent[Account]{
							Attr:     getters.FormBubbling(createName),
							Title:    "Create account",
							Subtitle: "Chart of accounts entry",
							ChildrenInput: []components.PageInterface{
								nameInput,
								codeInput,
								groupInput,
								balanceTypeInput,
								parentPicker,
							},
							ChildrenAction: []components.PageInterface{
								&components.ButtonSubmit{Label: "Save"},
							},
						},
					},
				},
			},
		}},
		{Key: "finance_accounts.AccountUpdateForm", Value: &components.ShellScaffold{
			Page:    components.Page{Roles: []string{"superuser"}},
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_accounts.AccountDetailMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name: updateName,
					ActionURL: lamu.RoutePath("finance_accounts.AccountUpdateRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("account.ID")),
					}),
					Children: []components.PageInterface{
						&components.FormComponent[Account]{
							Getter:   getters.Key[Account]("account"),
							Attr:     getters.FormBubbling(updateName),
							Title:    "Edit account",
							Subtitle: "Update account fields",
							ChildrenInput: []components.PageInterface{
								nameInput,
								codeInput,
								groupInput,
								balanceTypeInput,
								parentPicker,
							},
							ChildrenAction: []components.PageInterface{
								&components.ContainerRow{
									Classes: "flex flex-wrap justify-between gap-2 mt-2 items-center",
									Children: []components.PageInterface{
										&components.ContainerRow{
											Classes: "flex justify-end gap-2",
											Children: []components.PageInterface{
												&components.ButtonSubmit{Label: "Update"},
												&components.ButtonModalForm{
													Page:  components.Page{Roles: []string{"superuser"}},
													Label: "Delete",
													Icon:  "trash",
													Name:  deleteName,
													Url: lamu.RoutePath("finance_accounts.AccountDeleteRoute", map[string]getters.Getter[any]{
														"id": getters.Any(getters.Key[uint]("account.ID")),
													}),
													FormPostURL: lamu.RoutePath("finance_accounts.AccountDeleteRoute", map[string]getters.Getter[any]{
														"id": getters.Any(getters.Key[uint]("account.ID")),
													}),
													ModalUID: "finance-account-delete-modal",
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
		{Key: "finance_accounts.AccountDeleteForm", Value: &components.Modal{
			Page: components.Page{Roles: []string{"superuser"}},
			UID:  "finance-account-delete-modal",
			Children: []components.PageInterface{
				&components.DeleteConfirmation{
					Title:   "Delete account?",
					Message: "This removes the account. Child accounts that referenced it will have parent cleared.",
					Attr:    getters.FormBubbling(getters.Key[string]("$get.name")),
				},
			},
		}},
		{Key: "finance_accounts.AccountDetail", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_accounts.AccountDetailMenu"}},
			Children: []components.PageInterface{
				&components.Detail[Account]{
					Getter: getters.Key[Account]("account"),
					Children: []components.PageInterface{
						&components.ContainerColumn{
							Classes: "p-4",
							Children: []components.PageInterface{
								&components.FieldTitle{Getter: getters.Key[string]("$in.Name")},
								&components.FieldSubtitle{Getter: getters.Format("Code %d", getters.Any(getters.Key[int]("$in.Code")))},
								&components.ShowIf{
									Getter: func(ctx context.Context) (any, error) {
										return getters.Key[bool]("$in.IsGroup")(ctx)
									},
									Children: []components.PageInterface{
										&components.FieldSubtitle{Getter: getters.Static("Group")},
									},
								},
								&components.LabelInline{
									Title:   "Balance type",
									Classes: "mt-2",
									Children: []components.PageInterface{
										&components.FieldText{Getter: getters.Format("%s", getters.Any(getters.Key[BalanceType]("$in.BalanceType")))},
									},
								},
								&components.LabelInline{
									Title: "Parent",
									Children: []components.PageInterface{
										&components.FieldLink{
											Href:  accountDetailParentHref(),
											Label: accountParentLabel("$in"),
										},
									},
								},
							},
						},
					},
				},
				&components.ShowIf{
					Getter: func(ctx context.Context) (any, error) {
						return getters.Key[bool]("account.IsGroup")(ctx)
					},
					Children: []components.PageInterface{
						&components.DataTable[Account]{
							UID:      "finance-account-children-table",
							Title:    "Child accounts",
							Subtitle: "Direct children of this account",
							Classes:  "w-full",
							Data:     getters.Key[components.ObjectList[Account]](accountChildrenContextKey),
							RowAttr: getters.RowAttrNavigate(lamu.RoutePath("finance_accounts.AccountDetailRoute", map[string]getters.Getter[any]{
								"id": getters.Any(getters.Key[uint]("$row.ID")),
							})),
							Columns: []components.TableColumn{
								{Label: "Name", Name: "Name", Children: accountTableNameCell("$row")},
								{Label: "Code", Name: "Code", Children: []components.PageInterface{
									&components.FieldText{Getter: getters.Format("%d", getters.Any(getters.Key[int]("$row.Code")))},
								}},
								{Label: "Balance type", Name: "BalanceType", Children: []components.PageInterface{
									&components.FieldText{Getter: getters.Format("%s", getters.Any(getters.Key[BalanceType]("$row.BalanceType")))},
								}},
							},
						},
					},
				},
			},
		}},
		{Key: "finance_accounts.AccountSelectionTable", Value: &components.Modal{
			UID: "finance-account-selection-modal",
			Children: []components.PageInterface{
				&components.DataTable[Account]{
					UID:   "finance-account-selection-table",
					Title: "Select account",
					Data:  getters.Key[components.ObjectList[Account]]("accounts"),
					Actions: []components.PageInterface{
						&components.TableButtonFilter{Child: lamu.DynamicPage{Name: "finance_accounts.AccountSelectionFilter"}},
					},
					RowAttr: accountSelectionTableRowAttr(
						getters.IfOrElse(getters.Key[string]("$get.target_input"), getters.Static("ParentID")),
						getters.Key[uint]("$row.ID"),
						getters.Format("%s (#%d)",
							getters.Any(getters.Key[string]("$row.Name")),
							getters.Any(getters.Key[uint]("$row.ID")),
						),
						getters.Key[bool]("$row.IsGroup"),
					),
					Columns: []components.TableColumn{
						{Label: "Code", Name: "Code", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Format("%d", getters.Any(getters.Key[int]("$row.Code")))},
						}},
						{Label: "Name", Name: "Name", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Name")},
						}},
						{Label: "Type", Name: "BalanceType", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Format("%s", getters.Any(getters.Key[BalanceType]("$row.BalanceType")))},
						}},
					},
				},
			},
		}},
	}
}
