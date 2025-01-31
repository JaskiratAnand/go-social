-- name: CreateInvitation :exec
INSERT 
INTO user_invitations (token, user_id, expiary)
VALUES ($1, $2, $3)
RETURNING token;

-- name: GetInvitationByToken :one
SELECT *
FROM user_invitations
WHERE token = $1
LIMIT 1;

-- name: DeleteInvitationByUserId :exec
DELETE
FROM user_invitations
WHERE user_id = $1;