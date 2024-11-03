package schemas

import (
	"github.com/google/uuid"
	"github.com/guregu/null/v5"
)

type CreateArticle struct {
	Name            string      `json:"name" binding:"required"`
	Desc            null.String `json:"desc"`
	PurchasePrice   float64     `json:"purchase_price" binding:"required"`
	ResellPrice     float64     `json:"resell_price" binding:"required"`
	ArticleTypeUuid uuid.UUID   `json:"article_type_uuid" binding:"required"`
}

type UpdateArticle struct {
	Name            null.String   `json:"name"`
	Desc            null.String   `json:"desc"`
	PurchasePrice   null.Float    `json:"purchase_price"`
	ResellPrice     null.Float    `json:"resell_price"`
	ArticleTypeUuid uuid.NullUUID `json:"article_type_uuid"`
}
