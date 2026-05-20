package p_uniquity_finance_accounts

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/UniquityVentures/lamu/components"
	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/views"
	"gorm.io/gorm"
	"maragu.dev/gomponents"
)

// accountParentUpRowID is the synthetic row id for the ".." parent directory entry.
const accountParentUpRowID uint = 0

func accountParentUpRow() Account {
	return Account{
		Model:   gorm.Model{ID: accountParentUpRowID},
		Name:    "..",
		IsGroup: true,
	}
}

func prependAccountParentUpRow(list components.ObjectList[Account]) components.ObjectList[Account] {
	if list.Number != 1 {
		return list
	}
	items := append([]Account{accountParentUpRow()}, list.Items...)
	return components.ObjectList[Account]{
		Items:    items,
		Number:   list.Number,
		NumPages: list.NumPages,
		Total:    list.Total + 1,
	}
}

func parentIDFromGet(ctx context.Context) uint {
	m, ok := ctx.Value("$get").(map[string]any)
	if !ok {
		return 0
	}
	v, ok := m["ParentID"]
	if !ok || v == nil {
		return 0
	}
	switch t := v.(type) {
	case uint:
		return t
	case int:
		if t <= 0 {
			return 0
		}
		return uint(t)
	case int64:
		if t <= 0 {
			return 0
		}
		return uint(t)
	case uint64:
		return uint(t)
	case string:
		s := strings.TrimSpace(t)
		if s == "" {
			return 0
		}
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil || n == 0 {
			return 0
		}
		return uint(n)
	default:
		s := strings.TrimSpace(fmt.Sprint(v))
		if s == "" {
			return 0
		}
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil || n == 0 {
			return 0
		}
		return uint(n)
	}
}

func accountSelectBuildParentUpURL(ctx context.Context) (string, error) {
	currentParentID := parentIDFromGet(ctx)
	if currentParentID == 0 {
		return "", fmt.Errorf("already at root")
	}
	db, err := getters.DBFromContext(ctx)
	if err != nil {
		return "", err
	}
	var folder Account
	if err := db.Select("parent_id").First(&folder, currentParentID).Error; err != nil {
		return "", err
	}
	base, err := lamu.RoutePath("finance_accounts.AccountSelectRoute", nil)(ctx)
	if err != nil {
		return "", err
	}
	u, err := url.Parse(base)
	if err != nil {
		return "", fmt.Errorf("parse account select URL: %w", err)
	}
	q := u.Query()
	if m, ok := ctx.Value("$get").(map[string]any); ok {
		for k, v := range m {
			if k == "page" || k == "ParentID" {
				continue
			}
			s := strings.TrimSpace(fmt.Sprint(v))
			if s == "" {
				continue
			}
			q.Set(k, s)
		}
	}
	if folder.ParentID == nil || *folder.ParentID == 0 {
		q.Del("ParentID")
	} else {
		q.Set("ParentID", strconv.FormatUint(uint64(*folder.ParentID), 10))
	}
	q.Set("page", "1")
	u.RawQuery = q.Encode()
	return u.String(), nil
}

// accountSelectParentUpLayer prepends a ".." row when the picker is scoped to a parent account.
type accountSelectParentUpLayer struct{}

func (accountSelectParentUpLayer) Next(_ views.View, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if parentIDFromGet(ctx) == 0 {
			next.ServeHTTP(w, r)
			return
		}
		list, ok := ctx.Value("accounts").(components.ObjectList[Account])
		if !ok {
			next.ServeHTTP(w, r)
			return
		}
		ctx = context.WithValue(ctx, "accounts", prependAccountParentUpRow(list))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func accountTableRowIsParentUp(rowPrefix string) getters.Getter[bool] {
	return func(ctx context.Context) (bool, error) {
		id, err := getters.Key[uint](rowPrefix + ".ID")(ctx)
		if err != nil {
			return false, err
		}
		return id == accountParentUpRowID, nil
	}
}

func accountTableCodeCell(rowPrefix string) []components.PageInterface {
	return []components.PageInterface{
		&components.ShowIf{
			Getter: func(ctx context.Context) (any, error) {
				return accountTableRowIsParentUp(rowPrefix)(ctx)
			},
			Children: []components.PageInterface{
				&components.FieldText{Getter: getters.Static("—")},
			},
		},
		&components.ShowIf{
			Getter: func(ctx context.Context) (any, error) {
				isUp, err := accountTableRowIsParentUp(rowPrefix)(ctx)
				if err != nil {
					return false, err
				}
				return !isUp, nil
			},
			Children: []components.PageInterface{
				&components.FieldText{Getter: getters.Format("%d", getters.Any(getters.Key[int](rowPrefix+".Code")))},
			},
		},
	}
}

func accountTableBalanceTypeCell(rowPrefix string) []components.PageInterface {
	return []components.PageInterface{
		&components.ShowIf{
			Getter: func(ctx context.Context) (any, error) {
				return accountTableRowIsParentUp(rowPrefix)(ctx)
			},
			Children: []components.PageInterface{
				&components.FieldText{Getter: getters.Static("—")},
			},
		},
		&components.ShowIf{
			Getter: func(ctx context.Context) (any, error) {
				isUp, err := accountTableRowIsParentUp(rowPrefix)(ctx)
				if err != nil {
					return false, err
				}
				return !isUp, nil
			},
			Children: []components.PageInterface{
				&components.FieldText{Getter: getters.Format("%s", getters.Any(getters.Key[BalanceType](rowPrefix+".BalanceType")))},
			},
		},
	}
}

func accountChildrenTableRowAttr() getters.Getter[gomponents.Node] {
	detailAttr := getters.RowAttrNavigate(lamu.RoutePath("finance_accounts.AccountDetailRoute", map[string]getters.Getter[any]{
		"id": getters.Any(getters.Key[uint]("$row.ID")),
	}))
	return func(ctx context.Context) (gomponents.Node, error) {
		id, err := getters.Key[uint]("$row.ID")(ctx)
		if err != nil {
			return nil, err
		}
		if id != accountParentUpRowID {
			return detailAttr(ctx)
		}
		acc, ok := ctx.Value("account").(Account)
		if !ok {
			return nil, fmt.Errorf("parent up row without account context")
		}
		if acc.ParentID != nil && *acc.ParentID != 0 {
			pid := *acc.ParentID
			return getters.RowAttrNavigate(lamu.RoutePath("finance_accounts.AccountDetailRoute", map[string]getters.Getter[any]{
				"id": func(context.Context) (any, error) { return pid, nil },
			}))(ctx)
		}
		return getters.RowAttrNavigate(lamu.RoutePath("finance_accounts.DefaultRoute", nil))(ctx)
	}
}
