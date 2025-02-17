-- name: CreateUser :one
INSERT 
INTO users (username, email, password, role_id) 
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: GetUserByUserId :one
SELECT * 
FROM users 
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * 
FROM users 
WHERE email = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * 
FROM users 
WHERE username = $1 LIMIT 1;

-- name: ActivateUser :exec 
UPDATE users
SET verified = true
WHERE id = $1;

-- name: DeleteUser :exec
DELETE
FROM users
WHERE id = $1;