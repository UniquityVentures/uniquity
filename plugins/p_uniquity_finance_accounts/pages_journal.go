package p_uniquity_finance_accounts

import (
	"context"
	"time"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/fields"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

var journalTypeChoices = getters.Static([]registry.Pair[JournalType, string]{
	{Key: JournalTypeGeneral, Value: "General"},
})

func journalTypeSelectGetter(ctxKey string) getters.Getter[registry.Pair[JournalType, string]] {
	return func(ctx context.Context) (registry.Pair[JournalType, string], error) {
		jt, err := getters.Key[JournalType](ctxKey)(ctx)
		if err != nil {
			return registry.Pair[JournalType, string]{}, err
		}
		if jt == "" {
			return registry.Pair[JournalType, string]{}, nil
		}
		return registry.Pair[JournalType, string]{Key: jt, Value: string(jt)}, nil
	}
}

func journalCurrencyDetailHref() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		cid, err := getters.Key[uint]("$in.CurrencyID")(ctx)
		if err != nil {
			return "", err
		}
		if cid == 0 {
			return "", nil
		}
		return lamu.RoutePath("finance_accounts.CurrencyDetailRoute", map[string]getters.Getter[any]{
			"id": func(context.Context) (any, error) { return cid, nil },
		})(ctx)
	}
}

func journalCurrencySummary(rowPrefix string) getters.Getter[string] {
	return getters.Format("%s — %s (%d)",
		getters.Any(getters.Key[string](rowPrefix+".Currency.Symbol")),
		getters.Any(getters.Key[string](rowPrefix+".Currency.Name")),
		getters.Any(getters.Key[int](rowPrefix+".Currency.Code")),
	)
}

func journalEntryDatetimeText(rowPrefix string) getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		t, err := getters.Key[time.Time](rowPrefix + ".Datetime")(ctx)
		if err != nil {
			return "", err
		}
		if t.IsZero() {
			return "", nil
		}
		return t.Format(time.DateTime), nil
	}
}

func journalEntrySourceDocSummary(rowPrefix string) getters.Getter[string] {
	return getters.Format("%s · id %d",
		getters.Any(getters.Key[string](rowPrefix+".SourceDoc.Type")),
		getters.Any(getters.Key[uint](rowPrefix+".SourceDocID")),
	)
}

func journalEntryCreateDatetimeGetter(ctxKey string) getters.Getter[time.Time] {
	return func(ctx context.Context) (time.Time, error) {
		t, err := getters.Key[time.Time](ctxKey)(ctx)
		if err != nil {
			return time.Time{}, err
		}
		if t.IsZero() {
			return time.Now(), nil
		}
		return t, nil
	}
}

func journalEntryCreateHref() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		jid, err := getters.Key[uint]("journal.ID")(ctx)
		if err != nil {
			return "", err
		}
		return lamu.RoutePath("finance_accounts.JournalEntryCreateRoute", map[string]getters.Getter[any]{
			"journal_id": func(context.Context) (any, error) { return jid, nil },
		})(ctx)
	}
}

func journalEntryParentJournalHref() getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		jid, err := getters.Key[uint]("$in.JournalID")(ctx)
		if err != nil {
			return "", err
		}
		return lamu.RoutePath("finance_accounts.JournalDetailRoute", map[string]getters.Getter[any]{
			"id": func(context.Context) (any, error) { return jid, nil },
		})(ctx)
	}
}

func journalEntryParentJournalLabel() getters.Getter[string] {
	return getters.Format("%s (#%d)",
		getters.Any(getters.Key[string]("$in.Journal.Name")),
		getters.Any(getters.Key[uint]("$in.Journal.ID")),
	)
}

func journalEntryItemAmountText(rowPrefix string) getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		a, err := getters.Key[fields.DecimalSix](rowPrefix + ".Amount")(ctx)
		if err != nil {
			return "", err
		}
		return a.String(), nil
	}
}

func pageJournalCRUD() []registry.Pair[string, components.PageInterface] {
	createName := getters.Static("finance_accounts.JournalCreateForm")
	updateName := getters.Static("finance_accounts.JournalUpdateForm")
	deleteName := getters.Static("finance_accounts.JournalDeleteForm")

	nameInput := &components.InputText{
		Name:     "Name",
		Label:    "Name",
		Required: true,
		Getter:   getters.Key[string]("$in.Name"),
	}
	activeInput := &components.InputCheckbox{
		Name:   "IsActive",
		Label:  "Active",
		Getter: getters.Key[bool]("$in.IsActive"),
	}
	currencyPicker := &components.InputForeignKey[Currency]{
		Name:        "CurrencyID",
		Label:       "Currency",
		Url:         lamu.RoutePath("finance_accounts.CurrencySelectRoute", nil),
		Display:     getters.Format("%s — %s (%d)", getters.Any(getters.Key[string]("$in.Symbol")), getters.Any(getters.Key[string]("$in.Name")), getters.Any(getters.Key[int]("$in.Code"))),
		Placeholder: "Select currency…",
		Required:    true,
		Getter:      getters.Association[Currency, uint](getters.Key[uint]("$in.CurrencyID")),
	}
	typeInput := &components.InputSelect[JournalType]{
		Name:     "Type",
		Label:    "Type",
		Required: true,
		Choices:  journalTypeChoices,
		Getter:   journalTypeSelectGetter("$in.Type"),
	}

	return []registry.Pair[string, components.PageInterface]{
		{Key: "finance_accounts.JournalTable", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_accounts.MainMenu"}},
			Children: []components.PageInterface{
				&components.DataTable[Journal]{
					UID:     "finance-journals-table",
					Title:   "Journals",
					Classes: "w-full",
					Data:    getters.Key[components.ObjectList[Journal]]("journals"),
					Actions: []components.PageInterface{
						&components.TableButtonFilter{Child: lamu.DynamicPage{Name: "finance_accounts.JournalFilter"}},
						&components.TableButtonCreate{
							Link: lamu.RoutePath("finance_accounts.JournalCreateRoute", nil),
							Page: components.Page{Roles: []string{"superuser"}},
						},
					},
					RowAttr: getters.RowAttrNavigate(lamu.RoutePath("finance_accounts.JournalDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("$row.ID")),
					})),
					Columns: []components.TableColumn{
						{Label: "Name", Name: "Name", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Name")},
						}},
						{Label: "Active", Name: "IsActive", Children: []components.PageInterface{
							&components.FieldCheckbox{Getter: getters.Key[bool]("$row.IsActive")},
						}},
						{Label: "Currency", Name: "Currency", Children: []components.PageInterface{
							&components.FieldText{Getter: journalCurrencySummary("$row")},
						}},
						{Label: "Type", Name: "Type", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Format("%s", getters.Any(getters.Key[JournalType]("$row.Type")))},
						}},
					},
				},
			},
		}},
		{Key: "finance_accounts.JournalCreateForm", Value: &components.ShellScaffold{
			Page:    components.Page{Roles: []string{"superuser"}},
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_accounts.MainMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name:      createName,
					ActionURL: lamu.RoutePath("finance_accounts.JournalCreateRoute", nil),
					Children: []components.PageInterface{
						&components.FormComponent[Journal]{
							Attr:     getters.FormBubbling(createName),
							Title:    "Create journal",
							Subtitle: "Journal in one currency",
							ChildrenInput: []components.PageInterface{
								nameInput,
								activeInput,
								currencyPicker,
								typeInput,
							},
							ChildrenAction: []components.PageInterface{
								&components.ButtonSubmit{Label: "Save"},
							},
						},
					},
				},
			},
		}},
		{Key: "finance_accounts.JournalUpdateForm", Value: &components.ShellScaffold{
			Page:    components.Page{Roles: []string{"superuser"}},
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_accounts.JournalDetailMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name: updateName,
					ActionURL: lamu.RoutePath("finance_accounts.JournalUpdateRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("journal.ID")),
					}),
					Children: []components.PageInterface{
						&components.FormComponent[Journal]{
							Getter:   getters.Key[Journal]("journal"),
							Attr:     getters.FormBubbling(updateName),
							Title:    "Edit journal",
							Subtitle: "Update journal settings",
							ChildrenInput: []components.PageInterface{
								nameInput,
								activeInput,
								currencyPicker,
								typeInput,
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
													Url: lamu.RoutePath("finance_accounts.JournalDeleteRoute", map[string]getters.Getter[any]{
														"id": getters.Any(getters.Key[uint]("journal.ID")),
													}),
													FormPostURL: lamu.RoutePath("finance_accounts.JournalDeleteRoute", map[string]getters.Getter[any]{
														"id": getters.Any(getters.Key[uint]("journal.ID")),
													}),
													ModalUID: "finance-journal-delete-modal",
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
		{Key: "finance_accounts.JournalDeleteForm", Value: &components.Modal{
			Page: components.Page{Roles: []string{"superuser"}},
			UID:  "finance-journal-delete-modal",
			Children: []components.PageInterface{
				&components.DeleteConfirmation{
					Title:   "Delete journal?",
					Message: "This removes the journal record.",
					Attr:    getters.FormBubbling(getters.Key[string]("$get.name")),
				},
			},
		}},
		{Key: "finance_accounts.JournalDetail", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_accounts.JournalDetailMenu"}},
			Children: []components.PageInterface{
				&components.Detail[Journal]{
					Getter: getters.Key[Journal]("journal"),
					Children: []components.PageInterface{
						&components.ContainerColumn{
							Classes: "p-4",
							Children: []components.PageInterface{
								&components.FieldTitle{Getter: getters.Key[string]("$in.Name")},
								&components.ShowIf{
									Getter: func(ctx context.Context) (any, error) {
										return getters.Key[bool]("$in.IsActive")(ctx)
									},
									Children: []components.PageInterface{
										&components.FieldSubtitle{Getter: getters.Static("Active")},
									},
								},
								&components.ShowIf{
									Getter: func(ctx context.Context) (any, error) {
										active, err := getters.Key[bool]("$in.IsActive")(ctx)
										if err != nil {
											return false, err
										}
										return !active, nil
									},
									Children: []components.PageInterface{
										&components.FieldSubtitle{Getter: getters.Static("Inactive")},
									},
								},
								&components.LabelInline{
									Title:   "Currency",
									Classes: "mt-2",
									Children: []components.PageInterface{
										&components.FieldLink{
											Href:  journalCurrencyDetailHref(),
											Label: journalCurrencySummary("$in"),
										},
									},
								},
								&components.LabelInline{
									Title: "Type",
									Children: []components.PageInterface{
										&components.FieldText{Getter: getters.Format("%s", getters.Any(getters.Key[JournalType]("$in.Type")))},
									},
								},
							},
						},
					},
				},
				&components.DataTable[JournalEntry]{
					UID:      "finance-journal-entries-table",
					Title:    "Journal entries",
					Subtitle: "Entries posted to this journal",
					Classes:  "w-full",
					Data:     getters.Key[components.ObjectList[JournalEntry]](journalDetailEntriesContextKey),
					RowAttr: getters.RowAttrNavigate(lamu.RoutePath("finance_accounts.JournalEntryDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("$row.ID")),
					})),
					Actions: []components.PageInterface{
						&components.TableButtonCreate{
							Link:  journalEntryCreateHref(),
							Page:  components.Page{Roles: []string{"superuser"}},
							Label: "New entry",
						},
					},
					Columns: []components.TableColumn{
						{Label: "Date & time", Name: "Datetime", Children: []components.PageInterface{
							&components.FieldText{Getter: journalEntryDatetimeText("$row")},
						}},
						{Label: "Source document", Name: "SourceDoc", Children: []components.PageInterface{
							&components.FieldText{Getter: journalEntrySourceDocSummary("$row")},
						}},
					},
				},
			},
		}},
	}
}

func pageJournalEntryCreatePages() []registry.Pair[string, components.PageInterface] {
	entryCreateName := getters.Static("finance_accounts.JournalEntryCreateForm")
	sourceDocPicker := &components.InputForeignKey[SourceDoc]{
		Name:        "SourceDocID",
		Label:       "Source document",
		Url:         lamu.RoutePath("finance_accounts.SourceDocSelectRoute", nil),
		Display:     getters.Format("%s · ref %d · #%d", getters.Any(getters.Key[string]("$in.Type")), getters.Any(getters.Key[uint]("$in.SourceDocID")), getters.Any(getters.Key[uint]("$in.ID"))),
		Placeholder: "Select source document…",
		Required:    true,
		Getter:      getters.Association[SourceDoc, uint](getters.Key[uint]("$in.SourceDocID")),
	}
	return []registry.Pair[string, components.PageInterface]{
		{Key: "finance_accounts.JournalEntryCreateForm", Value: &components.ShellScaffold{
			Page:    components.Page{Roles: []string{"superuser"}},
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_accounts.JournalDetailMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name: entryCreateName,
					ActionURL: lamu.RoutePath("finance_accounts.JournalEntryCreateRoute", map[string]getters.Getter[any]{
						"journal_id": getters.Any(getters.Key[uint]("journal.ID")),
					}),
					Children: []components.PageInterface{
						&components.FormComponent[JournalEntry]{
							Attr:     getters.FormBubbling(entryCreateName),
							Title:    "Create journal entry",
							Subtitle: "Set the entry time and linked source document.",
							ChildrenInput: []components.PageInterface{
								&components.InputDatetime{
									Name:     "Datetime",
									Label:    "Date & time",
									Required: true,
									Getter:   journalEntryCreateDatetimeGetter("$in.Datetime"),
								},
								sourceDocPicker,
							},
							ChildrenAction: []components.PageInterface{
								&components.ButtonSubmit{Label: "Save"},
							},
						},
					},
				},
			},
		}},
		{Key: "finance_accounts.SourceDocSelectionTable", Value: &components.Modal{
			UID: "finance-sourcedoc-selection-modal",
			Children: []components.PageInterface{
				&components.DataTable[SourceDoc]{
					UID:   "finance-sourcedoc-selection-table",
					Title: "Select source document",
					Data:  getters.Key[components.ObjectList[SourceDoc]]("source_docs"),
					RowAttr: getters.RowAttrSelect("SourceDocID", getters.Key[uint]("$row.ID"), getters.Format("%s · ref %d · #%d",
						getters.Any(getters.Key[string]("$row.Type")),
						getters.Any(getters.Key[uint]("$row.SourceDocID")),
						getters.Any(getters.Key[uint]("$row.ID")),
					)),
					Columns: []components.TableColumn{
						{Label: "Type", Name: "Type", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Type")},
						}},
						{Label: "Document id", Name: "SourceDocID", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Format("%d", getters.Any(getters.Key[uint]("$row.SourceDocID")))},
						}},
						{Label: "Row id", Name: "ID", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Format("%d", getters.Any(getters.Key[uint]("$row.ID")))},
						}},
					},
				},
			},
		}},
	}
}

func pageJournalEntryDetailPages() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "finance_accounts.JournalEntryDetail", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_accounts.JournalEntryDetailMenu"}},
			Children: []components.PageInterface{
				&components.Detail[JournalEntry]{
					Getter: getters.Key[JournalEntry]("journalEntry"),
					Children: []components.PageInterface{
						&components.ContainerColumn{
							Classes: "p-4",
							Children: []components.PageInterface{
								&components.FieldTitle{Getter: getters.Static("Journal entry")},
								&components.FieldSubtitle{Getter: journalEntryDatetimeText("$in")},
								&components.LabelInline{
									Title:   "Journal",
									Classes: "mt-2",
									Children: []components.PageInterface{
										&components.FieldLink{
											Href:  journalEntryParentJournalHref(),
											Label: journalEntryParentJournalLabel(),
										},
									},
								},
								&components.LabelInline{
									Title: "Source document",
									Children: []components.PageInterface{
										&components.FieldText{Getter: journalEntrySourceDocSummary("$in")},
									},
								},
							},
						},
					},
				},
				&components.DataTable[JournalEntryItem]{
					UID:      "finance-journal-entry-items-table",
					Title:    "Lines",
					Subtitle: "Journal entry items",
					Classes:  "w-full",
					Data:     getters.Key[components.ObjectList[JournalEntryItem]](journalEntryDetailItemsContextKey),
					Columns: []components.TableColumn{
						{Label: "Date & time", Name: "Datetime", Children: []components.PageInterface{
							&components.FieldText{Getter: journalEntryDatetimeText("$row")},
						}},
						{Label: "Account", Name: "Account", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Format("%d — %s",
								getters.Any(getters.Key[int]("$row.Account.Code")),
								getters.Any(getters.Key[string]("$row.Account.Name")),
							)},
						}},
						{Label: "Amount", Name: "Amount", Children: []components.PageInterface{
							&components.FieldText{Getter: journalEntryItemAmountText("$row")},
						}},
					},
				},
			},
		}},
	}
}
