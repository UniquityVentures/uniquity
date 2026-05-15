-- +goose Up
-- Enforce: child.balance_type = parent.balance_type when parent_id is set;
-- and parent balance_type cannot be updated while any (non-deleted) child has a different type.
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION accounts_enforce_parent_balance_type() RETURNS TRIGGER AS $fn$
BEGIN
  IF NEW.parent_id IS NOT NULL THEN
    IF NOT EXISTS (
      SELECT 1 FROM accounts AS p
      WHERE p.id = NEW.parent_id
        AND p.deleted_at IS NULL
        AND p.balance_type = NEW.balance_type
    ) THEN
      RAISE EXCEPTION 'balance_type must match the parent account balance_type';
    END IF;
  END IF;

  IF TG_OP = 'UPDATE' AND NEW.balance_type IS DISTINCT FROM OLD.balance_type THEN
    IF EXISTS (
      SELECT 1 FROM accounts AS c
      WHERE c.parent_id = NEW.id
        AND c.deleted_at IS NULL
        AND c.balance_type IS DISTINCT FROM NEW.balance_type
    ) THEN
      RAISE EXCEPTION 'cannot change balance_type while child accounts have a different balance_type';
    END IF;
  END IF;

  RETURN NEW;
END;
$fn$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS accounts_enforce_parent_balance_type_biud ON accounts;

CREATE TRIGGER accounts_enforce_parent_balance_type_biud
  BEFORE INSERT OR UPDATE ON accounts
  FOR EACH ROW EXECUTE FUNCTION accounts_enforce_parent_balance_type();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS accounts_enforce_parent_balance_type_biud ON accounts;
DROP FUNCTION IF EXISTS accounts_enforce_parent_balance_type();
-- +goose StatementEnd
