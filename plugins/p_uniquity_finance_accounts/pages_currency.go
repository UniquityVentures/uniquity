package p_uniquity_finance_accounts

import (
	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func pageCurrencyCRUD() []registry.Pair[string, components.PageInterface] {
	createName := getters.Static("finance_accounts.CurrencyCreateForm")
	updateName := getters.Static("finance_accounts.CurrencyUpdateForm")
	deleteName := getters.Static("finance_accounts.CurrencyDeleteForm")

	codeInput := &components.InputNumber[int]{
		Name:     "Code",
		Label:    "ISO 4217 numeric code",
		Required: true,
		Getter:   getters.Key[int]("$in.Code"),
	}
	nameInput := &components.InputText{
		Name:     "Name",
		Label:    "Name",
		Required: true,
		Getter:   getters.Key[string]("$in.Name"),
	}
	symbolInput := &components.InputText{
		Name:     "Symbol",
		Label:    "Symbol",
		Required: true,
		Getter:   getters.Key[string]("$in.Symbol"),
	}
	minorUnitInput := &components.InputNumber[int]{
		Name:     "MinorUnit",
		Label:    "Minor unit (decimal places)",
		Required: true,
		Getter:   getters.Key[int]("$in.MinorUnit"),
	}

	return []registry.Pair[string, components.PageInterface]{
		{Key: "finance_accounts.CurrencyTable", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_accounts.MainMenu"}},
			Children: []components.PageInterface{
				&components.DataTable[Currency]{
					UID:     "finance-currencies-table",
					Title:   "Currencies",
					Classes: "w-full",
					Data:    getters.Key[components.ObjectList[Currency]]("currencies"),
					Actions: []components.PageInterface{
						&components.TableButtonFilter{Child: lamu.DynamicPage{Name: "finance_accounts.CurrencyFilter"}},
						&components.TableButtonCreate{
							Link: lamu.RoutePath("finance_accounts.CurrencyCreateRoute", nil),
							Page: components.Page{Roles: []string{"superuser"}},
						},
					},
					RowAttr: getters.RowAttrNavigate(lamu.RoutePath("finance_accounts.CurrencyDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("$row.ID")),
					})),
					Columns: []components.TableColumn{
						{Label: "Code", Name: "Code", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Format("%d", getters.Any(getters.Key[int]("$row.Code")))},
						}},
						{Label: "Name", Name: "Name", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Name")},
						}},
						{Label: "Symbol", Name: "Symbol", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Symbol")},
						}},
						{Label: "Minor unit", Name: "MinorUnit", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Format("%d", getters.Any(getters.Key[int]("$row.MinorUnit")))},
						}},
					},
				},
			},
		}},
		{Key: "finance_accounts.CurrencyCreateForm", Value: &components.ShellScaffold{
			Page:    components.Page{Roles: []string{"superuser"}},
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_accounts.MainMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name:      createName,
					ActionURL: lamu.RoutePath("finance_accounts.CurrencyCreateRoute", nil),
					Children: []components.PageInterface{
						&components.FormComponent[Currency]{
							Attr:     getters.FormBubbling(createName),
							Title:    "Create currency",
							Subtitle: "ISO 4217 currency",
							ChildrenInput: []components.PageInterface{
								codeInput,
								nameInput,
								symbolInput,
								minorUnitInput,
							},
							ChildrenAction: []components.PageInterface{
								&components.ButtonSubmit{Label: "Save"},
							},
						},
					},
				},
			},
		}},
		{Key: "finance_accounts.CurrencyUpdateForm", Value: &components.ShellScaffold{
			Page:    components.Page{Roles: []string{"superuser"}},
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_accounts.CurrencyDetailMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name: updateName,
					ActionURL: lamu.RoutePath("finance_accounts.CurrencyUpdateRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("currency.ID")),
					}),
					Children: []components.PageInterface{
						&components.FormComponent[Currency]{
							Getter:   getters.Key[Currency]("currency"),
							Attr:     getters.FormBubbling(updateName),
							Title:    "Edit currency",
							Subtitle: "Update ISO 4217 fields",
							ChildrenInput: []components.PageInterface{
								codeInput,
								nameInput,
								symbolInput,
								minorUnitInput,
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
													Url: lamu.RoutePath("finance_accounts.CurrencyDeleteRoute", map[string]getters.Getter[any]{
														"id": getters.Any(getters.Key[uint]("currency.ID")),
													}),
													FormPostURL: lamu.RoutePath("finance_accounts.CurrencyDeleteRoute", map[string]getters.Getter[any]{
														"id": getters.Any(getters.Key[uint]("currency.ID")),
													}),
													ModalUID: "finance-currency-delete-modal",
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
		{Key: "finance_accounts.CurrencyDeleteForm", Value: &components.Modal{
			Page: components.Page{Roles: []string{"superuser"}},
			UID:  "finance-currency-delete-modal",
			Children: []components.PageInterface{
				&components.DeleteConfirmation{
					Title:   "Delete currency?",
					Message: "This removes the currency. Journals that use it must be changed first.",
					Attr:    getters.FormBubbling(getters.Key[string]("$get.name")),
				},
			},
		}},
		{Key: "finance_accounts.CurrencyDetail", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_accounts.CurrencyDetailMenu"}},
			Children: []components.PageInterface{
				&components.Detail[Currency]{
					Getter: getters.Key[Currency]("currency"),
					Children: []components.PageInterface{
						&components.ContainerColumn{
							Classes: "p-4",
							Children: []components.PageInterface{
								&components.FieldTitle{Getter: getters.Key[string]("$in.Name")},
								&components.FieldSubtitle{Getter: getters.Format("Code %d · %s", getters.Any(getters.Key[int]("$in.Code")), getters.Any(getters.Key[string]("$in.Symbol")))},
								&components.LabelInline{
									Title:   "Minor unit",
									Classes: "mt-2",
									Children: []components.PageInterface{
										&components.FieldText{Getter: getters.Format("%d", getters.Any(getters.Key[int]("$in.MinorUnit")))},
									},
								},
							},
						},
					},
				},
			},
		}},
		{Key: "finance_accounts.CurrencySelectionTable", Value: &components.Modal{
			UID: "finance-currency-selection-modal",
			Children: []components.PageInterface{
				&components.DataTable[Currency]{
					UID:   "finance-currency-selection-table",
					Title: "Select currency",
					Data:  getters.Key[components.ObjectList[Currency]]("currencies"),
					Actions: []components.PageInterface{
						&components.TableButtonFilter{Child: lamu.DynamicPage{Name: "finance_accounts.CurrencySelectionFilter"}},
					},
					RowAttr: getters.RowAttrSelect("CurrencyID", getters.Key[uint]("$row.ID"), getters.Format("%s — %s (%d)",
						getters.Any(getters.Key[string]("$row.Symbol")),
						getters.Any(getters.Key[string]("$row.Name")),
						getters.Any(getters.Key[int]("$row.Code")),
					)),
					Columns: []components.TableColumn{
						{Label: "Code", Name: "Code", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Format("%d", getters.Any(getters.Key[int]("$row.Code")))},
						}},
						{Label: "Name", Name: "Name", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Name")},
						}},
						{Label: "Symbol", Name: "Symbol", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Symbol")},
						}},
					},
				},
			},
		}},
	}
}
