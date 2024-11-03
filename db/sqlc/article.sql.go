// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: article.sql

package db

import (
	"context"

	"github.com/google/uuid"
	null "github.com/guregu/null/v5"
)

const createArticle = `-- name: CreateArticle :one
INSERT INTO article (
    "name",
    "desc",
    purchase_price,
    resell_price,
    article_type_uuid
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING uuid, name, "desc", purchase_price, resell_price, article_type_uuid
`

type CreateArticleParams struct {
	Name            string      `json:"name"`
	Desc            null.String `json:"desc"`
	PurchasePrice   float64     `json:"purchase_price"`
	ResellPrice     float64     `json:"resell_price"`
	ArticleTypeUuid uuid.UUID   `json:"article_type_uuid"`
}

func (q *Queries) CreateArticle(ctx context.Context, arg CreateArticleParams) (Article, error) {
	row := q.queryRow(ctx, q.createArticleStmt, createArticle,
		arg.Name,
		arg.Desc,
		arg.PurchasePrice,
		arg.ResellPrice,
		arg.ArticleTypeUuid,
	)
	var i Article
	err := row.Scan(
		&i.Uuid,
		&i.Name,
		&i.Desc,
		&i.PurchasePrice,
		&i.ResellPrice,
		&i.ArticleTypeUuid,
	)
	return i, err
}

const deleteArticle = `-- name: DeleteArticle :exec
DELETE FROM article
WHERE uuid = $1
`

func (q *Queries) DeleteArticle(ctx context.Context, argUuid uuid.UUID) error {
	_, err := q.exec(ctx, q.deleteArticleStmt, deleteArticle, argUuid)
	return err
}

const getArticleById = `-- name: GetArticleById :one
SELECT uuid, name, "desc", purchase_price, resell_price, article_type_uuid FROM article
WHERE uuid = $1 LIMIT 1
`

func (q *Queries) GetArticleById(ctx context.Context, argUuid uuid.UUID) (Article, error) {
	row := q.queryRow(ctx, q.getArticleByIdStmt, getArticleById, argUuid)
	var i Article
	err := row.Scan(
		&i.Uuid,
		&i.Name,
		&i.Desc,
		&i.PurchasePrice,
		&i.ResellPrice,
		&i.ArticleTypeUuid,
	)
	return i, err
}

const getArticles = `-- name: GetArticles :many
SELECT uuid, name, "desc", purchase_price, resell_price, article_type_uuid FROM article
`

func (q *Queries) GetArticles(ctx context.Context) ([]Article, error) {
	rows, err := q.query(ctx, q.getArticlesStmt, getArticles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Article{}
	for rows.Next() {
		var i Article
		if err := rows.Scan(
			&i.Uuid,
			&i.Name,
			&i.Desc,
			&i.PurchasePrice,
			&i.ResellPrice,
			&i.ArticleTypeUuid,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateArticle = `-- name: UpdateArticle :one
UPDATE article
SET
    "name" = COALESCE($1, "name"),
    "desc" = COALESCE($2, "desc"),
    purchase_price = COALESCE($3, purchase_price),
    resell_price = COALESCE($4, resell_price),
    article_type_uuid = COALESCE($5, article_type_uuid)
WHERE uuid = $6
RETURNING uuid, name, "desc", purchase_price, resell_price, article_type_uuid
`

type UpdateArticleParams struct {
	Name            null.String   `json:"name"`
	Desc            null.String   `json:"desc"`
	PurchasePrice   null.Float    `json:"purchase_price"`
	ResellPrice     null.Float    `json:"resell_price"`
	ArticleTypeUuid uuid.NullUUID `json:"article_type_uuid"`
	Uuid            uuid.UUID     `json:"uuid"`
}

func (q *Queries) UpdateArticle(ctx context.Context, arg UpdateArticleParams) (Article, error) {
	row := q.queryRow(ctx, q.updateArticleStmt, updateArticle,
		arg.Name,
		arg.Desc,
		arg.PurchasePrice,
		arg.ResellPrice,
		arg.ArticleTypeUuid,
		arg.Uuid,
	)
	var i Article
	err := row.Scan(
		&i.Uuid,
		&i.Name,
		&i.Desc,
		&i.PurchasePrice,
		&i.ResellPrice,
		&i.ArticleTypeUuid,
	)
	return i, err
}
