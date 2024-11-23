-- name: CreateUser :one
INSERT INTO resident (
    "name",
    code
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetUserById :one
SELECT * FROM resident
WHERE name = $1 LIMIT 1;

-- name: GetUserByCode :one
SELECT * FROM resident
WHERE code = $1 LIMIT 1;

-- name: GetUsers :many
SELECT * FROM resident;

-- name: UpdateUser :one
UPDATE resident
SET
    name = COALESCE(sqlc.narg('name'), "name"),
    code = COALESCE(sqlc.narg('code'), code)
WHERE name = sqlc.arg('name')
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM resident
WHERE name = $1;
