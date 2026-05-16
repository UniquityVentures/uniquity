-- +goose Up
-- Invoice numbers must be unique among posted invoices and unique among cancelled invoices.
-- (A posted row and its cancellation snapshot intentionally share the same number, so this is
-- per-table uniqueness, not one constraint across both tables.)
CREATE UNIQUE INDEX IF NOT EXISTS uix_posted_invoices_number_live
    ON posted_invoices (number)
    WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS uix_cancelled_invoices_number_live
    ON cancelled_invoices (number)
    WHERE deleted_at IS NULL;

-- +goose Down
DROP INDEX IF EXISTS uix_cancelled_invoices_number_live;
DROP INDEX IF EXISTS uix_posted_invoices_number_live;
