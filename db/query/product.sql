-- name: CreateProduct :one
INSERT INTO products (name, description, price, user_id)
VALUES ($1, $2, $3, $4)
RETURNING id, name, description, price, user_id, created_at, updated_at;

-- name: ListProducts :many
SELECT *
FROM products
ORDER BY created_at DESC;

-- name: ListProductsWithUsers :many
SELECT 
    p.id,
    p.name,
    p.description,
    p.price,
    p.user_id,
    p.created_at,
    p.updated_at,
    u.name as user_name,
    u.email as user_email,
    u.phone as user_phone
FROM products p
INNER JOIN users u ON p.user_id = u.id
ORDER BY p.created_at DESC;

-- name: GetProductByID :one
SELECT *
FROM products
WHERE id = $1;

-- name: GetProductWithUserByID :one
SELECT 
    p.id,
    p.name,
    p.description,
    p.price,
    p.user_id,
    p.created_at,
    p.updated_at,
    u.name as user_name,
    u.email as user_email,
    u.phone as user_phone
FROM products p
INNER JOIN users u ON p.user_id = u.id
WHERE p.id = $1;

-- name: GetProductsByUserID :many
SELECT *
FROM products
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateProduct :one
UPDATE products
SET name = $2, description = $3, price = $4, updated_at = NOW()
WHERE id = $1
RETURNING id, name, description, price, user_id, created_at, updated_at;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1; 