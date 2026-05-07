package main

import (
	"log/slog"

	"github.com/UniquityVentures/lago/lago"

	_ "github.com/UniquityVentures/lago/plugins/p_dashboard"
	_ "github.com/UniquityVentures/lago/plugins/p_filesystem"
	_ "github.com/UniquityVentures/lago/plugins/p_livereloading"
	_ "github.com/UniquityVentures/lago/plugins/p_otp"
	_ "github.com/UniquityVentures/lago/plugins/p_pwa"
	_ "github.com/UniquityVentures/lago/plugins/p_users"
	_ "github.com/UniquityVentures/uniquity_ventures/plugins/p_uniquity_employees"
	_ "github.com/UniquityVentures/uniquity_ventures/plugins/p_uniquity_video"
)

func main() {
	config, err := lago.LoadConfigFromFile("uniquity_ventures.toml")
	if err != nil {
		panic(err)
	}
	if err := lago.Start(config); err != nil {
		slog.Error(err.Error())
	}
}
