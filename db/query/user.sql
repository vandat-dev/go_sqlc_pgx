-- name: CreateUser :one
INSERT INTO users (name, phone, email)
VALUES ($1, $2, $3) RETURNING id, name, phone ,email, created_at;

-- name: ListUsers :many
SELECT *
FROM users
ORDER BY created_at DESC;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1;