package p_uniquity_finance_accounts

import (
	"context"
	"fmt"

	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/components"
	"github.com/lariv-in/lago/getters"
	"github.com/lariv-in/lago/registry"
)

// filterGETString reads one filter field from context $get for display.
// Empty when absent so boosted GET does not submit spurious zeros (unlike InputNumber + Key[int]("$get.X")).
func filterGETString(field string) getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		m, ok := ctx.Value("$get").(map[string]any)
		if !ok {
			return "", nil
		}
		v, ok := m[field]
		if !ok || v == nil {
			return "", nil
		}
		switch t := v.(type) {
		case string:
			return t, nil
		default:
			return fmt.Sprintf("%v", t), nil
		}
	}
}

func pageFilterPages() []registry.Pair[string, components.PageInterface] {
	filterButtons := []components.PageInterface{
		&components.ContainerRow{
			Classes: "flex gap-2",
			Children: []components.PageInterface{
				&components.ButtonSubmit{Label: "Apply Filters"},
				&components.ButtonClear{Label: "Clear"},
			},
		},
	}

	return []registry.Pair[string, components.PageInterface]{
		{Key: "finance_accounts.AccountFilter", Value: &components.FormComponent[Account]{
			Attr: getters.FormBoostedGet(lago.RoutePath("finance_accounts.DefaultRoute", nil)),
			ChildrenInput: []components.PageInterface{
				&components.InputText{Name: "Name", Label: "Name", Getter: getters.Key[string]("$get.Name")},
				&components.InputText{Name: "Code", Label: "Code", Getter: filterGETString("Code")},
				&components.InputCheckbox{Name: "IsGroup", Label: "Group account", Getter: getters.Key[bool]("$get.IsGroup")},
				&components.InputSelect[BalanceType]{
					Name:     "BalanceType",
					Label:    "Balance type",
					Required: false,
					Choices:  balanceTypeChoices,
					Getter:   registry.PairFromGetter(getters.Key[BalanceType]("$get.BalanceType"), balanceTypeChoiceList),
				},
			},
			ChildrenAction: filterButtons,
		}},
		{Key: "finance_accounts.AccountSelectionFilter", Value: &components.FormComponent[Account]{
			Attr: getters.FormBoostedGet(lago.RoutePath("finance_accounts.AccountSelectRoute", nil)),
			ChildrenInput: []components.PageInterface{
				&components.InputText{Hidden: true, Name: "target_input", Getter: getters.Key[string]("$get.target_input")},
				&components.InputText{Hidden: true, Name: balanceTypeScopeQueryParam, Getter: filterGETString(balanceTypeScopeQueryParam)},
				&components.InputText{Name: "Name", Label: "Name", Getter: getters.Key[string]("$get.Name")},
				&components.InputText{Name: "Code", Label: "Code", Getter: filterGETString("Code")},
				&components.InputSelect[BalanceType]{
					Name:     "BalanceType",
					Label:    "Balance type",
					Required: false,
					Choices:  balanceTypeChoices,
					Getter:   registry.PairFromGetter(getters.Key[BalanceType]("$get.BalanceType"), balanceTypeChoiceList),
				},
				&components.InputText{Name: "ParentID", Label: "Parent ID", Getter: filterGETString("ParentID")},
			},
			ChildrenAction: filterButtons,
		}},
		{Key: "finance_accounts.CurrencyFilter", Value: &components.FormComponent[Currency]{
			Attr: getters.FormBoostedGet(lago.RoutePath("finance_accounts.CurrencyListRoute", nil)),
			ChildrenInput: []components.PageInterface{
				&components.InputText{Name: "Code", Label: "Numeric code", Getter: filterGETString("Code")},
				&components.InputText{Name: "Name", Label: "Name", Getter: getters.Key[string]("$get.Name")},
				&components.InputText{Name: "Symbol", Label: "Symbol", Getter: getters.Key[string]("$get.Symbol")},
				&components.InputText{Name: "MinorUnit", Label: "Minor unit", Getter: filterGETString("MinorUnit")},
			},
			ChildrenAction: filterButtons,
		}},
		{Key: "finance_accounts.CurrencySelectionFilter", Value: &components.FormComponent[Currency]{
			Attr: getters.FormBoostedGet(lago.RoutePath("finance_accounts.CurrencySelectRoute", nil)),
			ChildrenInput: []components.PageInterface{
				&components.InputText{Name: "Code", Label: "Numeric code", Getter: filterGETString("Code")},
				&components.InputText{Name: "Name", Label: "Name", Getter: getters.Key[string]("$get.Name")},
				&components.InputText{Name: "Symbol", Label: "Symbol", Getter: getters.Key[string]("$get.Symbol")},
			},
			ChildrenAction: filterButtons,
		}},
		{Key: "finance_accounts.JournalFilter", Value: &components.FormComponent[Journal]{
			Attr: getters.FormBoostedGet(lago.RoutePath("finance_accounts.JournalListRoute", nil)),
			ChildrenInput: []components.PageInterface{
				&components.InputText{Name: "Name", Label: "Name", Getter: getters.Key[string]("$get.Name")},
				&components.InputCheckbox{Name: "IsActive", Label: "Active", Getter: getters.Key[bool]("$get.IsActive")},
				&components.InputText{Name: "CurrencyID", Label: "Currency ID", Getter: filterGETString("CurrencyID")},
				&components.InputSelect[JournalType]{
					Name:     "Type",
					Label:    "Type",
					Required: false,
					Choices:  journalTypeChoices,
					Getter:   registry.PairFromGetter(getters.Key[JournalType]("$get.Type"), journalTypeChoiceList),
				},
			},
			ChildrenAction: filterButtons,
		}},
	}
}
