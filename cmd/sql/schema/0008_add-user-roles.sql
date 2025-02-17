-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS users 
ADD COLUMN IF NOT EXISTS 
role_id int NOT NULL REFERENCES roles(id) DEFAULT 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN IF EXISTS role_id;
-- +goose StatementEnd