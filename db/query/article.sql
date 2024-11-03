-- name: CreateArticle :one
INSERT INTO article (
    "name",
    "desc",
    purchase_price,
    resell_price,
    article_type_uuid
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetArticleById :one
SELECT * FROM article
WHERE uuid = $1 LIMIT 1;

-- name: GetArticles :many
SELECT * FROM article;

-- name: UpdateArticle :one
UPDATE article
SET
    "name" = COALESCE(sqlc.narg('name'), "name"),
    "desc" = COALESCE(sqlc.narg('desc'), "desc"),
    purchase_price = COALESCE(sqlc.narg('purchase_price'), purchase_price),
    resell_price = COALESCE(sqlc.narg('resell_price'), resell_price),
    article_type_uuid = COALESCE(sqlc.narg('article_type_uuid'), article_type_uuid)
WHERE uuid = sqlc.arg('uuid')
RETURNING *;

-- name: DeleteArticle :exec
DELETE FROM article
WHERE uuid = $1;
