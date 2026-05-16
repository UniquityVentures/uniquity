package main

import (
	"log/slog"

	"github.com/UniquityVentures/lamu/lamu"
	"github.com/UniquityVentures/lamu/registry"

	"github.com/UniquityVentures/lamu/plugins/p_dashboard"
	"github.com/UniquityVentures/lamu/plugins/p_filesystem"
	"github.com/UniquityVentures/lamu/plugins/p_livereloading"
	"github.com/UniquityVentures/lamu/plugins/p_otp"
	"github.com/UniquityVentures/lamu/plugins/p_pwa"
	"github.com/UniquityVentures/lamu/plugins/p_users"
	"github.com/UniquityVentures/uniquity/plugins/p_uniquity_employees"
	"github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_accounts"
	"github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_creditnotes"
	"github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_customer"
	"github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_fiscal_year"
	"github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_invoices"
	"github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_products"
	"github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_taxes"
	"github.com/UniquityVentures/uniquity/plugins/p_uniquity_video"
)

func main() {
	plugins := []registry.Pair[string, lamu.Plugin]{
		p_dashboard.GetPlugin(),
		p_filesystem.GetPlugin(),
		p_users.GetPlugin(),
		p_uniquity_employees.GetPlugin(),
		p_uniquity_finance_accounts.GetPlugin(),
		p_uniquity_finance_customer.GetPlugin(),
		p_uniquity_finance_creditnotes.GetPlugin(),
		p_uniquity_finance_fiscal_year.GetPlugin(),
		p_uniquity_finance_taxes.GetPlugin(),
		p_uniquity_finance_products.GetPlugin(),
		p_uniquity_finance_invoices.GetPlugin(),
		p_uniquity_video.GetPlugin(),
		p_livereloading.GetPlugin(),
		p_otp.GetPlugin(),
		p_pwa.GetPlugin(),
	}
	config, err := lamu.LoadConfigFromFile("uniquity.toml", plugins)
	if err != nil {
		panic(err)
	}
	if err := lamu.Start(config, plugins); err != nil {
		slog.Error(err.Error())
	}
}
