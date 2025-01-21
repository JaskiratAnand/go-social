// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: posts.sql

package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const createPost = `-- name: CreatePost :one
INSERT 
INTO posts (title, content, user_id, tags) 
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, updated_at
`

type CreatePostParams struct {
	Title   string    `json:"title"`
	Content string    `json:"content"`
	UserID  uuid.UUID `json:"user_id"`
	Tags    []string  `json:"tags"`
}

type CreatePostRow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (CreatePostRow, error) {
	row := q.db.QueryRowContext(ctx, createPost,
		arg.Title,
		arg.Content,
		arg.UserID,
		pq.Array(arg.Tags),
	)
	var i CreatePostRow
	err := row.Scan(&i.ID, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const deletePostById = `-- name: DeletePostById :exec
DELETE FROM posts WHERE id = $1
`

func (q *Queries) DeletePostById(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deletePostById, id)
	return err
}

const getPostWithCommentsById = `-- name: GetPostWithCommentsById :one
SELECT p.title, p.content, p.user_id, author.username, p.tags, p.created_at, p.updated_at, 
    JSON_AGG(
        JSON_BUILD_OBJECT(
            'id', c.id,
            'user_id', c.user_id,
            'username', u.username,
            'content', c.content,
            'created_at', c.created_at
        )
    ) AS comments
FROM posts p
LEFT JOIN users author ON p.user_id = author.id
LEFT JOIN comments c ON p.id = c.post_id
LEFT JOIN users u ON c.user_id = u.id
WHERE p.id = $1
GROUP BY p.id, author.username
`

type GetPostWithCommentsByIdRow struct {
	Title     string          `json:"title"`
	Content   string          `json:"content"`
	UserID    uuid.UUID       `json:"user_id"`
	Username  sql.NullString  `json:"username"`
	Tags      []string        `json:"tags"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Comments  json.RawMessage `json:"comments"`
}

func (q *Queries) GetPostWithCommentsById(ctx context.Context, id uuid.UUID) (GetPostWithCommentsByIdRow, error) {
	row := q.db.QueryRowContext(ctx, getPostWithCommentsById, id)
	var i GetPostWithCommentsByIdRow
	err := row.Scan(
		&i.Title,
		&i.Content,
		&i.UserID,
		&i.Username,
		pq.Array(&i.Tags),
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Comments,
	)
	return i, err
}

const getPostsById = `-- name: GetPostsById :one
SELECT title, content, tags, user_id, created_at, updated_at
FROM posts 
WHERE id = $1 LIMIT 1
`

type GetPostsByIdRow struct {
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) GetPostsById(ctx context.Context, id uuid.UUID) (GetPostsByIdRow, error) {
	row := q.db.QueryRowContext(ctx, getPostsById, id)
	var i GetPostsByIdRow
	err := row.Scan(
		&i.Title,
		&i.Content,
		pq.Array(&i.Tags),
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPostsByUserId = `-- name: GetPostsByUserId :many
SELECT id, title, content, tags, created_at, updated_at 
FROM posts 
WHERE user_id = $1 
ORDER BY created_at DESC
`

type GetPostsByUserIdRow struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) GetPostsByUserId(ctx context.Context, userID uuid.UUID) ([]GetPostsByUserIdRow, error) {
	rows, err := q.db.QueryContext(ctx, getPostsByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPostsByUserIdRow
	for rows.Next() {
		var i GetPostsByUserIdRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Content,
			pq.Array(&i.Tags),
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserFeed = `-- name: GetUserFeed :many
SELECT p.id, p.title, p.content, p.tags, p.created_at, p.updated_at, 
    u.username,
    COALESCE(COUNT(c.id), 0) AS comments_count
FROM posts p
LEFT JOIN users u ON p.user_id = u.id
LEFT JOIN comments c ON p.id = c.post_id
JOIN follows f ON p.user_id = f.follow_id OR p.user_id = $1
WHERE 
    f.user_id = $1 AND
    (p.title ILIKE '%' || $2::TEXT || '%' OR p.content ILIKE '%' || $2::TEXT || '%') AND 
    (p.tags @> $3 OR $3 = '{}')
GROUP BY p.id, u.username
ORDER BY p.created_at DESC
LIMIT $4 OFFSET $5
`

type GetUserFeedParams struct {
	UserID  uuid.UUID `json:"user_id"`
	Column2 string    `json:"column_2"`
	Tags    []string  `json:"tags"`
	Limit   int64     `json:"limit"`
	Offset  int64     `json:"offset"`
}

type GetUserFeedRow struct {
	ID            uuid.UUID      `json:"id"`
	Title         string         `json:"title"`
	Content       string         `json:"content"`
	Tags          []string       `json:"tags"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	Username      sql.NullString `json:"username"`
	CommentsCount sql.NullInt64  `json:"comments_count"`
}

func (q *Queries) GetUserFeed(ctx context.Context, arg GetUserFeedParams) ([]GetUserFeedRow, error) {
	rows, err := q.db.QueryContext(ctx, getUserFeed,
		arg.UserID,
		arg.Column2,
		pq.Array(arg.Tags),
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserFeedRow
	for rows.Next() {
		var i GetUserFeedRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Content,
			pq.Array(&i.Tags),
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Username,
			&i.CommentsCount,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatePostById = `-- name: UpdatePostById :one
UPDATE posts
SET 
    title = COALESCE($1, title),
    content = COALESCE($2, content),
    tags = COALESCE($3, tags)
WHERE id = $4 AND updated_at = $5
RETURNING id, updated_at
`

type UpdatePostByIdParams struct {
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	ID        uuid.UUID `json:"id"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdatePostByIdRow struct {
	ID        uuid.UUID `json:"id"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) UpdatePostById(ctx context.Context, arg UpdatePostByIdParams) (UpdatePostByIdRow, error) {
	row := q.db.QueryRowContext(ctx, updatePostById,
		arg.Title,
		arg.Content,
		pq.Array(arg.Tags),
		arg.ID,
		arg.UpdatedAt,
	)
	var i UpdatePostByIdRow
	err := row.Scan(&i.ID, &i.UpdatedAt)
	return i, err
}
