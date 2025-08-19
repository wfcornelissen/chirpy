-- name: CreateToken :one
INSERT INTO refresh_tokens (
    token,
    created_at,
    updated_at,
    user_id,
    expires_at,
    revoked_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
) RETURNING *;

-- name: GetToken :one
SELECT * FROM refresh_tokens
WHERE token = $1 AND expires_at > now()
LIMIT 1;

-- name: DeleteToken :exec
DELETE FROM refresh_tokens
WHERE token = $1;