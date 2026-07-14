package p_uniquity_finance_taxes

import (
	"github.com/lariv-in/lago/registry"
	"github.com/lariv-in/lago/views"
)

// taxPreloadAccount loads the GL account for tax forms and list/detail views.
var taxPreloadAccount = views.QueryPatcherPreload[Tax]{Fields: []string{"Account"}}

func taxQueryPatchers() views.QueryPatchers[Tax] {
	return views.QueryPatchers[Tax]{
		registry.Pair[string, views.QueryPatcher[Tax]]{Key: "finance_taxes.preload_account", Value: taxPreloadAccount},
	}
}
