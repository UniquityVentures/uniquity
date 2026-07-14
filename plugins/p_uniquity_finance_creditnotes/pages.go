package p_uniquity_finance_creditnotes

import (
	"context"
	"fmt"
	"time"

	finance_accounts "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_accounts"
	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/components"
	"github.com/lariv-in/lago/getters"
	"github.com/lariv-in/lago/registry"
)

const financeAccountsMainMenuCreditNotesLinkKey = "finance_credit_notes.FinanceAccountsMainMenuLink"

func patchFinanceAccountsMainMenuForCreditNotes(page components.PageInterface) components.PageInterface {
	menu, ok := page.(*components.SidebarMenu)
	if !ok {
		panic("p_uniquity_finance_creditnotes: finance_accounts.MainMenu must be *components.SidebarMenu")
	}
	for _, ch := range menu.Children {
		if item, ok := ch.(*components.SidebarMenuItem); ok && item.GetKey() == financeAccountsMainMenuCreditNotesLinkKey {
			return menu
		}
	}
	newChildren := append([]components.PageInterface{}, menu.Children...)
	newChildren = append(newChildren, &components.SidebarMenuItem{
		Page:  components.Page{Key: financeAccountsMainMenuCreditNotesLinkKey, Roles: []string{"superuser"}},
		Title: getters.Static("Credit notes"),
		Url:   lago.RoutePath("finance_credit_notes.DefaultRoute", nil),
		Icon:  "arrow-uturn-left",
	})
	cloned := *menu
	cloned.Children = newChildren
	return &cloned
}

func pluginPages() lago.PluginFeatures[components.PageInterface] {
	e := pageEntriesCreditNotePages()
	e = append(e, pageEntriesJournalEntryFkSelectPages()...)
	return lago.PluginFeatures[components.PageInterface]{
		Entries: e,
		Patches: []registry.Pair[string, func(components.PageInterface) components.PageInterface]{
			{Key: "finance_accounts.MainMenu", Value: patchFinanceAccountsMainMenuForCreditNotes},
		},
	}
}

func creditNoteCreateDatetimeGetter() getters.Getter[time.Time] {
	return func(ctx context.Context) (time.Time, error) {
		t, err := getters.Key[time.Time]("$in.Datetime")(ctx)
		if err != nil {
			return time.Time{}, nil
		}
		return t, nil
	}
}

func journalEntryFkDisplayGetter() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		id, err := getters.Key[uint]("$in.ID")(ctx)
		if err != nil || id == 0 {
			return "", nil
		}
		jname, jerr := getters.Key[string]("$in.Journal.Name")(ctx)
		if jerr != nil || jname == "" {
			jname = "—"
		}
		dt, derr := getters.Key[time.Time]("$in.Datetime")(ctx)
		if derr != nil || dt.IsZero() {
			return fmt.Sprintf("#%d · %s", id, jname), nil
		}
		return fmt.Sprintf("#%d · %s · %s", id, jname, dt.Format(time.DateTime)), nil
	}
}

func creditNoteListJournalEntrySummary(rowPrefix string) getters.Getter[string] {
	return getters.Format("#%d", getters.Any(getters.Key[uint](rowPrefix+".JournalEntryID")))
}

func creditNoteListReversedEntrySummary(rowPrefix string) getters.Getter[string] {
	return getters.Format("#%d", getters.Any(getters.Key[uint](rowPrefix+".ReversedJournalEntryID")))
}

func creditNoteCreateFormInputs() []components.PageInterface {
	return []components.PageInterface{
		&components.ContainerError{
			Error: getters.Key[error]("$error.Datetime"),
			Children: []components.PageInterface{
				&components.InputDatetime{
					Label:    "Credit note date & time",
					Name:     "Datetime",
					Required: true,
					Getter:   creditNoteCreateDatetimeGetter(),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Reason"),
			Children: []components.PageInterface{
				&components.InputTextarea{Name: "Reason", Label: "Reason", Rows: 4},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.JournalEntryID"),
			Children: []components.PageInterface{
				&components.InputForeignKey[finance_accounts.JournalEntry]{
					Label:       "Journal entry to reverse",
					Name:        "JournalEntryID",
					Required:    true,
					Url:         lago.RoutePath("finance_credit_notes.JournalEntryFkSelectRoute", nil),
					Display:     journalEntryFkDisplayGetter(),
					Placeholder: "Select journal entry…",
					Getter:      getters.Association[finance_accounts.JournalEntry, uint](getters.Key[uint]("$in.JournalEntryID")),
				},
			},
		},
	}
}

func pageEntriesCreditNotePages() []registry.Pair[string, components.PageInterface] {
	createName := getters.Static("finance_credit_notes.CreditNoteCreateForm")

	return []registry.Pair[string, components.PageInterface]{
		{Key: "finance_credit_notes.CreditNoteTable", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lago.DynamicPage{Name: "finance_accounts.MainMenu"}},
			Children: []components.PageInterface{
				&components.DataTable[CreditNote]{
					UID:     "finance-credit-notes-table",
					Classes: "w-full",
					Data:    getters.Key[components.ObjectList[CreditNote]]("credit_notes"),
					RowAttr: getters.RowAttrNavigate(lago.RoutePath("finance_credit_notes.CreditNoteDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("$row.ID")),
					})),
					Actions: []components.PageInterface{
						&components.TableButtonCreate{
							Link: lago.RoutePath("finance_credit_notes.CreditNoteCreateRoute", nil),
							Page: components.Page{Roles: []string{"superuser"}},
						},
					},
					Columns: []components.TableColumn{
						{Label: "Date", Name: "Datetime", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Format("%s", getters.Any(getters.Key[time.Time]("$row.Datetime")))},
						}},
						{Label: "Reason", Name: "Reason", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Reason")},
						}},
						{Label: "Original entry", Name: "JournalEntryID", Children: []components.PageInterface{
							&components.FieldText{Getter: creditNoteListJournalEntrySummary("$row")},
						}},
						{Label: "Reversal entry", Name: "ReversedJournalEntryID", Children: []components.PageInterface{
							&components.FieldText{Getter: creditNoteListReversedEntrySummary("$row")},
						}},
					},
				},
			},
		}},
		{Key: "finance_credit_notes.CreditNoteDetail", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lago.DynamicPage{Name: "finance_credit_notes.CreditNoteDetailMenu"}},
			Children: []components.PageInterface{
				&components.ContainerError{
					Error: getters.Key[error]("$error._global"),
					Children: []components.PageInterface{
						&components.Detail[CreditNote]{
							Getter: getters.Key[CreditNote]("credit_note"),
							Children: []components.PageInterface{
								&components.ContainerColumn{
									Classes: "p-4",
									Children: []components.PageInterface{
										&components.LabelInline{Title: "Date", Children: []components.PageInterface{
											&components.FieldText{Getter: getters.Format("%s", getters.Any(getters.Key[time.Time]("$in.Datetime")))},
										}},
										&components.LabelInline{Title: "Reason", Children: []components.PageInterface{
											&components.FieldText{Getter: getters.Key[string]("$in.Reason")},
										}},
										&components.LabelInline{Title: "Original journal entry", Children: []components.PageInterface{
											&components.FieldLink{
												Href: lago.RoutePath("finance_accounts.JournalEntryDetailRoute", map[string]getters.Getter[any]{
													"id": getters.Any(getters.Key[uint]("$in.JournalEntryID")),
												}),
												Label:   getters.Format("Entry #%d", getters.Any(getters.Key[uint]("$in.JournalEntryID"))),
												Classes: "link link-hover",
											},
										}},
										&components.LabelInline{Title: "Reversal journal entry", Children: []components.PageInterface{
											&components.FieldLink{
												Href: lago.RoutePath("finance_accounts.JournalEntryDetailRoute", map[string]getters.Getter[any]{
													"id": getters.Any(getters.Key[uint]("$in.ReversedJournalEntryID")),
												}),
												Label:   getters.Format("Entry #%d", getters.Any(getters.Key[uint]("$in.ReversedJournalEntryID"))),
												Classes: "link link-hover",
											},
										}},
									},
								},
							},
						},
					},
				},
			},
		}},
		{Key: "finance_credit_notes.CreditNoteDetailMenu", Value: &components.SidebarMenu{
			Title: getters.Format("Credit note #%d", getters.Any(getters.Key[uint]("credit_note.ID"))),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("Credit notes"),
				Url:   lago.RoutePath("finance_credit_notes.DefaultRoute", nil),
			},
			Children: []components.PageInterface{
				&components.SidebarMenuItem{
					Title: getters.Static("Detail"),
					Url: lago.RoutePath("finance_credit_notes.CreditNoteDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("credit_note.ID")),
					}),
				},
			},
		}},
		{Key: "finance_credit_notes.CreditNoteCreateForm", Value: &components.ShellScaffold{
			Page:    components.Page{Roles: []string{"superuser"}},
			Sidebar: []components.PageInterface{lago.DynamicPage{Name: "finance_accounts.MainMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name:      createName,
					ActionURL: lago.RoutePath("finance_credit_notes.CreditNoteCreateRoute", nil),
					Children: []components.PageInterface{
						&components.FormComponent[CreditNote]{
							Attr:          getters.FormBubbling(createName),
							Title:         "Create credit note",
							Subtitle:      "Reverses the selected journal entry with an opposite posting",
							ChildrenInput: creditNoteCreateFormInputs(),
							ChildrenAction: []components.PageInterface{
								&components.ButtonSubmit{Label: "Save"},
							},
						},
					},
				},
			},
		}},
	}
}

func pageEntriesJournalEntryFkSelectPages() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "finance_credit_notes.JournalEntryFkSelectionTable", Value: &components.Modal{
			UID: "finance-journal-entry-fk-select-modal",
			Children: []components.PageInterface{
				&components.DataTable[finance_accounts.JournalEntry]{
					UID:   "finance-journal-entry-fk-select-table",
					Title: "Select journal entry",
					Data:  getters.Key[components.ObjectList[finance_accounts.JournalEntry]]("journal_entries"),
					RowAttr: getters.RowAttrSelect(
						"JournalEntryID",
						getters.Key[uint]("$row.ID"),
						journalEntryFkRowLabel(),
					),
					Columns: []components.TableColumn{
						{Label: "When", Name: "Datetime", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Format("%s", getters.Any(getters.Key[time.Time]("$row.Datetime")))},
						}},
						{Label: "Journal", Name: "Journal", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Journal.Name")},
						}},
						{Label: "Entry", Name: "ID", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Format("#%d", getters.Any(getters.Key[uint]("$row.ID")))},
						}},
					},
				},
			},
		}},
	}
}

func journalEntryFkRowLabel() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		id, err := getters.Key[uint]("$row.ID")(ctx)
		if err != nil {
			return "", err
		}
		jname, jerr := getters.Key[string]("$row.Journal.Name")(ctx)
		if jerr != nil || jname == "" {
			jname = "—"
		}
		dt, derr := getters.Key[time.Time]("$row.Datetime")(ctx)
		if derr != nil || dt.IsZero() {
			return fmt.Sprintf("#%d · %s", id, jname), nil
		}
		return fmt.Sprintf("#%d · %s · %s", id, jname, dt.Format(time.DateTime)), nil
	}
}
