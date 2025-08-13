-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL DEFAULT 'unset'
);


-- +goose Down
DROP TABLE users;
