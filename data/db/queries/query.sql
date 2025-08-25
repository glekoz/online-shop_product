-- name: Create :one
INSERT INTO product(
    id,
    user_id,
    name,
    description,
    price
) VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- imname: LimitImageNumber :one
--INSERT INTO product_image(product_id, max_count)
--VALUES ($1, $2)
--RETURNING *;

-- imname: GetOne :one
--SELECT p.id, p.user_id, p.name, p.description, p.price, p.created_at, pi.max_count
--FROM product AS p INNER JOIN product_image AS pi ON p.id = pi.product_id
--WHERE p.id = $1;

-- name: GetAll :many
SELECT *
FROM product;

-- name: GetLatest :many
SELECT *
FROM product
ORDER BY created_at DESC
LIMIT 10;

-- name: Delete :exec
DELETE 
FROM product
WHERE id = $1;

-- name: DeleteByUser :exec
DELETE
FROM product
WHERE user_id = $1;