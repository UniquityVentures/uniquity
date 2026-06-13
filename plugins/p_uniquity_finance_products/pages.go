package p_uniquity_finance_products

import (
	"context"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/fields"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
	finance_taxes "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_taxes"
)

const financeAccountsMainMenuProductsLinkKey = "finance_products.FinanceAccountsMainMenuLink"

var productTypeChoiceList = []registry.Pair[ProductType, string]{
	{Key: ProductTypeGoods, Value: "Goods"},
	{Key: ProductTypeServices, Value: "Services"},
	{Key: ProductTypeBoth, Value: "Both"},
}

var productTypeChoices = getters.Static(productTypeChoiceList)

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
	e = append(e, pageEntriesProductFkSelectPages()...)
	return lamu.PluginFeatures[components.PageInterface]{
		Entries: e,
		Patches: []registry.Pair[string, func(components.PageInterface) components.PageInterface]{
			{Key: "finance_accounts.MainMenu", Value: patchFinanceAccountsMainMenuForProducts},
			{Key: "finance_accounts.AccountingPreferencesForm", Value: patchAccountingPreferencesForm},
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

func productFormInputs() []components.PageInterface {
	return []components.PageInterface{
		&components.ContainerError{
			Error: getters.Key[error]("$error.Name"),
			Children: []components.PageInterface{
				&components.InputText{Name: "Name", Label: "Name", Required: true, Getter: getters.Key[string]("$in.Name")},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Type"),
			Children: []components.PageInterface{
				&components.InputSelect[ProductType]{
					Name:     "Type",
					Label:    "Type",
					Required: true,
					Choices:  productTypeChoices,
					Getter:   registry.PairFromGetter(getters.Key[ProductType]("$in.Type"), productTypeChoiceList),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Reference"),
			Children: []components.PageInterface{
				&components.InputText{Name: "Reference", Label: "Reference", Required: true, Getter: getters.Key[string]("$in.Reference")},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Remarks"),
			Children: []components.PageInterface{
				&components.InputTextarea{Name: "Remarks", Label: "Remarks", Getter: getters.Key[string]("$in.Remarks"), Rows: 4},
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
	}
}

func productCreateFormInputs() []components.PageInterface {
	return productFormInputs()
}

func productUpdateFormInputs() []components.PageInterface {
	return productFormInputs()
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
						{Label: "Type", Name: "Type", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Format("%s", getters.Any(getters.Key[ProductType]("$row.Type")))},
						}},
						{Label: "Reference", Name: "Reference", Children: []components.PageInterface{
							&components.FieldText{Getter: getters.Key[string]("$row.Reference")},
						}},
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
								&components.LabelInline{Title: "Type", Children: []components.PageInterface{
									&components.FieldText{Getter: getters.Format("%s", getters.Any(getters.Key[ProductType]("$in.Type")))},
								}},
								&components.LabelInline{Title: "Reference", Children: []components.PageInterface{
									&components.FieldText{Getter: getters.Key[string]("$in.Reference")},
								}},
								&components.LabelInline{Title: "Remarks", Children: []components.PageInterface{
									&components.FieldText{Getter: getters.Key[string]("$in.Remarks")},
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
							},
						},
					},
				},
			},
		}},
	}
}
