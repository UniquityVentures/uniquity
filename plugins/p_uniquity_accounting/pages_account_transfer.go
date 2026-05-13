package p_uniquity_accounting

import (
	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/fields"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pageEntriesAccountTransferPages() []registry.Pair[string, components.PageInterface] {
	transferFormName := getters.Static("accounting.AccountTransferForm")

	journalInput := &components.InputForeignKey[Journal]{
		Name:        "JournalID",
		Label:       "Journal",
		Url:         lamu.RoutePath("accounting.JournalEntryJournalSelectRoute", nil),
		Display:     getters.Key[string]("$in.Name"),
		Placeholder: "Select journal...",
		Required:    true,
		Getter:      getters.Association[Journal, uint](getters.Key[uint]("$in.JournalID")),
	}

	toAccountInput := &components.InputForeignKey[Account]{
		Name:        "ToAccountID",
		Label:       "To account",
		Url: lamu.RoutePath("accounting.AccountTransferToAccountSelectRoute", map[string]getters.Getter[any]{
			"id": getters.Any(getters.Key[uint]("account.ID")),
		}),
		Display:     getters.Format("%s — %s", getters.Any(getters.Key[string]("$in.ToAccount.Code")), getters.Any(getters.Key[string]("$in.ToAccount.Name"))),
		Placeholder: "Select destination account...",
		Required:    true,
		Getter:      getters.Association[Account, uint](getters.Key[uint]("$in.ToAccountID")),
	}

	return []registry.Pair[string, components.PageInterface]{
		{Key: "accounting.AccountTransferForm", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "accounting.AccountDetailMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name: transferFormName,
					ActionURL: lamu.RoutePath("accounting.AccountTransferRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("account.ID")),
					}),
					Children: []components.PageInterface{
						&components.FormComponent[AccountTransferForm]{
							Attr:     getters.FormBubbling(transferFormName),
							Title:    "Transfer",
							Subtitle: "Creates one journal entry with two lines: this account (negative amount) and the selected account (positive), equal and opposite.",
							ChildrenInput: []components.PageInterface{
								&components.ContainerError{
									Error: getters.Key[error]("$error.JournalID"),
									Children: []components.PageInterface{
										journalInput,
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
		{Key: "accounting.AccountTransferToAccountSelectionTable", Value: &components.Modal{
			UID: "accounting-account-transfer-to-account-modal",
			Children: []components.PageInterface{
				&components.DataTable[Account]{
					UID:   "accounting-account-transfer-to-account-table",
					Title: "Select destination account",
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
