package schemas

import (
	db "github.com/KevinGruber2001/rupay-bar-backend/db/sqlc"
	"github.com/guregu/null/v5"
)

type CreateArticleType struct {
	Name          string      `json:"name" binding:"required"`
	Desc          null.String `json:"desc"`
	IconCodepoint int32       `json:"icon_codepoint binding:"required"`
	Color         string      `json:"color" binding:"required"`
}

type UpdateArticleType struct {
	Name          null.String `json:"name"`
	Desc          null.String `json:"desc"`
	IconCodepoint null.Int32  `json:"icon_codepoint`
	Color         null.String `json:"color"`
}

type ArticleTypeWithArticles struct {
	db.ArticleType
	Articles []db.Article `json:"articles"`
}
