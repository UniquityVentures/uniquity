-- +goose Up
CREATE TYPE "ProductType" AS ENUM ('Goods', 'Services', 'Both');

ALTER TABLE products
    ADD COLUMN product_type "ProductType" NOT NULL DEFAULT 'Goods',
    ADD COLUMN reference TEXT,
    ADD COLUMN remarks TEXT;

UPDATE products
SET reference = 'LEGACY-' || id::text
WHERE reference IS NULL;

ALTER TABLE products
    ALTER COLUMN reference SET NOT NULL;

ALTER TABLE products ADD CONSTRAINT uq_products_reference UNIQUE (reference);

-- +goose Down
ALTER TABLE products
    DROP COLUMN IF EXISTS remarks,
    DROP COLUMN IF EXISTS reference,
    DROP COLUMN IF EXISTS product_type;

DROP TYPE IF EXISTS "ProductType";
