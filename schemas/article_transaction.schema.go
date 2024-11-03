package schemas

import (
	"github.com/google/uuid"
	"github.com/guregu/null/v5"
)

type CreateArticleTransaction struct {
	ArticleUuid     uuid.UUID `json:"article_uuid" binding:"required"`
	TransactionUuid uuid.UUID `json:"transaction_uuid" binding:"required"`
	Amount          int32     `json:"amount" binding:"required"`
	Price           float64   `json:"price" binding:"required"`
}

type UpdateArticleTransaction struct {
	ArticleUuid     uuid.NullUUID `json:"article_uuid"`
	TransactionUuid uuid.NullUUID `json:"transaction_uuid"`
	Amount          null.Int32    `json:"amount"`
	Price           null.Float    `json:"price"`
}
