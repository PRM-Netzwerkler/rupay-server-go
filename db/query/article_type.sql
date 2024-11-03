-- name: CreateArticleType :one
INSERT INTO article_type (
    "name",
    "desc",
    "icon_codepoint",
    color
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetArticleTypeById :one
SELECT * FROM article_type
WHERE uuid = $1 LIMIT 1;

-- name: GetArticleTypes :many
SELECT * FROM article_type;

-- name: GetArticleTypesWithArticles :many
select sqlc.embed(article_type), article.* from article_type left join article on article.article_type_uuid = article_type.uuid;

-- name: UpdateArticleType :one
UPDATE article_type
SET
    "name" = COALESCE(sqlc.narg('name'), "name"),
    "desc" = COALESCE(sqlc.narg('desc'), "desc"),
    "icon_codepoint" = COALESCE(sqlc.narg('icon_codepoint'), "icon_codepoint"),
    "color" = COALESCE(sqlc.narg('color'), "color")
WHERE uuid = sqlc.arg('uuid')
RETURNING *;

-- name: DeleteArticleType :exec
DELETE FROM article_type
WHERE uuid = $1;
