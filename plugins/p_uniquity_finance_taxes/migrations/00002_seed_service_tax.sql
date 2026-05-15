-- +goose Up
INSERT INTO taxes (created_at, updated_at, name, percentage)
VALUES (now(), now(), 'Service Tax', 18);

-- +goose Down
DELETE FROM taxes WHERE name = 'Service Tax' AND percentage = 18;
