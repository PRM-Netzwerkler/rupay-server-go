-- name: CreateEvent :one
INSERT INTO event (
    "name",
    "desc",
    from_date,
    to_date
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetEventById :one
SELECT * FROM event
WHERE uuid = $1 LIMIT 1;

-- name: GetEvents :many
SELECT * FROM event;

-- name: UpdateEvent :one
UPDATE event
SET
    "name" = COALESCE(sqlc.narg('name'), "name"),
    "desc" = COALESCE(sqlc.narg('desc'), "desc"),
    from_date = COALESCE(sqlc.narg('from_date'), from_date),
    to_date = COALESCE(sqlc.narg('to_date'), to_date)
WHERE uuid = sqlc.arg('uuid')
RETURNING *;

-- name: DeleteEvent :exec
DELETE FROM event
WHERE uuid = $1;
