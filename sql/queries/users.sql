-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, password_hash, email, is_chirpy_red)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    false
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

-- name: UpgradeToRed :exec
UPDATE users SET is_chirpy_red = true
WHERE id = $1;