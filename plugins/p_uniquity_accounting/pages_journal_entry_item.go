package p_uniquity_accounting

import (
	"slices"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/fields"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

const accountingMainMenuJournalEntryItemsKey = "accounting.MainMenu.item.journal_entry_items"

func patchAccountingMainMenuJournalEntryItems(p components.PageInterface) components.PageInterface {
	m, ok := p.(*components.SidebarMenu)
	if !ok {
		panic("accounting.MainMenu patch expected *components.SidebarMenu")
	}
	children := slices.Clone(m.Children)
	if !accountingSidebarHasChildKey(children, accountingMainMenuJournalEntryItemsKey) {
		children = append(children, &components.SidebarMenuItem{
			Page:  components.Page{Key: accountingMainMenuJournalEntryItemsKey},
			Title: getters.Static("Journal entry items"),
			Url:   lamu.RoutePath("accounting.JournalEntryItemListRoute", nil),
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

func pageEntriesJournalEntryItemPages() []registry.Pair[string, components.PageInterface] {
	createName := getters.Static("accounting.JournalEntryItemCreateForm")

	accountInput := &components.InputForeignKey[Account]{
		Name:        "AccountID",
		Label:       "Account",
		Url:         lamu.RoutePath("accounting.JournalEntryItemAccountSelectRoute", nil),
		Display:     getters.Format("%s — %s", getters.Any(getters.Key[string]("$in.Account.Code")), getters.Any(getters.Key[string]("$in.Account.Name"))),
		Placeholder: "Select account...",
		Required:    true,
		Getter:      getters.Association[Account, uint](getters.Key[uint]("$in.AccountID")),
	}

	journalEntryInput := &components.InputForeignKey[JournalEntry]{
		Name:        "JournalEntryID",
		Label:       "Journal entry",
		Url:         lamu.RoutePath("accounting.JournalEntryItemJournalEntrySelectRoute", nil),
		Display:     getters.Format("#%d — %s", getters.Any(getters.Key[uint]("$in.JournalEntry.ID")), getters.Any(getters.Key[string]("$in.JournalEntry.Journal.Name"))),
		Placeholder: "Select journal entry...",
		Required:    true,
		Getter:      getters.Association[JournalEntry, uint](getters.Key[uint]("$in.JournalEntryID")),
	}

	return []registry.Pair[string, components.PageInterface]{
		{Key: "accounting.JournalEntryItemDetailMenu", Value: &components.SidebarMenu{
			Title: getters.Format("Journal entry item #%d", getters.Any(getters.Key[uint]("journalEntryItem.ID"))),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("All journal entry items"),
				Url:   lamu.RoutePath("accounting.JournalEntryItemListRoute", nil),
			},
			Children: []components.PageInterface{
				&components.SidebarMenuItem{
					Title: getters.Static("Detail"),
					Url: lamu.RoutePath("accounting.JournalEntryItemDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("journalEntryItem.ID")),
					}),
				},
			},
		}},
		{Key: "accounting.JournalEntryItemTable", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "accounting.MainMenu"}},
			Children: []components.PageInterface{
				&components.DataTable[JournalEntryItem]{
					UID:     "accounting-journal-entry-item-table",
					Classes: "w-full",
					Data:    getters.Key[components.ObjectList[JournalEntryItem]]("journalEntryItems"),
					Actions: []components.PageInterface{
						&components.TableButtonCreate{
							Link: lamu.RoutePath("accounting.JournalEntryItemCreateRoute", nil),
						},
					},
					RowAttr: getters.RowAttrNavigate(lamu.RoutePath("accounting.JournalEntryItemDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("$row.ID")),
					})),
					Columns: []components.TableColumn{
						{Label: "Amount", Name: "Amount", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Format("%v", getters.Any(getters.Key[fields.DecimalSix]("$row.Amount")))},
						}},
						{Label: "Journal", Name: "Journal", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.JournalEntry.Journal.Name")},
						}},
						{Label: "Account", Name: "Account", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Account.Name")},
						}},
					},
				},
			},
		}},
		{Key: "accounting.JournalEntryItemCreateForm", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "accounting.MainMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name:      createName,
					ActionURL: lamu.RoutePath("accounting.JournalEntryItemCreateRoute", nil),
					Children: []components.PageInterface{
						&components.FormComponent[JournalEntryItem]{
							Attr:     getters.FormBubbling(createName),
							Title:    "Create journal entry item",
							Subtitle: "Add a line to an existing journal entry header",
							ChildrenInput: []components.PageInterface{
								&components.InputPointsDecimal{
									Label:    "Amount",
									Name:     "Amount",
									Required: true,
									Getter:   getters.Key[fields.DecimalSix]("$in.Amount"),
								},
								journalEntryInput,
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
		{Key: "accounting.JournalEntryItemDetail", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "accounting.JournalEntryItemDetailMenu"}},
			Children: []components.PageInterface{
				&components.Detail[JournalEntryItem]{
					Getter: getters.Key[JournalEntryItem]("journalEntryItem"),
					Children: []components.PageInterface{
						&components.ContainerColumn{
							Classes: "p-4",
							Children: []components.PageInterface{
								&components.LabelInline{
									Title: "Amount",
									Children: []components.PageInterface{
										&components.FieldText{Getter: getters.Format("%v", getters.Any(getters.Key[fields.DecimalSix]("$in.Amount")))},
									},
								},
								&components.LabelInline{
									Title: "Journal",
									Children: []components.PageInterface{
										&components.FieldText{Getter: getters.Key[string]("$in.JournalEntry.Journal.Name")},
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
		{Key: "accounting.JournalEntryItemAccountSelectionTable", Value: &components.Modal{
			UID: "accounting-journal-entry-item-account-select-modal",
			Children: []components.PageInterface{
				&components.DataTable[Account]{
					UID:   "accounting-journal-entry-item-account-select-table",
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
		{Key: "accounting.JournalEntryItemJournalEntrySelectionTable", Value: &components.Modal{
			UID: "accounting-journal-entry-item-journal-entry-select-modal",
			Children: []components.PageInterface{
				&components.DataTable[JournalEntry]{
					UID:   "accounting-journal-entry-item-journal-entry-select-table",
					Title: "Select journal entry",
					Data:  getters.Key[components.ObjectList[JournalEntry]]("journalEntries"),
					RowAttr: getters.RowAttrSelect("JournalEntryID",
						getters.Key[uint]("$row.ID"),
						getters.Format("#%d — %s", getters.Any(getters.Key[uint]("$row.ID")), getters.Any(getters.Key[string]("$row.Journal.Name"))),
					),
					Columns: []components.TableColumn{
						{Label: "ID", Name: "ID", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Format("%d", getters.Any(getters.Key[uint]("$row.ID")))},
						}},
						{Label: "Journal", Name: "Journal", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Journal.Name")},
						}},
					},
				},
			},
		}},
	}
}
