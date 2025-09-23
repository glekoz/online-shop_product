-- Я использую pgx, в котором если не найдена строка
-- при запросе одной строки, то возвращается ошибка;
-- а если запрашивается несколько строк и не возвращается
-- ни одной, то ошибки нет - будет пустой срез.

-- name: Create :exec
INSERT INTO products(id, name, price, description)
VALUES ($1, $2, $3, $4);

-- name: Get :one
SELECT name, price, description
FROM products
WHERE id = $1;

-- name: GetAll :many
SELECT id, name, price
FROM products;
