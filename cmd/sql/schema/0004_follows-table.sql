-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS follows (
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  follow_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
  PRIMARY KEY (user_id, follow_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS follows;
-- +goose StatementEnd 