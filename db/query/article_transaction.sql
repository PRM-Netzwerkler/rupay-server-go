-- name: CreateArticleTransaction :one
INSERT INTO article_transaction (
    article_uuid,
    transaction_uuid,
    amount,
    price
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetArticleTransactionById :one
SELECT * FROM article_transaction
WHERE uuid = $1 LIMIT 1;

-- name: GetArticleTransactions :many
SELECT * FROM article_transaction;

-- name: UpdateArticleTransaction :one
UPDATE article_transaction
SET
    article_uuid = COALESCE(sqlc.narg('article_uuid'), article_uuid),
    transaction_uuid = COALESCE(sqlc.narg('transaction_uuid'), transaction_uuid),
    amount = COALESCE(sqlc.narg('amount'), amount),
    price = COALESCE(sqlc.narg('price'), price)
WHERE uuid = sqlc.arg('uuid')
RETURNING *;

-- name: DeleteArticleTransaction :exec
DELETE FROM article_transaction
WHERE uuid = $1;

-- name: GetArticleTransactionsGroupedByArticle :many
select article_uuid, sum(amount) as amount from article_transaction
group by article_uuid;
