-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, password_hash, email)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: FindUser :one
SELECT * FROM users
WHERE email = $1;

-- name: FindUserByUUID :one
SELECT * FROM users
WHERE id = $1;

-- name: UpdateUserDetails :exec
UPDATE users SET email = $1, password_hash = $2, updated_at = NOW()
WHERE id = $3;