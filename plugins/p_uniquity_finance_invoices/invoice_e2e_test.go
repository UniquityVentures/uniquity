package p_uniquity_finance_invoices

import (
	"context"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"testing"
	"time"

	finance_employees "github.com/UniquityVentures/uniquity/plugins/p_uniquity_employees"
	finance_accounts "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_accounts"
	finance_creditnotes "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_creditnotes"
	finance_customer "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_customer"
	finance_fiscal_year "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_fiscal_year"
	finance_indian "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_indian"
	finance_products "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_products"
	finance_taxes "github.com/UniquityVentures/uniquity/plugins/p_uniquity_finance_taxes"
	finance_video "github.com/UniquityVentures/uniquity/plugins/p_uniquity_video"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/lariv-in/lago"
	"github.com/lariv-in/lago/fields"
	"github.com/lariv-in/lago/plugins/p_dashboard"
	"github.com/lariv-in/lago/plugins/p_filesystem"
	"github.com/lariv-in/lago/plugins/p_users"
	"github.com/lariv-in/lago/registry"
	gorm_postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestCreateInvoiceE2E(t *testing.T) {
	// Connect to master postgres database to create a new temporary test database
	dbMaster, err := gorm.Open(gorm_postgres.Open("host=localhost user=postgres dbname=postgres sslmode=disable"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to master db: %v", err)
	}

	dbName := fmt.Sprintf("uniquity_test_%d", time.Now().UnixNano())
	if err := dbMaster.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)).Error; err != nil {
		t.Fatalf("failed to create database %s: %v", dbName, err)
	}

	t.Cleanup(func() {
		if err := dbMaster.Exec(fmt.Sprintf("DROP DATABASE %s WITH (FORCE)", dbName)).Error; err != nil {
			t.Logf("failed to drop database %s: %v", dbName, err)
		}
	})

	dsn := fmt.Sprintf("host=localhost user=postgres dbname=%s sslmode=disable", dbName)
	db, err := gorm.Open(gorm_postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	pgConfig := &gorm_postgres.Config{
		DSN: dsn,
	}

	config := lago.LagoConfig{
		Debug:          true,
		DBType:         lago.DBTypePostgres,
		PostgresConfig: pgConfig,
	}

	plugins := []registry.Pair[string, lago.Plugin]{
		p_dashboard.GetPlugin(),
		p_filesystem.GetPlugin(),
		p_users.GetPlugin(),
		finance_employees.GetPlugin(),
		finance_accounts.GetPlugin(),
		finance_customer.GetPlugin(),
		finance_creditnotes.GetPlugin(),
		finance_fiscal_year.GetPlugin(),
		finance_taxes.GetPlugin(),
		finance_products.GetPlugin(),
		GetPlugin(), // this plugin (p_uniquity_finance_invoices)
		finance_indian.GetPlugin(),
		finance_video.GetPlugin(),
	}

	lago.BuildAllRegistries(append([]registry.Pair[string, lago.Plugin]{lago.CorePlugin(db, config)}, plugins...))

	if err := lago.InitDB(db, config); err != nil {
		t.Fatalf("failed to initialize db: %v", err)
	}

	// Create necessary seed records in DB
	// 1. A superuser to bypass authentication
	user := p_users.User{
		Email:        "admin@example.com",
		IsSuperuser:  true,
		Timezone:     "UTC",
		PasswordHash: []byte("hash"),
		PasswordSalt: []byte("salt"),
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to create superuser: %v", err)
	}

	// 2. A payment term
	pt := PaymentTerm{
		Type:      "Immediate",
		BackingID: 1,
	}
	if err := db.Create(&pt).Error; err != nil {
		t.Fatalf("failed to create payment term: %v", err)
	}

	// 3. A customer
	cust := finance_customer.Customer{
		Name: "Test Customer",
	}
	if err := db.Create(&cust).Error; err != nil {
		t.Fatalf("failed to create customer: %v", err)
	}

	// 4. A tax
	tax := finance_taxes.Tax{
		Name:       "GST 18%",
		Percentage: fields.DecimalSix{R: big.NewRat(18, 1)},
		TaxType:    finance_taxes.TaxKindLevied,
	}
	if err := db.Create(&tax).Error; err != nil {
		t.Fatalf("failed to create tax: %v", err)
	}

	// 5. A product
	prod := finance_products.Product{
		Name:       "Test Product",
		SalesPrice: fields.DecimalSix{R: big.NewRat(100, 1)},
		Type:       finance_products.ProductTypeGoods,
		Reference:  "REF-001",
	}
	if err := db.Create(&prod).Error; err != nil {
		t.Fatalf("failed to create product: %v", err)
	}
	// Associate tax with product
	if err := db.Model(&prod).Association("Taxes").Append(&tax); err != nil {
		t.Fatalf("failed to associate tax with product: %v", err)
	}

	// Start HTTP server on a random free port
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer ln.Close()

	addr := ln.Addr().String()

	// Get router and wrap layers
	layers := *lago.RegistryLayer.AllStable()
	var router http.Handler = lago.GetRouter(config)
	for _, layer := range layers {
		router = layer.Value.Next(router)
	}
	router = http.NewCrossOriginProtection().Handler(router)

	srv := &http.Server{Handler: router}
	go func() {
		_ = srv.Serve(ln)
	}()
	defer srv.Shutdown(context.Background())

	// Launch Chrome via go-rod
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	// Generate JWT auth token for our superuser
	token, err := user.GetJwt(time.Now(), time.Now().Add(time.Hour*24))
	if err != nil {
		t.Fatalf("failed to generate JWT: %v", err)
	}

	// Set auth-token cookie so the browser is authenticated as the superuser
	cookie := &proto.NetworkCookie{
		Name:   "auth-token",
		Value:  token,
		Domain: "127.0.0.1",
		Path:   "/",
	}
	browser.MustSetCookies(cookie)

	// Navigate to draft invoice create page
	createURL := fmt.Sprintf("http://%s/finance-invoices/create/", addr)
	page := browser.MustPage(createURL)
	page.MustWaitLoad()

	// 1. Select Customer
	t.Log("Opening customer modal...")
	page.MustElementR("div.input", "Select customer…").MustClick()
	t.Log("Waiting for customer modal...")
	page.MustElement("#finance-customer-fk-select-modal").MustWaitVisible()
	t.Log("Clicking customer...")
	page.MustElementR("#finance-customer-fk-select-modal td, #finance-customer-fk-select-modal div", "Test Customer").MustClick()
	for page.MustHas("#finance-customer-fk-select-modal") {
		time.Sleep(50 * time.Millisecond)
	}

	// 2. Select Payment Term
	t.Log("Opening payment term modal...")
	page.MustElementR("div.input", "Select payment term…").MustClick()
	t.Log("Waiting for payment term modal...")
	page.MustElement("#finance-invoice-payment-term-fk-modal").MustWaitVisible()
	t.Log("Clicking payment term...")
	page.MustElementR("#finance-invoice-payment-term-fk-modal td, #finance-invoice-payment-term-fk-modal div", "#1").MustClick()
	for page.MustHas("#finance-invoice-payment-term-fk-modal") {
		time.Sleep(50 * time.Millisecond)
	}

	// 3. Select Product in line 1
	t.Log("Opening product modal...")
	page.MustElementR("div.input", "Select…").MustClick()
	t.Log("Waiting for product modal...")
	page.MustElement("#finance-product-fk-select-modal").MustWaitVisible()
	t.Log("Clicking product...")
	page.MustElementR("#finance-product-fk-select-modal td, #finance-product-fk-select-modal div", "Test Product").MustClick()
	for page.MustHas("#finance-product-fk-select-modal") {
		time.Sleep(50 * time.Millisecond)
	}

	// 4. Update quantity
	inputs := page.MustElements("table tbody tr td input[inputmode='decimal']")
	if len(inputs) < 2 {
		t.Fatalf("expected at least 2 input fields in lines table, got %d", len(inputs))
	}
	inputs[0].MustSelectAllText().MustInput("2")    // Quantity
	inputs[1].MustSelectAllText().MustInput("50.0") // Rate

	// Click Save (Submit form)
	page.MustElement("button[type='submit']").MustClick()
	page.MustWaitNavigation()

	// Verify that the draft invoice was created successfully in the DB
	var count int64
	if err := db.Model(&DraftInvoice{}).Count(&count).Error; err != nil {
		t.Fatalf("failed to query draft invoices count: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 draft invoice in database, got %d", count)
	}

	// Verify the invoice values are correct
	var inv DraftInvoice
	if err := db.Preload("Lines.Taxes").First(&inv).Error; err != nil {
		t.Fatalf("failed to load created draft invoice: %v", err)
	}
	if inv.CustomerID != cust.ID {
		t.Errorf("expected CustomerID %d, got %d", cust.ID, inv.CustomerID)
	}
	if len(inv.Lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(inv.Lines))
	}
	line := inv.Lines[0]
	if line.ProductID != prod.ID {
		t.Errorf("expected Line ProductID %d, got %d", prod.ID, line.ProductID)
	}
	if line.Quantity.String() != "2.000000" {
		t.Errorf("expected line quantity 2, got %s", line.Quantity.String())
	}
	if line.Rate.String() != "50.000000" {
		t.Errorf("expected line rate 50, got %s", line.Rate.String())
	}
	// Verify that the line tax was auto-populated and saved!
	if len(line.Taxes) != 1 || line.Taxes[0].ID != tax.ID {
		t.Errorf("expected line to have tax ID %d, got %v", tax.ID, line.Taxes)
	}
}
