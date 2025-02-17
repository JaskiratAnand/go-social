-- name: CreatePost :one
INSERT 
INTO posts (title, content, user_id, tags) 
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, updated_at;

-- name: GetPostsByUserId :many
SELECT id, title, content, tags, created_at, updated_at 
FROM posts 
WHERE user_id = $1 
ORDER BY created_at DESC;

-- name: GetPostsById :one
SELECT id, title, content, tags, user_id, created_at, updated_at
FROM posts 
WHERE id = $1 LIMIT 1;

-- name: GetPostWithCommentsById :one
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
JOIN users author ON p.user_id = author.id
LEFT JOIN comments c ON p.id = c.post_id
LEFT JOIN users u ON c.user_id = u.id
WHERE p.id = $1
GROUP BY p.id, author.username;

-- name: DeletePostById :exec
DELETE FROM posts WHERE id = $1;

-- name: UpdatePostById :one
UPDATE posts
SET 
    title = COALESCE($1, title),
    content = COALESCE($2, content),
    tags = COALESCE($3, tags)
WHERE id = $4 AND updated_at = $5
RETURNING id, updated_at;

-- name: GetUserFeed :many
SELECT p.id, p.title, p.content, p.tags, p.created_at, p.updated_at, 
    u.username,
    COUNT(c.id) AS comments_count
FROM posts p
JOIN users u ON p.user_id = u.id
LEFT JOIN comments c ON p.id = c.post_id
JOIN follows f ON p.user_id = f.follow_id OR p.user_id = $1
WHERE 
    f.user_id = $1 AND
    (p.title ILIKE '%' || $2::TEXT || '%' OR p.content ILIKE '%' || $2::TEXT || '%') AND 
    (p.tags @> $3 OR $3 = '{}')
GROUP BY p.id, u.username
ORDER BY p.created_at DESC
LIMIT $4 OFFSET $5;