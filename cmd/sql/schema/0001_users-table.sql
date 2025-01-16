-- goose postgres postgres://admin:adminpassword@localhost:5432/go-social up

-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS citext;
CREATE TABLE IF NOT EXISTS users (
  id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  email citext UNIQUE NOT NULL,
  username VARCHAR(255) UNIQUE NOT NULL,
  password bytea NOT NULL,
  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP EXTENSION IF EXISTS citext;
-- +goose StatementEnd
