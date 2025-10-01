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

-- name: Delete :execrows
DELETE
FROM products
WHERE id = $1;

-- тут сначала будет Гет, потом из переданных в функцию 
-- аргументов выбираются ненулевые и заменяются в структуре из Гет
-- и отправляются в БД

-- name: Update :execrows
UPDATE products
SET name = $2, price = $3, description = $4
WHERE id = $1;

-- name: OrderedOffsetGetAll :many
SELECT id, name, price, description
FROM products
ORDER BY $1
LIMIT $2
OFFSET $3;