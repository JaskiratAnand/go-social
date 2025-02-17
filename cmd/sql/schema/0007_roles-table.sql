-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS roles (
    id int UNIQUE NOT NULL DEFAULT 0 PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT NOT NULL
);

INSERT INTO 
roles (id, name, description) 
VALUES (
    1,
    'USER', 
    'A user can create posts and comments'
);
INSERT INTO 
roles (id, name, description) 
VALUES (
    2,
    'MODERATOR', 
    'A moderator can update other users posts'
);
INSERT INTO 
roles (id, name, description) 
VALUES (
    3,
    'ADMIN', 
    'An admin can update and delete other users posts'
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS roles;
-- +goose StatementEnd