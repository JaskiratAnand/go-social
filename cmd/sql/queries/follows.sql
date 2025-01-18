-- name: FollowUser :exec
INSERT 
INTO follows (user_id, follow_id) 
VALUES ($1, $2);

-- name: UnfollowUser :exec
DELETE FROM follows 
WHERE user_id = $1 AND follow_id = $2;