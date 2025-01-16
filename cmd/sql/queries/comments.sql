-- name: CreateComment :one
INSERT 
INTO comments (post_id, user_id, content) 
VALUES ($1, $2, $3)
RETURNING id, created_at;

-- name: GetCommentsByPostId :many
SELECT c.id, c.content, c.created_at, c.user_id, u.username 
FROM comments c
JOIN users u ON u.id = c.user_id
WHERE c.post_id = $1
ORDER BY c.created_at DESC;