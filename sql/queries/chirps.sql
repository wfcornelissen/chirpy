-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, email)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1
)
RETURNING *;

-- name: DeleteAllChirps :exec
DELETE FROM users;