package p_uniquity_finance_taxes

import (
	"context"
	"fmt"
	"strings"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/fields"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
	finance_accounts "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_accounts"
)

const financeAccountsMainMenuTaxesLinkKey = "finance_taxes.FinanceAccountsMainMenuLink"

func patchFinanceAccountsMainMenuForTaxes(page components.PageInterface) components.PageInterface {
	menu, ok := page.(*components.SidebarMenu)
	if !ok {
		panic("p_uniquity_finance_taxes: finance_accounts.MainMenu must be *components.SidebarMenu")
	}
	for _, ch := range menu.Children {
		if item, ok := ch.(*components.SidebarMenuItem); ok && item.GetKey() == financeAccountsMainMenuTaxesLinkKey {
			return menu
		}
	}
	newChildren := append([]components.PageInterface{}, menu.Children...)
	newChildren = append(newChildren, &components.SidebarMenuItem{
		Page:  components.Page{Key: financeAccountsMainMenuTaxesLinkKey, Roles: []string{"superuser"}},
		Title: getters.Static("Taxes"),
		Url:   lamu.RoutePath("finance_taxes.DefaultRoute", nil),
		Icon:  "calculator",
	})
	cloned := *menu
	cloned.Children = newChildren
	return &cloned
}

func pluginPages() lamu.PluginFeatures[components.PageInterface] {
	e := pageEntriesTaxMenus()
	e = append(e, pageEntriesTaxPages()...)
	e = append(e, pageEntriesTaxMultiSelectPages()...)
	return lamu.PluginFeatures[components.PageInterface]{
		Entries: e,
		Patches: []registry.Pair[string, func(components.PageInterface) components.PageInterface]{
			{Key: "finance_accounts.MainMenu", Value: patchFinanceAccountsMainMenuForTaxes},
		},
	}
}

func taxDecimalStringGetter(ctxKey string) getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		pd, err := getters.Key[fields.DecimalSix](ctxKey)(ctx)
		if err != nil {
			return "", err
		}
		return pd.String(), nil
	}
}

func taxDecimalGetter(ctxKey string) getters.Getter[fields.DecimalSix] {
	return func(ctx context.Context) (fields.DecimalSix, error) {
		return getters.Key[fields.DecimalSix](ctxKey)(ctx)
	}
}

var taxKindChoices = getters.Static([]registry.Pair[TaxKind, string]{
	{Key: TaxKindLevied, Value: "Levied"},
	{Key: TaxKindWithholding, Value: "Withholding"},
})

func taxKindSelectGetter(ctxKey string) getters.Getter[registry.Pair[TaxKind, string]] {
	return func(ctx context.Context) (registry.Pair[TaxKind, string], error) {
		k, err := getters.Key[TaxKind](ctxKey)(ctx)
		if err != nil {
			return registry.Pair[TaxKind, string]{}, err
		}
		if k == "" {
			return registry.Pair[TaxKind, string]{}, nil
		}
		label := string(k)
		if k == TaxKindLevied {
			label = "Levied"
		}
		if k == TaxKindWithholding {
			label = "Withholding"
		}
		return registry.Pair[TaxKind, string]{Key: k, Value: label}, nil
	}
}

func taxKindLabel(ctxKey string) getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		k, err := getters.Key[TaxKind](ctxKey)(ctx)
		if err != nil {
			return "", err
		}
		switch k {
		case TaxKindLevied:
			return "Levied", nil
		case TaxKindWithholding:
			return "Withholding", nil
		default:
			return string(k), nil
		}
	}
}

func taxAccountLabel(rowPrefix string) getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		aid, err := getters.Key[*uint](rowPrefix + ".AccountID")(ctx)
		if err != nil {
			return "", err
		}
		if aid == nil || *aid == 0 {
			return "—", nil
		}
		name, err := getters.Key[string](rowPrefix + ".Account.Name")(ctx)
		if err != nil || strings.TrimSpace(name) == "" {
			return fmt.Sprintf("#%d", *aid), nil
		}
		code, err := getters.Key[int](rowPrefix + ".Account.Code")(ctx)
		if err == nil && code != 0 {
			return fmt.Sprintf("%d — %s", code, name), nil
		}
		return fmt.Sprintf("%s (#%d)", name, *aid), nil
	}
}

func pageEntriesTaxMenus() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "finance_taxes.TaxDetailMenu", Value: &components.SidebarMenu{
			Title: getters.Format("%s", getters.Any(getters.Key[string]("tax.Name"))),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("All taxes"),
				Url:   lamu.RoutePath("finance_taxes.DefaultRoute", nil),
			},
			Children: []components.PageInterface{
				&components.SidebarMenuItem{
					Title: getters.Static("Detail"),
					Url: lamu.RoutePath("finance_taxes.TaxDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("tax.ID")),
					}),
				},
				&components.SidebarMenuItem{
					Page:  components.Page{Roles: []string{"superuser"}},
					Title: getters.Static("Edit"),
					Url: lamu.RoutePath("finance_taxes.TaxUpdateRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("tax.ID")),
					}),
				},
			},
		}},
	}
}

func taxCreateFormInputs() []components.PageInterface {
	return []components.PageInterface{
		&components.ContainerError{
			Error: getters.Key[error]("$error.Name"),
			Children: []components.PageInterface{
				&components.InputText{Name: "Name", Label: "Name", Required: true},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.TaxType"),
			Children: []components.PageInterface{
				&components.InputSelect[TaxKind]{
					Name:     "TaxType",
					Label:    "Type",
					Required: true,
					Choices:  taxKindChoices,
					Getter:   taxKindSelectGetter("$in.TaxType"),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Percentage"),
			Children: []components.PageInterface{
				&components.InputPointsDecimal{
					Label:    "Percentage",
					Name:     "Percentage",
					Required: true,
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.AccountID"),
			Children: []components.PageInterface{
				&components.InputForeignKey[finance_accounts.Account]{
					Label:       "Account",
					Name:        "AccountID",
					Url:         lamu.RoutePath("finance_accounts.AccountSelectRoute", nil),
					Display:     getters.Format("%s (#%d)", getters.Any(getters.Key[string]("$in.Name")), getters.Any(getters.Key[uint]("$in.ID"))),
					Placeholder: "Select…",
					Required:    true,
					Getter:      getters.Association[finance_accounts.Account, *uint](getters.Key[*uint]("$in.AccountID")),
				},
			},
		},
	}
}

func taxUpdateFormInputs() []components.PageInterface {
	return []components.PageInterface{
		&components.ContainerError{
			Error: getters.Key[error]("$error.Name"),
			Children: []components.PageInterface{
				&components.InputText{Name: "Name", Label: "Name", Required: true, Getter: getters.Key[string]("$in.Name")},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.TaxType"),
			Children: []components.PageInterface{
				&components.InputSelect[TaxKind]{
					Name:     "TaxType",
					Label:    "Type",
					Required: true,
					Choices:  taxKindChoices,
					Getter:   taxKindSelectGetter("$in.TaxType"),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Percentage"),
			Children: []components.PageInterface{
				&components.InputPointsDecimal{
					Label:    "Percentage",
					Name:     "Percentage",
					Required: true,
					Getter:   taxDecimalGetter("$in.Percentage"),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.AccountID"),
			Children: []components.PageInterface{
				&components.InputForeignKey[finance_accounts.Account]{
					Label:       "Account",
					Name:        "AccountID",
					Url:         lamu.RoutePath("finance_accounts.AccountSelectRoute", nil),
					Display:     getters.Format("%s (#%d)", getters.Any(getters.Key[string]("$in.Name")), getters.Any(getters.Key[uint]("$in.ID"))),
					Placeholder: "Select…",
					Required:    true,
					Getter:      getters.Association[finance_accounts.Account, *uint](getters.Key[*uint]("$in.AccountID")),
				},
			},
		},
	}
}

func pageEntriesTaxPages() []registry.Pair[string, components.PageInterface] {
	createName := getters.Static("finance_taxes.TaxCreateForm")
	updateName := getters.Static("finance_taxes.TaxUpdateForm")
	deleteName := getters.Static("finance_taxes.TaxDeleteForm")

	return []registry.Pair[string, components.PageInterface]{
		{Key: "finance_taxes.TaxTable", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_accounts.MainMenu"}},
			Children: []components.PageInterface{
				&components.DataTable[Tax]{
					UID:     "finance-tax-table",
					Classes: "w-full",
					Data:    getters.Key[components.ObjectList[Tax]]("taxes"),
					Actions: []components.PageInterface{
						&components.TableButtonCreate{
							Link: lamu.RoutePath("finance_taxes.TaxCreateRoute", nil),
							Page: components.Page{Roles: []string{"superuser"}},
						},
					},
					RowAttr: getters.RowAttrNavigate(lamu.RoutePath("finance_taxes.TaxDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("$row.ID")),
					})),
					Columns: []components.TableColumn{
						{Label: "Name", Name: "Name", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Name")},
						}},
						{Label: "Type", Name: "TaxType", Children: []components.PageInterface{
							&components.FieldText{Getter: taxKindLabel("$row.TaxType")},
						}},
						{Label: "Percentage", Name: "Percentage", Children: []components.PageInterface{
							&components.FieldText{Getter: taxDecimalStringGetter("$row.Percentage")},
						}},
						{Label: "Account", Name: "AccountID", Children: []components.PageInterface{
							&components.FieldText{Getter: taxAccountLabel("$row")},
						}},
					},
				},
			},
		}},
		{Key: "finance_taxes.TaxCreateForm", Value: &components.ShellScaffold{
			Page:    components.Page{Roles: []string{"superuser"}},
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_accounts.MainMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name:      createName,
					ActionURL: lamu.RoutePath("finance_taxes.TaxCreateRoute", nil),
					Children: []components.PageInterface{
						&components.FormComponent[Tax]{
							Attr:          getters.FormBubbling(createName),
							Title:         "Create tax",
							Subtitle:      "Name, type, percentage, and GL account",
							ChildrenInput: taxCreateFormInputs(),
							ChildrenAction: []components.PageInterface{
								&components.ButtonSubmit{Label: "Save"},
							},
						},
					},
				},
			},
		}},
		{Key: "finance_taxes.TaxUpdateForm", Value: &components.ShellScaffold{
			Page:    components.Page{Roles: []string{"superuser"}},
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_taxes.TaxDetailMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name: updateName,
					ActionURL: lamu.RoutePath("finance_taxes.TaxUpdateRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("tax.ID")),
					}),
					Children: []components.PageInterface{
						&components.FormComponent[Tax]{
							Getter:        getters.Key[Tax]("tax"),
							Attr:          getters.FormBubbling(updateName),
							Title:         "Edit tax",
							Subtitle:      "Update name, type, percentage, and GL account",
							ChildrenInput: taxUpdateFormInputs(),
							ChildrenAction: []components.PageInterface{
								&components.ContainerRow{
									Classes: "flex flex-wrap justify-between gap-2 mt-2 items-center",
									Children: []components.PageInterface{
										&components.ContainerRow{
											Classes: "flex justify-end gap-2",
											Children: []components.PageInterface{
												&components.ButtonSubmit{Label: "Update"},
												&components.ButtonModalForm{
													Page:        components.Page{Roles: []string{"superuser"}},
													Label:       "Delete",
													Icon:        "trash",
													Name:        deleteName,
													Url:         lamu.RoutePath("finance_taxes.TaxDeleteRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("tax.ID"))}),
													FormPostURL: lamu.RoutePath("finance_taxes.TaxDeleteRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("tax.ID"))}),
													ModalUID:    "finance-tax-delete-modal",
													Classes:     "btn-error",
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
		{Key: "finance_taxes.TaxDeleteForm", Value: &components.Modal{
			Page: components.Page{Roles: []string{"superuser"}},
			UID:  "finance-tax-delete-modal",
			Children: []components.PageInterface{
				&components.DeleteConfirmation{
					Title:   "Delete tax?",
					Message: "This permanently removes the tax record.",
					Attr:    getters.FormBubbling(getters.Key[string]("$get.name")),
				},
			},
		}},
		{Key: "finance_taxes.TaxDetail", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_taxes.TaxDetailMenu"}},
			Children: []components.PageInterface{
				&components.Detail[Tax]{
					Getter: getters.Key[Tax]("tax"),
					Children: []components.PageInterface{
						&components.ContainerColumn{
							Classes: "p-4",
							Children: []components.PageInterface{
								&components.LabelInline{Title: "Name", Children: []components.PageInterface{
									&components.FieldText{Getter: getters.Key[string]("$in.Name")},
								}},
								&components.LabelInline{Title: "Type", Children: []components.PageInterface{
									&components.FieldText{Getter: taxKindLabel("$in.TaxType")},
								}},
								&components.LabelInline{Title: "Percentage", Children: []components.PageInterface{
									&components.FieldText{Getter: taxDecimalStringGetter("$in.Percentage")},
								}},
								&components.LabelInline{Title: "Account", Children: []components.PageInterface{
									&components.FieldText{Getter: taxAccountLabel("$in")},
								}},
							},
						},
					},
				},
			},
		}},
	}
}
