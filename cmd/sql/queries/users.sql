-- name: CreateUser :one
INSERT 
INTO users (username, email, password) 
VALUES ($1, $2, $3)
RETURNING id;

-- name: GetUserByUserId :one
SELECT * 
FROM users 
WHERE id = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * 
FROM users 
WHERE  username = $1 LIMIT 1;