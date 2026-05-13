package p_uniquity_accounting

import (
	"slices"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/fields"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

const (
	accountingMainMenuJournalsKey       = "accounting.MainMenu.item.journals"
	accountingMainMenuJournalEntriesKey = "accounting.MainMenu.item.journal_entries"
)

func patchAccountingMainMenuJournals(p components.PageInterface) components.PageInterface {
	m, ok := p.(*components.SidebarMenu)
	if !ok {
		panic("accounting.MainMenu patch expected *components.SidebarMenu")
	}
	children := slices.Clone(m.Children)
	if !accountingSidebarHasChildKey(children, accountingMainMenuJournalsKey) {
		children = append(children, &components.SidebarMenuItem{
			Page:  components.Page{Key: accountingMainMenuJournalsKey},
			Title: getters.Static("Journals"),
			Url:   lamu.RoutePath("accounting.JournalListRoute", nil),
			Icon:  "book-open",
		})
	}
	return &components.SidebarMenu{
		Page:     m.Page,
		Title:    m.Title,
		Back:     m.Back,
		Children: children,
	}
}

func patchAccountingMainMenuJournalEntries(p components.PageInterface) components.PageInterface {
	m, ok := p.(*components.SidebarMenu)
	if !ok {
		panic("accounting.MainMenu patch expected *components.SidebarMenu")
	}
	children := slices.Clone(m.Children)
	if !accountingSidebarHasChildKey(children, accountingMainMenuJournalEntriesKey) {
		children = append(children, &components.SidebarMenuItem{
			Page:  components.Page{Key: accountingMainMenuJournalEntriesKey},
			Title: getters.Static("Journal entries"),
			Url:   lamu.RoutePath("accounting.JournalEntryListRoute", nil),
			Icon:  "document-text",
		})
	}
	return &components.SidebarMenu{
		Page:     m.Page,
		Title:    m.Title,
		Back:     m.Back,
		Children: children,
	}
}

func pageEntriesJournalPages() []registry.Pair[string, components.PageInterface] {
	createName := getters.Static("accounting.JournalCreateForm")

	return append([]registry.Pair[string, components.PageInterface]{
		{Key: "accounting.JournalDetailMenu", Value: &components.SidebarMenu{
			Title: getters.Format("Journal #%d", getters.Any(getters.Key[uint]("journal.ID"))),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("All journals"),
				Url:   lamu.RoutePath("accounting.JournalListRoute", nil),
			},
			Children: []components.PageInterface{
				&components.SidebarMenuItem{
					Title: getters.Static("Detail"),
					Url: lamu.RoutePath("accounting.JournalDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("journal.ID")),
					}),
				},
				&components.SidebarMenuItem{
					Title: getters.Static("Account transfer"),
					Url: lamu.RoutePath("accounting.JournalAccountTransferRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("journal.ID")),
					}),
					Icon: "arrows-right-left",
				},
			},
		}},
		{Key: "accounting.JournalTable", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "accounting.MainMenu"}},
			Children: []components.PageInterface{
				&components.DataTable[Journal]{
					UID:     "accounting-journal-table",
					Classes: "w-full",
					Data:    getters.Key[components.ObjectList[Journal]]("journals"),
					Actions: []components.PageInterface{
						&components.TableButtonCreate{
							Link: lamu.RoutePath("accounting.JournalCreateRoute", nil),
						},
					},
					RowAttr: getters.RowAttrNavigate(lamu.RoutePath("accounting.JournalDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("$row.ID")),
					})),
					Columns: []components.TableColumn{
						{Label: "Name", Name: "Name", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Name")},
						}},
					},
				},
			},
		}},
		{Key: "accounting.JournalCreateForm", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "accounting.MainMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name:      createName,
					ActionURL: lamu.RoutePath("accounting.JournalCreateRoute", nil),
					Children: []components.PageInterface{
						&components.FormComponent[Journal]{
							Attr:     getters.FormBubbling(createName),
							Title:    "Create journal",
							Subtitle: "Add a new journal",
							ChildrenInput: []components.PageInterface{
								&components.InputText{
									Label:    "Name",
									Name:     "Name",
									Required: true,
									Getter:   getters.Key[string]("$in.Name"),
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
		{Key: "accounting.JournalDetail", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "accounting.JournalDetailMenu"}},
			Children: []components.PageInterface{
				&components.Detail[Journal]{
					Getter: getters.Key[Journal]("journal"),
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
							},
						},
					},
				},
			},
		}},
	}, journalAccountTransferFormPages()...)
}

func journalAccountTransferFormPages() []registry.Pair[string, components.PageInterface] {
	formName := getters.Static("accounting.JournalAccountTransferForm")

	fromAccountInput := journalAccountTransferForeignKey{
		FK: components.InputForeignKey[Account]{
			Name:        "FromAccountID",
			Label:       "From account",
			Url: lamu.RoutePath("accounting.JournalAccountTransferSelectFromRoute", map[string]getters.Getter[any]{
				"id": getters.Any(getters.Key[uint]("journal.ID")),
			}),
			Display:     getters.Format("%s — %s", getters.Any(getters.Key[string]("$in.FromAccount.Code")), getters.Any(getters.Key[string]("$in.FromAccount.Name"))),
			Placeholder: "Select source account...",
			Required:    true,
			Getter:      getters.Association[Account, uint](getters.Key[uint]("$in.FromAccountID")),
		},
		Attr: journalTransferPickerIncludeClosestForm,
	}

	toAccountInput := journalAccountTransferForeignKey{
		FK: components.InputForeignKey[Account]{
			Name:        "ToAccountID",
			Label:       "To account",
			Url: lamu.RoutePath("accounting.JournalAccountTransferSelectToRoute", map[string]getters.Getter[any]{
				"id": getters.Any(getters.Key[uint]("journal.ID")),
			}),
			Display:     getters.Format("%s — %s", getters.Any(getters.Key[string]("$in.ToAccount.Code")), getters.Any(getters.Key[string]("$in.ToAccount.Name"))),
			Placeholder: "Select destination account...",
			Required:    true,
			Getter:      getters.Association[Account, uint](getters.Key[uint]("$in.ToAccountID")),
		},
		Attr: journalTransferPickerIncludeClosestForm,
	}

	return []registry.Pair[string, components.PageInterface]{
		{Key: "accounting.JournalAccountTransferForm", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "accounting.JournalDetailMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name: formName,
					ActionURL: lamu.RoutePath("accounting.JournalAccountTransferRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("journal.ID")),
					}),
					Children: []components.PageInterface{
						&components.FormComponent[JournalAccountTransferForm]{
							Attr:     getters.FormBubbling(formName),
							Title:    "Account transfer",
							Subtitle: "Posts one journal entry in this journal: from account (negative amount) and to account (positive). Amount must be greater than zero.",
							ChildrenInput: []components.PageInterface{
								&components.ContainerError{
									Error: getters.Key[error]("$error.FromAccountID"),
									Children: []components.PageInterface{
										fromAccountInput,
									},
								},
								&components.ContainerError{
									Error: getters.Key[error]("$error.ToAccountID"),
									Children: []components.PageInterface{
										toAccountInput,
									},
								},
								&components.ContainerError{
									Error: getters.Key[error]("$error.Amount"),
									Children: []components.PageInterface{
										&components.InputPointsDecimal{
											Label:    "Amount",
											Name:     "Amount",
											Required: true,
											Getter:   getters.Key[fields.DecimalSix]("$in.Amount"),
										},
									},
								},
							},
							ChildrenAction: []components.PageInterface{
								&components.ButtonSubmit{Label: "Post transfer"},
							},
						},
					},
				},
			},
		}},
		{Key: "accounting.JournalAccountTransferFromAccountSelectionTable", Value: &components.Modal{
			UID: "accounting-journal-account-transfer-from-modal",
			Children: []components.PageInterface{
				&components.DataTable[Account]{
					UID:   "accounting-journal-account-transfer-from-table",
					Title: "Select from account",
					Data:  getters.Key[components.ObjectList[Account]]("accounts"),
					RowAttr: getters.RowAttrSelect("FromAccountID",
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
		{Key: "accounting.JournalAccountTransferToAccountSelectionTable", Value: &components.Modal{
			UID: "accounting-journal-account-transfer-to-modal",
			Children: []components.PageInterface{
				&components.DataTable[Account]{
					UID:   "accounting-journal-account-transfer-to-table",
					Title: "Select to account",
					Data:  getters.Key[components.ObjectList[Account]]("accounts"),
					RowAttr: getters.RowAttrSelect("ToAccountID",
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
	}
}

func pageEntriesJournalEntryPages() []registry.Pair[string, components.PageInterface] {
	createName := getters.Static("accounting.JournalEntryCreateForm")

	journalInput := &components.InputForeignKey[Journal]{
		Name:        "JournalID",
		Label:       "Journal",
		Url:         lamu.RoutePath("accounting.JournalEntryJournalSelectRoute", nil),
		Display:     getters.Key[string]("$in.Journal.Name"),
		Placeholder: "Select journal...",
		Required:    true,
		Getter:      getters.Association[Journal, uint](getters.Key[uint]("$in.JournalID")),
	}

	return []registry.Pair[string, components.PageInterface]{
		{Key: "accounting.JournalEntryDetailMenu", Value: &components.SidebarMenu{
			Title: getters.Format("Journal entry #%d", getters.Any(getters.Key[uint]("journalEntry.ID"))),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("All journal entries"),
				Url:   lamu.RoutePath("accounting.JournalEntryListRoute", nil),
			},
			Children: []components.PageInterface{
				&components.SidebarMenuItem{
					Title: getters.Static("Detail"),
					Url: lamu.RoutePath("accounting.JournalEntryDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("journalEntry.ID")),
					}),
				},
			},
		}},
		{Key: "accounting.JournalEntryTable", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "accounting.MainMenu"}},
			Children: []components.PageInterface{
				&components.DataTable[JournalEntry]{
					UID:     "accounting-journal-entry-table",
					Classes: "w-full",
					Data:    getters.Key[components.ObjectList[JournalEntry]]("journalEntries"),
					Actions: []components.PageInterface{
						&components.TableButtonCreate{
							Link: lamu.RoutePath("accounting.JournalEntryCreateRoute", nil),
						},
					},
					RowAttr: getters.RowAttrNavigate(lamu.RoutePath("accounting.JournalEntryDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("$row.ID")),
					})),
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
		{Key: "accounting.JournalEntryCreateForm", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "accounting.MainMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name:      createName,
					ActionURL: lamu.RoutePath("accounting.JournalEntryCreateRoute", nil),
					Children: []components.PageInterface{
						&components.FormComponent[JournalEntry]{
							Attr:     getters.FormBubbling(createName),
							Title:    "Create journal entry",
							Subtitle: "Header for a balanced posting group",
							ChildrenInput: []components.PageInterface{
								journalInput,
							},
							ChildrenAction: []components.PageInterface{
								&components.ButtonSubmit{Label: "Save"},
							},
						},
					},
				},
			},
		}},
		{Key: "accounting.JournalEntryDetail", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "accounting.JournalEntryDetailMenu"}},
			Children: []components.PageInterface{
				&components.Detail[JournalEntry]{
					Getter: getters.Key[JournalEntry]("journalEntry"),
					Children: []components.PageInterface{
						&components.ContainerColumn{
							Classes: "p-4",
							Children: []components.PageInterface{
								&components.LabelInline{
									Title: "Journal",
									Children: []components.PageInterface{
										&components.FieldText{Getter: getters.Key[string]("$in.Journal.Name")},
									},
								},
							},
						},
					},
				},
			},
		}},
		{Key: "accounting.JournalEntryJournalSelectionTable", Value: &components.Modal{
			UID: "accounting-journal-select-modal",
			Children: []components.PageInterface{
				&components.DataTable[Journal]{
					UID:   "accounting-journal-select-table",
					Title: "Select journal",
					Data:  getters.Key[components.ObjectList[Journal]]("journals"),
					RowAttr: getters.RowAttrSelect("JournalID",
						getters.Key[uint]("$row.ID"),
						getters.Key[string]("$row.Name"),
					),
					Columns: []components.TableColumn{
						{Label: "Name", Name: "Name", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Name")},
						}},
					},
				},
			},
		}},
	}
}
