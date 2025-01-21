// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package store

import (
	"context"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT 
INTO users (username, email, password) 
VALUES ($1, $2, $3)
RETURNING id
`

type CreateUserParams struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Username, arg.Email, arg.Password)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const getUserByUserId = `-- name: GetUserByUserId :one
SELECT id, email, username, password, created_at 
FROM users 
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUserByUserId(ctx context.Context, id uuid.UUID) (Users, error) {
	row := q.db.QueryRowContext(ctx, getUserByUserId, id)
	var i Users
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Password,
		&i.CreatedAt,
	)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id, email, username, password, created_at 
FROM users 
WHERE  username = $1 LIMIT 1
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (Users, error) {
	row := q.db.QueryRowContext(ctx, getUserByUsername, username)
	var i Users
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Password,
		&i.CreatedAt,
	)
	return i, err
}
