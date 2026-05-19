-- +goose Up
-- Indian GST output-tax rates aligned with CBIC goods rate schedules (CGST + SGST/UTGST = IGST for each slab).
-- Source: https://cbic-gst.gov.in/gst-goods-services-rates.html — Schedules I–VII (as notified, including Sept 2025 structure).

-- Posting accounts under Liabilities for GST collected on supplies (output tax payable).
WITH liab AS (
    SELECT id FROM accounts WHERE code = 20000 LIMIT 1
),
ins_group AS (
    INSERT INTO accounts (created_at, updated_at, name, code, is_group, balance_type, parent_id)
    SELECT now(), now(), 'India GST output payable', 20500, TRUE, 'Credit'::"BalanceType", liab.id
    FROM liab
    RETURNING id
)
INSERT INTO accounts (created_at, updated_at, name, code, is_group, balance_type, parent_id)
SELECT now(), now(), v.name, v.code, FALSE, 'Credit'::"BalanceType", g.id
FROM ins_group g
CROSS JOIN (VALUES
    ('CGST output payable', 20501),
    ('SGST and UTGST output payable', 20502),
    ('IGST output payable', 20503)
) AS v(name, code);

INSERT INTO taxes (created_at, updated_at, name, percentage, tax_type, account_id)
SELECT now(), now(), t.name, t.pct::NUMERIC(19, 6), 'levied'::"TaxKind", acct.id
FROM (VALUES
    ('CGST 2.5%', 2.5::NUMERIC, 20501),
    ('SGST 2.5%', 2.5::NUMERIC, 20502),
    ('IGST 5%', 5::NUMERIC, 20503),
    ('CGST 6%', 6::NUMERIC, 20501),
    ('SGST 6%', 6::NUMERIC, 20502),
    ('IGST 12%', 12::NUMERIC, 20503),
    ('CGST 9%', 9::NUMERIC, 20501),
    ('SGST 9%', 9::NUMERIC, 20502),
    ('IGST 18%', 18::NUMERIC, 20503),
    ('CGST 14%', 14::NUMERIC, 20501),
    ('SGST 14%', 14::NUMERIC, 20502),
    ('IGST 28%', 28::NUMERIC, 20503),
    ('CGST 1.5%', 1.5::NUMERIC, 20501),
    ('SGST 1.5%', 1.5::NUMERIC, 20502),
    ('IGST 3%', 3::NUMERIC, 20503),
    ('CGST 0.125%', 0.125::NUMERIC, 20501),
    ('SGST 0.125%', 0.125::NUMERIC, 20502),
    ('IGST 0.25%', 0.25::NUMERIC, 20503),
    ('CGST 0.75%', 0.75::NUMERIC, 20501),
    ('SGST 0.75%', 0.75::NUMERIC, 20502),
    ('IGST 1.5%', 1.5::NUMERIC, 20503)
) AS t(name, pct, acct_code)
JOIN accounts AS acct ON acct.code = t.acct_code;

-- +goose Down
DELETE FROM taxes
WHERE name IN (
    'CGST 2.5%', 'SGST 2.5%', 'IGST 5%',
    'CGST 6%', 'SGST 6%', 'IGST 12%',
    'CGST 9%', 'SGST 9%', 'IGST 18%',
    'CGST 14%', 'SGST 14%', 'IGST 28%',
    'CGST 1.5%', 'SGST 1.5%', 'IGST 3%',
    'CGST 0.125%', 'SGST 0.125%', 'IGST 0.25%',
    'CGST 0.75%', 'SGST 0.75%', 'IGST 1.5%'
);

DELETE FROM accounts WHERE code IN (20501, 20502, 20503, 20500);
