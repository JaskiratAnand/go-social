-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_invitations (
  token UUID NOT NULL PRIMARY KEY,
  user_id UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  expiary TIMESTAMP(0) WITH TIME ZONE NOT NULL
);

ALTER TABLE users ADD COLUMN IF NOT EXISTS verified BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN IF EXISTS verified;
DROP TABLE IF EXISTS user_invitations;
-- +goose StatementEnd