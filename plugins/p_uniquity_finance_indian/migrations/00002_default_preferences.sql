-- +goose Up

-- 1. Product preferences (inventory_account_id = 10301 Merchandise, cost_of_sales_account_id = 50201 Cost Of Sales)
INSERT INTO product_preferences (id, created_at, updated_at, inventory_account_id, cost_of_sales_account_id)
SELECT 1, now(), now(), inv.id, cos.id
FROM (SELECT id FROM accounts WHERE code = 10301 LIMIT 1) AS inv
CROSS JOIN (SELECT id FROM accounts WHERE code = 50201 LIMIT 1) AS cos
ON CONFLICT (id) DO UPDATE
SET inventory_account_id = EXCLUDED.inventory_account_id,
    cost_of_sales_account_id = EXCLUDED.cost_of_sales_account_id,
    updated_at = now();

-- 2. Invoice preferences (account_receivable_id = 10201 Accounts Receivable, account_revenue_id = 40101 Goods, account_tax_payable_id = 20203 Accrued Taxes)
INSERT INTO invoice_preferences (id, created_at, updated_at, account_receivable_id, account_revenue_id, account_tax_payable_id)
SELECT 1, now(), now(), ar.id, rev.id, tax.id
FROM (SELECT id FROM accounts WHERE code = 10201 LIMIT 1) AS ar
CROSS JOIN (SELECT id FROM accounts WHERE code = 40101 LIMIT 1) AS rev
CROSS JOIN (SELECT id FROM accounts WHERE code = 20203 LIMIT 1) AS tax
ON CONFLICT (id) DO UPDATE
SET account_receivable_id = EXCLUDED.account_receivable_id,
    account_revenue_id = EXCLUDED.account_revenue_id,
    account_tax_payable_id = EXCLUDED.account_tax_payable_id,
    updated_at = now();

-- 3. Payment preferences (payment_account_id = 10101 Cash and Cash Equivalents)
INSERT INTO payment_preferences (id, created_at, updated_at, payment_account_id)
SELECT 1, now(), now(), pay.id
FROM (SELECT id FROM accounts WHERE code = 10101 LIMIT 1) AS pay
ON CONFLICT (id) DO UPDATE
SET payment_account_id = EXCLUDED.payment_account_id,
    updated_at = now();

-- +goose Down
UPDATE product_preferences SET inventory_account_id = NULL, cost_of_sales_account_id = NULL WHERE id = 1;
UPDATE invoice_preferences SET account_receivable_id = NULL, account_revenue_id = NULL, account_tax_payable_id = NULL WHERE id = 1;
UPDATE payment_preferences SET payment_account_id = NULL WHERE id = 1;
