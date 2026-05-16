package p_uniquity_finance_products

import (
	"context"
	"strings"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/fields"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
	finance_accounts "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_accounts"
	finance_taxes "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_taxes"
)

const financeAccountsMainMenuProductsLinkKey = "finance_products.FinanceAccountsMainMenuLink"

func patchFinanceAccountsMainMenuForProducts(page components.PageInterface) components.PageInterface {
	menu, ok := page.(*components.SidebarMenu)
	if !ok {
		panic("p_uniquity_finance_products: finance_accounts.MainMenu must be *components.SidebarMenu")
	}
	for _, ch := range menu.Children {
		if item, ok := ch.(*components.SidebarMenuItem); ok && item.GetKey() == financeAccountsMainMenuProductsLinkKey {
			return menu
		}
	}
	newChildren := append([]components.PageInterface{}, menu.Children...)
	newChildren = append(newChildren, &components.SidebarMenuItem{
		Page:  components.Page{Key: financeAccountsMainMenuProductsLinkKey, Roles: []string{"superuser"}},
		Title: getters.Static("Products"),
		Url:   lamu.RoutePath("finance_products.DefaultRoute", nil),
		Icon:  "cube",
	})
	cloned := *menu
	cloned.Children = newChildren
	return &cloned
}

func pluginPages() lamu.PluginFeatures[components.PageInterface] {
	e := pageEntriesProductMenus()
	e = append(e, pageEntriesProductPages()...)
	return lamu.PluginFeatures[components.PageInterface]{
		Entries: e,
		Patches: []registry.Pair[string, func(components.PageInterface) components.PageInterface]{
			{Key: "finance_accounts.MainMenu", Value: patchFinanceAccountsMainMenuForProducts},
		},
	}
}

func productDecimalStringGetter(ctxKey string) getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		pd, err := getters.Key[fields.DecimalSix](ctxKey)(ctx)
		if err != nil {
			return "", err
		}
		return pd.String(), nil
	}
}

func productDecimalGetter(ctxKey string) getters.Getter[fields.DecimalSix] {
	return func(ctx context.Context) (fields.DecimalSix, error) {
		return getters.Key[fields.DecimalSix](ctxKey)(ctx)
	}
}

func pageEntriesProductMenus() []registry.Pair[string, components.PageInterface] {
	return []registry.Pair[string, components.PageInterface]{
		{Key: "finance_products.ProductDetailMenu", Value: &components.SidebarMenu{
			Title: getters.Format("%s", getters.Any(getters.Key[string]("product.Name"))),
			Back: &components.SidebarMenuItem{
				Title: getters.Static("All products"),
				Url:   lamu.RoutePath("finance_products.DefaultRoute", nil),
			},
			Children: []components.PageInterface{
				&components.SidebarMenuItem{
					Title: getters.Static("Detail"),
					Url: lamu.RoutePath("finance_products.ProductDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("product.ID")),
					}),
				},
				&components.SidebarMenuItem{
					Page:  components.Page{Roles: []string{"superuser"}},
					Title: getters.Static("Edit"),
					Url: lamu.RoutePath("finance_products.ProductUpdateRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("product.ID")),
					}),
				},
			},
		}},
	}
}

func productCreateFormInputs() []components.PageInterface {
	return []components.PageInterface{
		&components.ContainerError{
			Error: getters.Key[error]("$error.Name"),
			Children: []components.PageInterface{
				&components.InputText{Name: "Name", Label: "Name", Required: true},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Taxes"),
			Children: []components.PageInterface{
				&components.InputManyToMany[finance_taxes.Tax]{
					Label:       "Taxes",
					Name:        "Taxes",
					Display:     getters.Key[string]("$in.Name"),
					Url:         lamu.RoutePath("finance_taxes.TaxMultiSelectRoute", nil),
					Placeholder: "Select taxes…",
					Classes:     "w-full",
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.BaseCost"),
			Children: []components.PageInterface{
				&components.InputPointsDecimal{
					Label:    "Base cost",
					Name:     "BaseCost",
					Required: true,
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.SalesPrice"),
			Children: []components.PageInterface{
				&components.InputPointsDecimal{
					Label:    "Sales price",
					Name:     "SalesPrice",
					Required: true,
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.HSNCode"),
			Children: []components.PageInterface{
				&components.InputNumber[int64]{Label: "HSN code", Name: "HSNCode", Required: true},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.InventoryAccountID"),
			Children: []components.PageInterface{
				&components.InputForeignKey[finance_accounts.Account]{
					Label:       "Inventory account",
					Name:        "InventoryAccountID",
					Url:         lamu.RoutePath("finance_accounts.AccountSelectRoute", nil),
					Display:     getters.Format("%s (#%d)", getters.Any(getters.Key[string]("$in.Name")), getters.Any(getters.Key[uint]("$in.ID"))),
					Placeholder: "Select…",
					Getter:      getters.Association[finance_accounts.Account, uint](getters.Key[uint]("$in.InventoryAccountID")),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.CostOfSalesAcctID"),
			Children: []components.PageInterface{
				&components.InputForeignKey[finance_accounts.Account]{
					Label:       "Cost of sales account",
					Name:        "CostOfSalesAcctID",
					Url:         lamu.RoutePath("finance_accounts.AccountSelectRoute", nil),
					Display:     getters.Format("%s (#%d)", getters.Any(getters.Key[string]("$in.Name")), getters.Any(getters.Key[uint]("$in.ID"))),
					Placeholder: "Select…",
					Getter:      getters.Association[finance_accounts.Account, uint](getters.Key[uint]("$in.CostOfSalesAcctID")),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.InputTaxAccountID"),
			Children: []components.PageInterface{
				&components.InputForeignKey[finance_accounts.Account]{
					Label:       "Input tax (ITC) account",
					Name:        "InputTaxAccountID",
					Url:         lamu.RoutePath("finance_accounts.AccountSelectRoute", nil),
					Display:     getters.Format("%s (#%d)", getters.Any(getters.Key[string]("$in.Name")), getters.Any(getters.Key[uint]("$in.ID"))),
					Placeholder: "Optional",
					Getter:      getters.Association[finance_accounts.Account, uint](getters.Key[uint]("$in.InputTaxAccountID")),
				},
			},
		},
	}
}

func productUpdateFormInputs() []components.PageInterface {
	return []components.PageInterface{
		&components.ContainerError{
			Error: getters.Key[error]("$error.Name"),
			Children: []components.PageInterface{
				&components.InputText{Name: "Name", Label: "Name", Required: true, Getter: getters.Key[string]("$in.Name")},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Taxes"),
			Children: []components.PageInterface{
				&components.InputManyToMany[finance_taxes.Tax]{
					Label:       "Taxes",
					Name:        "Taxes",
					Getter:      getters.Key[[]finance_taxes.Tax]("$in.Taxes"),
					Display:     getters.Key[string]("$in.Name"),
					Url:         lamu.RoutePath("finance_taxes.TaxMultiSelectRoute", nil),
					Placeholder: "Select taxes…",
					Classes:     "w-full",
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.BaseCost"),
			Children: []components.PageInterface{
				&components.InputPointsDecimal{
					Label:    "Base cost",
					Name:     "BaseCost",
					Required: true,
					Getter:   productDecimalGetter("$in.BaseCost"),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.SalesPrice"),
			Children: []components.PageInterface{
				&components.InputPointsDecimal{
					Label:    "Sales price",
					Name:     "SalesPrice",
					Required: true,
					Getter:   productDecimalGetter("$in.SalesPrice"),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.HSNCode"),
			Children: []components.PageInterface{
				&components.InputNumber[int64]{Label: "HSN code", Name: "HSNCode", Required: true, Getter: getters.Key[int64]("$in.HSNCode")},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.InventoryAccountID"),
			Children: []components.PageInterface{
				&components.InputForeignKey[finance_accounts.Account]{
					Label:       "Inventory account",
					Name:        "InventoryAccountID",
					Url:         lamu.RoutePath("finance_accounts.AccountSelectRoute", nil),
					Display:     getters.Format("%s (#%d)", getters.Any(getters.Key[string]("$in.Name")), getters.Any(getters.Key[uint]("$in.ID"))),
					Placeholder: "Select…",
					Getter:      getters.Association[finance_accounts.Account, uint](getters.Key[uint]("$in.InventoryAccountID")),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.CostOfSalesAcctID"),
			Children: []components.PageInterface{
				&components.InputForeignKey[finance_accounts.Account]{
					Label:       "Cost of sales account",
					Name:        "CostOfSalesAcctID",
					Url:         lamu.RoutePath("finance_accounts.AccountSelectRoute", nil),
					Display:     getters.Format("%s (#%d)", getters.Any(getters.Key[string]("$in.Name")), getters.Any(getters.Key[uint]("$in.ID"))),
					Placeholder: "Select…",
					Getter:      getters.Association[finance_accounts.Account, uint](getters.Key[uint]("$in.CostOfSalesAcctID")),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.InputTaxAccountID"),
			Children: []components.PageInterface{
				&components.InputForeignKey[finance_accounts.Account]{
					Label:       "Input tax (ITC) account",
					Name:        "InputTaxAccountID",
					Url:         lamu.RoutePath("finance_accounts.AccountSelectRoute", nil),
					Display:     getters.Format("%s (#%d)", getters.Any(getters.Key[string]("$in.Name")), getters.Any(getters.Key[uint]("$in.ID"))),
					Placeholder: "Optional",
					Getter:      getters.Association[finance_accounts.Account, uint](getters.Key[uint]("$in.InputTaxAccountID")),
				},
			},
		},
	}
}

func pageEntriesProductPages() []registry.Pair[string, components.PageInterface] {
	createName := getters.Static("finance_products.ProductCreateForm")
	updateName := getters.Static("finance_products.ProductUpdateForm")
	deleteName := getters.Static("finance_products.ProductDeleteForm")

	return []registry.Pair[string, components.PageInterface]{
		{Key: "finance_products.ProductTable", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_accounts.MainMenu"}},
			Children: []components.PageInterface{
				&components.DataTable[Product]{
					UID:     "finance-product-table",
					Classes: "w-full",
					Data:    getters.Key[components.ObjectList[Product]]("products"),
					Actions: []components.PageInterface{
						&components.TableButtonCreate{
							Link: lamu.RoutePath("finance_products.ProductCreateRoute", nil),
							Page: components.Page{Roles: []string{"superuser"}},
						},
					},
					RowAttr: getters.RowAttrNavigate(lamu.RoutePath("finance_products.ProductDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("$row.ID")),
					})),
					Columns: []components.TableColumn{
						{Label: "Name", Name: "Name", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Name")},
						}},
						{Label: "Base cost", Name: "BaseCost", Children: []components.PageInterface{
							&components.FieldText{Getter: productDecimalStringGetter("$row.BaseCost")},
						}},
						{Label: "Sales price", Name: "SalesPrice", Children: []components.PageInterface{
							&components.FieldText{Getter: productDecimalStringGetter("$row.SalesPrice")},
						}},
						{Label: "HSN code", Name: "HSNCode", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Format("%d", getters.Any(getters.Key[int64]("$row.HSNCode")))},
						}},
					},
				},
			},
		}},
		{Key: "finance_products.ProductCreateForm", Value: &components.ShellScaffold{
			Page:    components.Page{Roles: []string{"superuser"}},
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_accounts.MainMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name:      createName,
					ActionURL: lamu.RoutePath("finance_products.ProductCreateRoute", nil),
					Children: []components.PageInterface{
						&components.FormComponent[Product]{
							Attr:          getters.FormBubbling(createName),
							Title:         "Create product",
							Subtitle:      "Pricing, HSN, and applicable taxes",
							ChildrenInput: productCreateFormInputs(),
							ChildrenAction: []components.PageInterface{
								&components.ButtonSubmit{Label: "Save"},
							},
						},
					},
				},
			},
		}},
		{Key: "finance_products.ProductUpdateForm", Value: &components.ShellScaffold{
			Page:    components.Page{Roles: []string{"superuser"}},
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_products.ProductDetailMenu"}},
			Children: []components.PageInterface{
				&components.FormListenBoostedPost{
					Name: updateName,
					ActionURL: lamu.RoutePath("finance_products.ProductUpdateRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("product.ID")),
					}),
					Children: []components.PageInterface{
						&components.FormComponent[Product]{
							Getter:        getters.Key[Product]("product"),
							Attr:          getters.FormBubbling(updateName),
							Title:         "Edit product",
							Subtitle:      "Update pricing, HSN, and taxes",
							ChildrenInput: productUpdateFormInputs(),
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
													Url:         lamu.RoutePath("finance_products.ProductDeleteRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("product.ID"))}),
													FormPostURL: lamu.RoutePath("finance_products.ProductDeleteRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("product.ID"))}),
													ModalUID:    "finance-product-delete-modal",
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
		{Key: "finance_products.ProductDeleteForm", Value: &components.Modal{
			Page: components.Page{Roles: []string{"superuser"}},
			UID:  "finance-product-delete-modal",
			Children: []components.PageInterface{
				&components.DeleteConfirmation{
					Title:   "Delete product?",
					Message: "This permanently removes the product record.",
					Attr:    getters.FormBubbling(getters.Key[string]("$get.name")),
				},
			},
		}},
		{Key: "finance_products.ProductDetail", Value: &components.ShellScaffold{
			Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "finance_products.ProductDetailMenu"}},
			Children: []components.PageInterface{
				&components.Detail[Product]{
					Getter: getters.Key[Product]("product"),
					Children: []components.PageInterface{
						&components.ContainerColumn{
							Classes: "p-4",
							Children: []components.PageInterface{
								&components.LabelInline{Title: "Name", Children: []components.PageInterface{
									&components.FieldText{Getter: getters.Key[string]("$in.Name")},
								}},
								&components.FieldManyToMany[finance_taxes.Tax]{
									Label:   "Taxes",
									Getter:  getters.Key[[]finance_taxes.Tax]("$in.Taxes"),
									Display: getters.Key[string]("$in.Name"),
									Link: lamu.RoutePath("finance_taxes.TaxDetailRoute", map[string]getters.Getter[any]{
										"id": getters.Any(getters.Key[uint]("$in.ID")),
									}),
									Classes: "w-full",
								},
								&components.LabelInline{Title: "Base cost", Children: []components.PageInterface{
									&components.FieldText{Getter: productDecimalStringGetter("$in.BaseCost")},
								}},
								&components.LabelInline{Title: "Sales price", Children: []components.PageInterface{
									&components.FieldText{Getter: productDecimalStringGetter("$in.SalesPrice")},
								}},
								&components.LabelInline{Title: "HSN code", Children: []components.PageInterface{
									&components.FieldText{Getter: getters.Format("%d", getters.Any(getters.Key[int64]("$in.HSNCode")))},
								}},
								&components.LabelInline{Title: "Inventory GL", Children: []components.PageInterface{
									&components.FieldText{Getter: getters.Format("%s", getters.Any(accountNameOrDash("$in.InventoryAccount.Name")))},
								}},
								&components.LabelInline{Title: "COGS GL", Children: []components.PageInterface{
									&components.FieldText{Getter: getters.Format("%s", getters.Any(accountNameOrDash("$in.CostOfSalesAccount.Name")))},
								}},
								&components.LabelInline{Title: "Input tax GL", Children: []components.PageInterface{
									&components.FieldText{Getter: getters.Format("%s", getters.Any(accountNameOrDash("$in.InputTaxAccount.Name")))},
								}},
							},
						},
					},
				},
			},
		}},
	}
}

func accountNameOrDash(ctxKey string) getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		s, err := getters.Key[string](ctxKey)(ctx)
		if err != nil || strings.TrimSpace(s) == "" {
			return "—", nil
		}
		return s, nil
	}
}
