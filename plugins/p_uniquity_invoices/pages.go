package p_uniquity_invoices

import (
	"context"
	"time"

	acct "github.com/UniquityVentures/uniquity/plugins/p_uniquity_accounting"
	currencies "github.com/UniquityVentures/uniquity/plugins/p_uniquity_currencies"
	ent "github.com/UniquityVentures/uniquity/plugins/p_uniquity_entities"
	prod "github.com/UniquityVentures/uniquity/plugins/p_uniquity_products"
	tax "github.com/UniquityVentures/uniquity/plugins/p_uniquity_tax_rates"
	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/fields"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"
)

func invoiceKindChoices() []registry.Pair[string, string] {
	return []registry.Pair[string, string]{
		{Key: InvoiceTypeOutInvoice, Value: "Customer invoice"},
		{Key: InvoiceTypeOutRefund, Value: "Customer credit note"},
		{Key: InvoiceTypeInInvoice, Value: "Vendor bill"},
		{Key: InvoiceTypeInRefund, Value: "Vendor credit note"},
	}
}

func invoiceStateChoices() []registry.Pair[string, string] {
	return []registry.Pair[string, string]{
		{Key: InvoiceStateDraft, Value: "Draft"},
		{Key: InvoiceStatePosted, Value: "Posted"},
		{Key: InvoiceStatePaid, Value: "Paid"},
		{Key: InvoiceStateCancelled, Value: "Cancelled"},
	}
}

func invoiceTypePairGetter(g getters.Getter[string]) getters.Getter[registry.Pair[string, string]] {
	labels := map[string]string{}
	for _, p := range invoiceKindChoices() {
		labels[p.Key] = p.Value
	}
	return func(ctx context.Context) (registry.Pair[string, string], error) {
		k, err := g(ctx)
		if err != nil {
			return registry.Pair[string, string]{}, err
		}
		lab := labels[k]
		if lab == "" {
			lab = k
		}
		return registry.Pair[string, string]{Key: k, Value: lab}, nil
	}
}

func invoiceStatePairGetter(g getters.Getter[string]) getters.Getter[registry.Pair[string, string]] {
	labels := map[string]string{}
	for _, p := range invoiceStateChoices() {
		labels[p.Key] = p.Value
	}
	return func(ctx context.Context) (registry.Pair[string, string], error) {
		k, err := g(ctx)
		if err != nil {
			return registry.Pair[string, string]{}, err
		}
		lab := labels[k]
		if lab == "" {
			lab = k
		}
		return registry.Pair[string, string]{Key: k, Value: lab}, nil
	}
}

func invoiceFormInputs() []components.PageInterface {
	return []components.PageInterface{
		&components.ContainerError{
			Error: getters.Key[error]("$error.EntityID"),
			Children: []components.PageInterface{
				&components.InputForeignKey[ent.Entity]{
					Name:        "EntityID",
					Label:       "Entity",
					Url:         lamu.RoutePath("entities.EntitySelectRoute", nil),
					Display:     getters.Key[string]("$in.Entity.Name"),
					Placeholder: "Select entity…",
					Required:    true,
					Getter:      getters.Association[ent.Entity, uint](getters.Key[uint]("$in.EntityID")),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.PartnerID"),
			Children: []components.PageInterface{
				&components.InputForeignKey[Contact]{
					Name:        "PartnerID",
					Label:       "Partner (contact)",
					Url:         lamu.RoutePath("invoices.ContactSelectRoute", nil),
					Display:     getters.Key[string]("$in.Partner.Name"),
					Placeholder: "Select contact…",
					Required:    true,
					Getter:      getters.Association[Contact, uint](getters.Key[uint]("$in.PartnerID")),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.JournalID"),
			Children: []components.PageInterface{
				&components.InputForeignKey[acct.Journal]{
					Name:        "JournalID",
					Label:       "Journal",
					Url:         lamu.RoutePath("accounting.JournalEntryJournalSelectRoute", nil),
					Display:     getters.Key[string]("$in.Journal.Name"),
					Placeholder: "Select journal…",
					Required:    true,
					Getter:      getters.Association[acct.Journal, uint](getters.Key[uint]("$in.JournalID")),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.InvoiceType"),
			Children: []components.PageInterface{
				&components.InputSelect[string]{
					Label:    "Document type",
					Name:     "InvoiceType",
					Required: true,
					Choices:  getters.Static(invoiceKindChoices()),
					Getter:   invoiceTypePairGetter(getters.Key[string]("$in.InvoiceType")),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.State"),
			Children: []components.PageInterface{
				&components.InputSelect[string]{
					Label:    "State",
					Name:     "State",
					Required: true,
					Choices:  getters.Static(invoiceStateChoices()),
					Getter:   invoiceStatePairGetter(getters.Key[string]("$in.State")),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Reference"),
			Children: []components.PageInterface{
				&components.InputText{
					Label:  "Reference",
					Name:   "Reference",
					Getter: getters.Key[string]("$in.Reference"),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Number"),
			Children: []components.PageInterface{
				&components.InputText{
					Label:  "Number",
					Name:   "Number",
					Getter: getters.Deref(getters.Key[*string]("$in.Number")),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.InvoiceDate"),
			Children: []components.PageInterface{
				&components.InputDate{
					Label:    "Invoice date",
					Name:     "InvoiceDate",
					Required: true,
					Getter:   getters.Key[time.Time]("$in.InvoiceDate"),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.PaymentTermID"),
			Children: []components.PageInterface{
				&components.InputForeignKey[PaymentTerm]{
					Name:        "PaymentTermID",
					Label:       "Payment term",
					Url:         lamu.RoutePath("invoices.PaymentTermSelectRoute", nil),
					Display:     getters.Key[string]("$in.PaymentTerm.Name"),
					Placeholder: "Optional…",
					Required:    false,
					Getter:      getters.Association[PaymentTerm, *uint](getters.Key[*uint]("$in.PaymentTermID")),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.DueDate"),
			Children: []components.PageInterface{
				&components.InputDate{
					Label:    "Due date",
					Name:     "DueDate",
					Required: false,
					Getter:   getters.Deref(getters.Key[*time.Time]("$in.DueDate")),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.CurrencyID"),
			Children: []components.PageInterface{
				&components.InputForeignKey[currencies.Currency]{
					Name:        "CurrencyID",
					Label:       "Currency",
					Url:         lamu.RoutePath("currencies.CurrencySelectRoute", nil),
					Display:     getters.Format("%s — %s", getters.Any(getters.Key[string]("$in.Currency.Code")), getters.Any(getters.Key[string]("$in.Currency.Name"))),
					Placeholder: "Select currency…",
					Required:    true,
					Getter:      getters.Association[currencies.Currency, uint](getters.Key[uint]("$in.CurrencyID")),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.MoveID"),
			Children: []components.PageInterface{
				&components.InputForeignKey[acct.JournalEntry]{
					Name:        "MoveID",
					Label:       "Posted entry (journal entry)",
					Url:         lamu.RoutePath("accounting.JournalEntryItemJournalEntrySelectRoute", nil),
					Display:     getters.Format("#%d — %s", getters.Any(getters.Key[uint]("$in.Move.ID")), getters.Any(getters.Key[string]("$in.Move.Journal.Name"))),
					Placeholder: "Optional…",
					Required:    false,
					Getter:      getters.Association[acct.JournalEntry, *uint](getters.Key[*uint]("$in.MoveID")),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.AmountUntaxed"),
			Children: []components.PageInterface{
				&components.InputPointsDecimal{
					Label:    "Amount untaxed",
					Name:     "AmountUntaxed",
					Required: true,
					Getter:   getters.Key[fields.DecimalSix]("$in.AmountUntaxed"),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.AmountTax"),
			Children: []components.PageInterface{
				&components.InputPointsDecimal{
					Label:    "Amount tax",
					Name:     "AmountTax",
					Required: true,
					Getter:   getters.Key[fields.DecimalSix]("$in.AmountTax"),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.AmountTotal"),
			Children: []components.PageInterface{
				&components.InputPointsDecimal{
					Label:    "Amount total",
					Name:     "AmountTotal",
					Required: true,
					Getter:   getters.Key[fields.DecimalSix]("$in.AmountTotal"),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.AmountResidual"),
			Children: []components.PageInterface{
				&components.InputPointsDecimal{
					Label:    "Amount residual",
					Name:     "AmountResidual",
					Required: true,
					Getter:   getters.Key[fields.DecimalSix]("$in.AmountResidual"),
				},
			},
		},
	}
}

func invoiceLineFormInputs() []components.PageInterface {
	return []components.PageInterface{
		&components.ContainerError{
			Error: getters.Key[error]("$error.ProductID"),
			Children: []components.PageInterface{
				&components.InputForeignKey[prod.Product]{
					Name:        "ProductID",
					Label:       "Product",
					Url:         lamu.RoutePath("products.ProductSelectRoute", nil),
					Display:     getters.Format("%s — %s", getters.Any(getters.Key[string]("$in.Product.Code")), getters.Any(getters.Key[string]("$in.Product.Name"))),
					Placeholder: "Optional…",
					Required:    false,
					Getter:      getters.Association[prod.Product, *uint](getters.Key[*uint]("$in.ProductID")),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Label"),
			Children: []components.PageInterface{
				&components.InputText{
					Label:  "Label",
					Name:   "Label",
					Getter: getters.Key[string]("$in.Label"),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Quantity"),
			Children: []components.PageInterface{
				&components.InputPointsDecimal{
					Label:    "Quantity",
					Name:     "Quantity",
					Required: true,
					Getter:   getters.Key[fields.DecimalSix]("$in.Quantity"),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.PriceUnit"),
			Children: []components.PageInterface{
				&components.InputPointsDecimal{
					Label:    "Unit price",
					Name:     "PriceUnit",
					Required: true,
					Getter:   getters.Key[fields.DecimalSix]("$in.PriceUnit"),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Discount"),
			Children: []components.PageInterface{
				&components.InputPointsDecimal{
					Label:    "Discount %",
					Name:     "Discount",
					Required: true,
					Getter:   getters.Key[fields.DecimalSix]("$in.Discount"),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.PriceSubtotal"),
			Children: []components.PageInterface{
				&components.InputPointsDecimal{
					Label:    "Line subtotal",
					Name:     "PriceSubtotal",
					Required: true,
					Getter:   getters.Key[fields.DecimalSix]("$in.PriceSubtotal"),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.AccountID"),
			Children: []components.PageInterface{
				&components.InputForeignKey[acct.Account]{
					Name:        "AccountID",
					Label:       "Account",
					Url:         lamu.RoutePath("accounting.JournalEntryItemAccountSelectRoute", nil),
					Display:     getters.Format("%s — %s", getters.Any(getters.Key[string]("$in.Account.Code")), getters.Any(getters.Key[string]("$in.Account.Name"))),
					Placeholder: "Select account…",
					Required:    true,
					Getter:      getters.Association[acct.Account, uint](getters.Key[uint]("$in.AccountID")),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Taxes"),
			Children: []components.PageInterface{
				&components.InputManyToMany[tax.TaxRate]{
					Label:       "Taxes",
					Name:        "Taxes",
					Required:    false,
					Getter:      getters.Key[[]tax.TaxRate]("$in.Taxes"),
					Display:     getters.Key[string]("$in.Name"),
					Url:         lamu.RoutePath("tax_rates.TaxRateSelectRoute", nil),
					Placeholder: "Add tax rates…",
				},
			},
		},
	}
}

func pluginPages() lamu.PluginFeatures[components.PageInterface] {
	invCreate := getters.Static("invoices.InvoiceCreateForm")
	invUpdate := getters.Static("invoices.InvoiceUpdateForm")
	invDelete := getters.Static("invoices.InvoiceDeleteForm")
	lineCreate := getters.Static("invoices.InvoiceLineCreateForm")
	lineUpdate := getters.Static("invoices.InvoiceLineUpdateForm")
	lineDelete := getters.Static("invoices.InvoiceLineDeleteForm")

	return lamu.PluginFeatures[components.PageInterface]{
		Entries: []registry.Pair[string, components.PageInterface]{
			{Key: "invoices.MainMenu", Value: &components.SidebarMenu{
				Title: getters.Static("Invoices"),
				Back: &components.SidebarMenuItem{
					Title: getters.Static("Back to Home"),
					Url:   lamu.RoutePath("dashboard.AppsPage", nil),
				},
				Children: []components.PageInterface{
					&components.SidebarMenuItem{
						Page:  components.Page{Roles: []string{"superuser"}},
						Title: getters.Static("All invoices"),
						Url:   lamu.RoutePath("invoices.DefaultRoute", nil),
						Icon:  "document-text",
					},
					&components.SidebarMenuItem{
						Page:  components.Page{Roles: []string{"superuser"}},
						Title: getters.Static("New invoice"),
						Url:   lamu.RoutePath("invoices.InvoiceCreateRoute", nil),
						Icon:  "plus",
					},
				},
			}},
			{Key: "invoices.InvoiceDetailMenu", Value: &components.SidebarMenu{
				Title: getters.Format("Invoice #%d", getters.Any(getters.Key[uint]("invoice.ID"))),
				Back: &components.SidebarMenuItem{
					Title: getters.Static("All invoices"),
					Url:   lamu.RoutePath("invoices.DefaultRoute", nil),
				},
				Children: []components.PageInterface{
					&components.SidebarMenuItem{
						Title: getters.Static("Detail"),
						Url: lamu.RoutePath("invoices.InvoiceDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("invoice.ID")),
						}),
					},
					&components.SidebarMenuItem{
						Page:  components.Page{Roles: []string{"superuser"}},
						Title: getters.Static("Edit"),
						Url: lamu.RoutePath("invoices.InvoiceUpdateRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("invoice.ID")),
						}),
					},
				},
			}},
			{Key: "invoices.InvoiceLineDetailMenu", Value: &components.SidebarMenu{
				Title: getters.Format("Line #%d", getters.Any(getters.Key[uint]("invoiceLine.ID"))),
				Back: &components.SidebarMenuItem{
					Title: getters.Static("Back to invoice"),
					Url: lamu.RoutePath("invoices.InvoiceDetailRoute", map[string]getters.Getter[any]{
						"id": getters.Any(getters.Key[uint]("invoiceLine.InvoiceID")),
					}),
				},
				Children: []components.PageInterface{
					&components.SidebarMenuItem{
						Title: getters.Static("Detail"),
						Url: lamu.RoutePath("invoices.InvoiceLineDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("invoiceLine.ID")),
						}),
					},
					&components.SidebarMenuItem{
						Page:  components.Page{Roles: []string{"superuser"}},
						Title: getters.Static("Edit"),
						Url: lamu.RoutePath("invoices.InvoiceLineUpdateRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("invoiceLine.ID")),
						}),
					},
				},
			}},
			{Key: "invoices.ContactSelectionTable", Value: &components.Modal{
				UID: "invoices-contact-select-modal",
				Children: []components.PageInterface{
					&components.DataTable[Contact]{
						UID:   "invoices-contact-select-table",
						Title: "Select contact",
						Data:  getters.Key[components.ObjectList[Contact]]("contacts"),
						RowAttr: getters.RowAttrSelect("PartnerID",
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
			{Key: "invoices.PaymentTermSelectionTable", Value: &components.Modal{
				UID: "invoices-payment-term-select-modal",
				Children: []components.PageInterface{
					&components.DataTable[PaymentTerm]{
						UID:   "invoices-payment-term-select-table",
						Title: "Select payment term",
						Data:  getters.Key[components.ObjectList[PaymentTerm]]("paymentTerms"),
						RowAttr: getters.RowAttrSelect("PaymentTermID",
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
			{Key: "invoices.InvoiceTable", Value: &components.ShellScaffold{
				Page:    components.Page{Roles: []string{"superuser"}},
				Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "invoices.MainMenu"}},
				Children: []components.PageInterface{
					&components.DataTable[Invoice]{
						UID:     "invoices-table",
						Classes: "w-full",
						Data:    getters.Key[components.ObjectList[Invoice]]("invoices"),
						Actions: []components.PageInterface{
							&components.TableButtonCreate{
								Link: lamu.RoutePath("invoices.InvoiceCreateRoute", nil),
								Page: components.Page{Roles: []string{"superuser"}},
							},
						},
						RowAttr: getters.RowAttrNavigate(lamu.RoutePath("invoices.InvoiceDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("$row.ID")),
						})),
						Columns: []components.TableColumn{
							{Label: "Number", Name: "Number", Children: []components.PageInterface{
								&components.FieldText{Getter: getters.Deref(getters.Key[*string]("$row.Number"))},
							}},
							{Label: "Type", Name: "InvoiceType", Children: []components.PageInterface{
								&components.FieldText{Getter: getters.Key[string]("$row.InvoiceType")},
							}},
							{Label: "State", Name: "State", Children: []components.PageInterface{
								&components.FieldText{Getter: getters.Key[string]("$row.State")},
							}},
							{Label: "Date", Name: "InvoiceDate", Children: []components.PageInterface{
								&components.FieldText{Getter: getters.Format("%v", getters.Any(getters.Key[time.Time]("$row.InvoiceDate")))},
							}},
							{Label: "Partner", Name: "Partner", Children: []components.PageInterface{
								&components.FieldText{Getter: getters.Key[string]("$row.Partner.Name")},
							}},
							{Label: "Total", Name: "AmountTotal", Children: []components.PageInterface{
								&components.FieldText{Getter: getters.Format("%v", getters.Any(getters.Key[fields.DecimalSix]("$row.AmountTotal")))},
							}},
						},
					},
				},
			}},
			{Key: "invoices.InvoiceCreateForm", Value: &components.ShellScaffold{
				Page:    components.Page{Roles: []string{"superuser"}},
				Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "invoices.MainMenu"}},
				Children: []components.PageInterface{
					&components.FormListenBoostedPost{
						Name:      invCreate,
						ActionURL: lamu.RoutePath("invoices.InvoiceCreateRoute", nil),
						Children: []components.PageInterface{
							&components.FormComponent[Invoice]{
								Attr:          getters.FormBubbling(invCreate),
								Title:         "Create invoice",
								Subtitle:      "Customer invoice, vendor bill, or credit note",
								ChildrenInput: invoiceFormInputs(),
								ChildrenAction: []components.PageInterface{
									&components.ButtonSubmit{Label: "Save"},
								},
							},
						},
					},
				},
			}},
			{Key: "invoices.InvoiceUpdateForm", Value: &components.ShellScaffold{
				Page:    components.Page{Roles: []string{"superuser"}},
				Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "invoices.InvoiceDetailMenu"}},
				Children: []components.PageInterface{
					&components.FormListenBoostedPost{
						Name: invUpdate,
						ActionURL: lamu.RoutePath("invoices.InvoiceUpdateRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("invoice.ID")),
						}),
						Children: []components.PageInterface{
							&components.FormComponent[Invoice]{
								Getter:        getters.Key[Invoice]("invoice"),
								Attr:          getters.FormBubbling(invUpdate),
								Title:         "Edit invoice",
								Subtitle:      "Update header and amounts",
								ChildrenInput: invoiceFormInputs(),
								ChildrenAction: []components.PageInterface{
									&components.ContainerRow{
										Classes: "flex flex-wrap justify-between gap-2 mt-2 items-center",
										Children: []components.PageInterface{
											&components.ContainerRow{
												Classes: "flex justify-end gap-2",
												Children: []components.PageInterface{
													&components.ButtonSubmit{Label: "Save"},
													&components.ButtonModalForm{
														Page:        components.Page{Roles: []string{"superuser"}},
														Label:       "Delete",
														Icon:        "trash",
														Name:        invDelete,
														Url: lamu.RoutePath("invoices.InvoiceDeleteRoute", map[string]getters.Getter[any]{
															"id": getters.Any(getters.Key[uint]("invoice.ID")),
														}),
														FormPostURL: lamu.RoutePath("invoices.InvoiceDeleteRoute", map[string]getters.Getter[any]{
															"id": getters.Any(getters.Key[uint]("invoice.ID")),
														}),
														ModalUID: "invoice-delete-modal",
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
			{Key: "invoices.InvoiceDeleteForm", Value: &components.Modal{
				Page: components.Page{Roles: []string{"superuser"}},
				UID:  "invoice-delete-modal",
				Children: []components.PageInterface{
					&components.DeleteConfirmation{
						Title:   "Delete invoice?",
						Message: "This removes the invoice and its lines.",
						Attr:    getters.FormBubbling(invDelete),
					},
				},
			}},
			{Key: "invoices.InvoiceDetail", Value: &components.ShellScaffold{
				Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "invoices.InvoiceDetailMenu"}},
				Children: []components.PageInterface{
					&components.Detail[Invoice]{
						Getter: getters.Key[Invoice]("invoice"),
						Children: []components.PageInterface{
							&components.ContainerColumn{
								Classes: "p-4 space-y-2",
								Children: []components.PageInterface{
									&components.LabelInline{
										Title: "Entity",
										Children: []components.PageInterface{
											&components.FieldText{Getter: getters.Key[string]("$in.Entity.Name")},
										},
									},
									&components.LabelInline{
										Title: "Partner",
										Children: []components.PageInterface{
											&components.FieldText{Getter: getters.Key[string]("$in.Partner.Name")},
										},
									},
									&components.LabelInline{
										Title: "Journal",
										Children: []components.PageInterface{
											&components.FieldText{Getter: getters.Key[string]("$in.Journal.Name")},
										},
									},
									&components.LabelInline{
										Title: "Type",
										Children: []components.PageInterface{
											&components.FieldText{Getter: getters.Key[string]("$in.InvoiceType")},
										},
									},
									&components.LabelInline{
										Title: "State",
										Children: []components.PageInterface{
											&components.FieldText{Getter: getters.Key[string]("$in.State")},
										},
									},
									&components.LabelInline{
										Title: "Reference",
										Children: []components.PageInterface{
											&components.FieldText{Getter: getters.Key[string]("$in.Reference")},
										},
									},
									&components.LabelInline{
										Title: "Number",
										Children: []components.PageInterface{
											&components.FieldText{Getter: getters.Deref(getters.Key[*string]("$in.Number"))},
										},
									},
									&components.LabelInline{
										Title: "Invoice date",
										Children: []components.PageInterface{
											&components.FieldText{Getter: getters.Format("%v", getters.Any(getters.Key[time.Time]("$in.InvoiceDate")))},
										},
									},
									&components.LabelInline{
										Title: "Due date",
										Children: []components.PageInterface{
											&components.FieldText{Getter: getters.Format("%v", getters.Any(getters.Deref(getters.Key[*time.Time]("$in.DueDate"))))},
										},
									},
									&components.LabelInline{
										Title: "Payment term",
										Children: []components.PageInterface{
											&components.FieldText{Getter: getters.Key[string]("$in.PaymentTerm.Name")},
										},
									},
									&components.LabelInline{
										Title: "Currency",
										Children: []components.PageInterface{
											&components.FieldText{Getter: getters.Format("%s — %s", getters.Any(getters.Key[string]("$in.Currency.Code")), getters.Any(getters.Key[string]("$in.Currency.Name")))},
										},
									},
									&components.LabelInline{
										Title: "Posted entry",
										Children: []components.PageInterface{
											&components.FieldText{Getter: getters.Format("#%d — %s", getters.Any(getters.Key[uint]("$in.Move.ID")), getters.Any(getters.Key[string]("$in.Move.Journal.Name")))},
										},
									},
									&components.LabelInline{
										Title: "Untaxed",
										Children: []components.PageInterface{
											&components.FieldText{Getter: getters.Format("%v", getters.Any(getters.Key[fields.DecimalSix]("$in.AmountUntaxed")))},
										},
									},
									&components.LabelInline{
										Title: "Tax",
										Children: []components.PageInterface{
											&components.FieldText{Getter: getters.Format("%v", getters.Any(getters.Key[fields.DecimalSix]("$in.AmountTax")))},
										},
									},
									&components.LabelInline{
										Title: "Total",
										Children: []components.PageInterface{
											&components.FieldText{Getter: getters.Format("%v", getters.Any(getters.Key[fields.DecimalSix]("$in.AmountTotal")))},
										},
									},
									&components.LabelInline{
										Title: "Residual",
										Children: []components.PageInterface{
											&components.FieldText{Getter: getters.Format("%v", getters.Any(getters.Key[fields.DecimalSix]("$in.AmountResidual")))},
										},
									},
								},
							},
						},
					},
					&components.DataTable[InvoiceLine]{
						UID:     "invoices-lines-table",
						Classes: "w-full",
						Title:   "Invoice lines",
						Data:    getters.Key[components.ObjectList[InvoiceLine]]("invoiceLines"),
						Actions: []components.PageInterface{
							&components.TableButtonCreate{
								Link: lamu.RoutePath("invoices.InvoiceLineCreateRoute", map[string]getters.Getter[any]{
									"invoiceId": getters.Any(getters.Key[uint]("invoice.ID")),
								}),
								Page: components.Page{Roles: []string{"superuser"}},
							},
						},
						RowAttr: getters.RowAttrNavigate(lamu.RoutePath("invoices.InvoiceLineDetailRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("$row.ID")),
						})),
						Columns: []components.TableColumn{
							{Label: "Label", Name: "Label", Children: []components.PageInterface{
								&components.FieldText{Getter: getters.Key[string]("$row.Label")},
							}},
							{Label: "Qty", Name: "Quantity", Children: []components.PageInterface{
								&components.FieldText{Getter: getters.Format("%v", getters.Any(getters.Key[fields.DecimalSix]("$row.Quantity")))},
							}},
							{Label: "Unit price", Name: "PriceUnit", Children: []components.PageInterface{
								&components.FieldText{Getter: getters.Format("%v", getters.Any(getters.Key[fields.DecimalSix]("$row.PriceUnit")))},
							}},
							{Label: "Subtotal", Name: "PriceSubtotal", Children: []components.PageInterface{
								&components.FieldText{Getter: getters.Format("%v", getters.Any(getters.Key[fields.DecimalSix]("$row.PriceSubtotal")))},
							}},
							{Label: "Account", Name: "Account", Children: []components.PageInterface{
								&components.FieldText{Getter: getters.Format("%s — %s", getters.Any(getters.Key[string]("$row.Account.Code")), getters.Any(getters.Key[string]("$row.Account.Name")))},
							}},
						},
					},
				},
			}},
			{Key: "invoices.InvoiceLineCreateForm", Value: &components.ShellScaffold{
				Page:    components.Page{Roles: []string{"superuser"}},
				Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "invoices.InvoiceDetailMenu"}},
				Children: []components.PageInterface{
					&components.FormListenBoostedPost{
						Name: lineCreate,
						ActionURL: lamu.RoutePath("invoices.InvoiceLineCreateRoute", map[string]getters.Getter[any]{
							"invoiceId": getters.Any(getters.Key[uint]("invoice.ID")),
						}),
						Children: []components.PageInterface{
							&components.FormComponent[InvoiceLine]{
								Attr:          getters.FormBubbling(lineCreate),
								Title:         "Add invoice line",
								Subtitle:      "Line item on this invoice",
								ChildrenInput: invoiceLineFormInputs(),
								ChildrenAction: []components.PageInterface{
									&components.ButtonSubmit{Label: "Save"},
								},
							},
						},
					},
				},
			}},
			{Key: "invoices.InvoiceLineUpdateForm", Value: &components.ShellScaffold{
				Page:    components.Page{Roles: []string{"superuser"}},
				Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "invoices.InvoiceLineDetailMenu"}},
				Children: []components.PageInterface{
					&components.FormListenBoostedPost{
						Name: lineUpdate,
						ActionURL: lamu.RoutePath("invoices.InvoiceLineUpdateRoute", map[string]getters.Getter[any]{
							"id": getters.Any(getters.Key[uint]("invoiceLine.ID")),
						}),
						Children: []components.PageInterface{
							&components.FormComponent[InvoiceLine]{
								Getter:        getters.Key[InvoiceLine]("invoiceLine"),
								Attr:          getters.FormBubbling(lineUpdate),
								Title:         "Edit invoice line",
								Subtitle:      "Update quantities, pricing, and taxes",
								ChildrenInput: invoiceLineFormInputs(),
								ChildrenAction: []components.PageInterface{
									&components.ContainerRow{
										Classes: "flex flex-wrap justify-between gap-2 mt-2 items-center",
										Children: []components.PageInterface{
											&components.ContainerRow{
												Classes: "flex justify-end gap-2",
												Children: []components.PageInterface{
													&components.ButtonSubmit{Label: "Save"},
													&components.ButtonModalForm{
														Page:        components.Page{Roles: []string{"superuser"}},
														Label:       "Delete",
														Icon:        "trash",
														Name:        lineDelete,
														Url: lamu.RoutePath("invoices.InvoiceLineDeleteRoute", map[string]getters.Getter[any]{
															"id": getters.Any(getters.Key[uint]("invoiceLine.ID")),
														}),
														FormPostURL: lamu.RoutePath("invoices.InvoiceLineDeleteRoute", map[string]getters.Getter[any]{
															"id": getters.Any(getters.Key[uint]("invoiceLine.ID")),
														}),
														ModalUID: "invoice-line-delete-modal",
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
			{Key: "invoices.InvoiceLineDeleteForm", Value: &components.Modal{
				Page: components.Page{Roles: []string{"superuser"}},
				UID:  "invoice-line-delete-modal",
				Children: []components.PageInterface{
					&components.DeleteConfirmation{
						Title:   "Delete invoice line?",
						Message: "This removes the line from the invoice.",
						Attr:    getters.FormBubbling(lineDelete),
					},
				},
			}},
			{Key: "invoices.InvoiceLineDetail", Value: &components.ShellScaffold{
				Sidebar: []components.PageInterface{lamu.DynamicPage{Name: "invoices.InvoiceLineDetailMenu"}},
				Children: []components.PageInterface{
					&components.Detail[InvoiceLine]{
						Getter: getters.Key[InvoiceLine]("invoiceLine"),
						Children: []components.PageInterface{
							&components.ContainerColumn{
								Classes: "p-4 space-y-2",
								Children: []components.PageInterface{
									&components.LabelInline{
										Title: "Invoice",
										Children: []components.PageInterface{
											&components.FieldText{Getter: getters.Format("#%d", getters.Any(getters.Key[uint]("$in.Invoice.ID")))},
										},
									},
									&components.LabelInline{
										Title: "Product",
										Children: []components.PageInterface{
											&components.FieldText{Getter: getters.Format("%s — %s", getters.Any(getters.Key[string]("$in.Product.Code")), getters.Any(getters.Key[string]("$in.Product.Name")))},
										},
									},
									&components.LabelInline{
										Title: "Label",
										Children: []components.PageInterface{
											&components.FieldText{Getter: getters.Key[string]("$in.Label")},
										},
									},
									&components.LabelInline{
										Title: "Quantity",
										Children: []components.PageInterface{
											&components.FieldText{Getter: getters.Format("%v", getters.Any(getters.Key[fields.DecimalSix]("$in.Quantity")))},
										},
									},
									&components.LabelInline{
										Title: "Unit price",
										Children: []components.PageInterface{
											&components.FieldText{Getter: getters.Format("%v", getters.Any(getters.Key[fields.DecimalSix]("$in.PriceUnit")))},
										},
									},
									&components.LabelInline{
										Title: "Discount %",
										Children: []components.PageInterface{
											&components.FieldText{Getter: getters.Format("%v", getters.Any(getters.Key[fields.DecimalSix]("$in.Discount")))},
										},
									},
									&components.LabelInline{
										Title: "Subtotal",
										Children: []components.PageInterface{
											&components.FieldText{Getter: getters.Format("%v", getters.Any(getters.Key[fields.DecimalSix]("$in.PriceSubtotal")))},
										},
									},
									&components.LabelInline{
										Title: "Account",
										Children: []components.PageInterface{
											&components.FieldText{Getter: getters.Format("%s — %s", getters.Any(getters.Key[string]("$in.Account.Code")), getters.Any(getters.Key[string]("$in.Account.Name")))},
										},
									},
									&components.FieldManyToMany[tax.TaxRate]{
										Label:  "Taxes",
										Getter: getters.Key[[]tax.TaxRate]("$in.Taxes"),
										Display: getters.Format("%s (%v)", getters.Any(getters.Key[string]("$in.Name")),
											getters.Any(getters.Key[fields.DecimalSix]("$in.Amount"))),
									},
								},
							},
						},
					},
				},
			}},
		},
	}
}
