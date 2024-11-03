-- name: CreateTransaction :one
INSERT INTO transaction (
    "date",
    price
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetTransactionById :one
SELECT * FROM transaction
WHERE uuid = $1 LIMIT 1;

-- name: GetTransactions :many
SELECT * FROM transaction;

-- name: UpdateTransaction :one
UPDATE transaction
SET
    "date" = COALESCE(sqlc.narg('date'), "date"),
    price = COALESCE(sqlc.narg('price'), price)
WHERE uuid = sqlc.arg('uuid')
RETURNING *;

-- name: DeleteTransaction :exec
DELETE FROM transaction
WHERE uuid = $1;
