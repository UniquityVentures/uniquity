package p_uniquity_employees

import (
	"log"
	"net/url"

	"github.com/UniquityVentures/lago/lago"
)

const AppUrl = "/employees/"

func init() {
	u, err := url.Parse(AppUrl)
	if err != nil {
		log.Panic(err)
	}
	if err := lago.RegistryPlugin.Register("p_uniquity_employees", lago.Plugin{
		Type:        lago.PluginTypeApp,
		Icon:        "users",
		URL:         u,
		VerboseName: "Employees & points",
		Roles:       []string{"superuser"},
	}); err != nil {
		log.Panic(err)
	}
}
