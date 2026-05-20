package p_uniquity_finance_accounts

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/UniquityVentures/lamu/getters"
	"github.com/UniquityVentures/lamu/lamu"
	"maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

const accountSelectionModalElementID = "finance-account-selection-modal"

// accountSelectionTableRowAttr handles account picker rows: group accounts drill down via HTMX;
// leaf accounts (posting / non-group) dispatch fk-select like [getters.RowAttrSelectNamed].
func accountSelectionTableRowAttr(
	nameGetter getters.Getter[string],
	idGetter getters.Getter[uint],
	displayGetter getters.Getter[string],
	isGroupGetter getters.Getter[bool],
) getters.Getter[gomponents.Node] {
	return func(ctx context.Context) (gomponents.Node, error) {
		rowID, err := idGetter(ctx)
		if err != nil {
			return nil, err
		}
		if rowID == accountParentUpRowID {
			drillURL, err := accountSelectBuildParentUpURL(ctx)
			if err != nil {
				return nil, err
			}
			return gomponents.Group{
				ghtml.Class("cursor-pointer hover:bg-base-200 transition-colors"),
				gomponents.Attr("hx-get", drillURL),
				gomponents.Attr("hx-target", "#"+accountSelectionModalElementID),
				gomponents.Attr("hx-swap", "outerHTML"),
				gomponents.Attr("hx-push-url", "false"),
			}, nil
		}
		isGroup, err := isGroupGetter(ctx)
		if err != nil {
			return nil, err
		}
		if isGroup {
			parentID, err := idGetter(ctx)
			if err != nil {
				return nil, err
			}
			drillURL, err := accountSelectBuildDrillURL(ctx, parentID)
			if err != nil {
				return nil, err
			}
			return gomponents.Group{
				ghtml.Class("cursor-pointer hover:bg-base-200 transition-colors"),
				gomponents.Attr("hx-get", drillURL),
				gomponents.Attr("hx-target", "#"+accountSelectionModalElementID),
				gomponents.Attr("hx-swap", "outerHTML"),
				gomponents.Attr("hx-push-url", "false"),
			}, nil
		}
		click := getters.SelectNamed(nameGetter, idGetter, displayGetter)
		return getters.RowAttrClickWithClass(click, nil)(ctx)
	}
}

func accountSelectBuildDrillURL(ctx context.Context, parentID uint) (string, error) {
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
			if k == "page" {
				continue
			}
			s := strings.TrimSpace(fmt.Sprint(v))
			if s == "" {
				continue
			}
			q.Set(k, s)
		}
	}
	q.Set("ParentID", strconv.FormatUint(uint64(parentID), 10))
	q.Set("page", "1")
	u.RawQuery = q.Encode()
	return u.String(), nil
}
