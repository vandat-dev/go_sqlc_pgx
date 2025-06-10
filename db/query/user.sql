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

-- name: ListUsersWithProducts :many
SELECT
    u.id as user_id,
    u.name as user_name,
    u.phone as user_phone,
    u.email as user_email,
    u.created_at as user_created_at,
    p.id as product_id,
    p.name as product_name,
    p.description as product_description,
    p.price as product_price,
    p.created_at as product_created_at,
    p.updated_at as product_updated_at
FROM users u
LEFT JOIN products p ON u.id = p.user_id
ORDER BY u.created_at DESC, p.created_at DESC;