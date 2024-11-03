package controllers

import (
	"context"
	"database/sql"
	"net/http"

	db "github.com/KevinGruber2001/rupay-bar-backend/db/sqlc"
	e "github.com/KevinGruber2001/rupay-bar-backend/errors"
	"github.com/KevinGruber2001/rupay-bar-backend/schemas"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ArticleTypeController struct {
	db  *db.Queries
	ctx context.Context
}

func NewArticleTypeController(db *db.Queries, ctx context.Context) *ArticleTypeController {
	return &ArticleTypeController{db, ctx}
}

// CreateArticleType godoc
// @Summary Create a new article type
// @Description Create a new article type with the provided payload
// @Tags ArticleTypes
// @Accept json
// @Produce json
// @Param articleType body schemas.CreateArticleType true "Create article type payload"
// @Success 200 {object} db.ArticleType "Successfully created article type"
// @Failure 400 {object} e.ErrorResponse "Invalid payload"
// @Failure 500 {object} e.ErrorResponse "Failed to create article type"
// @Router /article-type [post]
func (cc *ArticleTypeController) CreateArticleType(ctx *gin.Context) {
	var payload *schemas.CreateArticleType

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, e.ErrorResponse{Code: e.InvalidPayload, Message: "Invalid Payload", Error: err.Error()})
		return
	}

	args := &db.CreateArticleTypeParams{
		Name:          payload.Name,
		Desc:          payload.Desc,
		IconCodepoint: payload.IconCodepoint,
		Color:         payload.Color,
	}

	articleType, err := cc.db.CreateArticleType(ctx, *args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to create ArticleType", Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, articleType)
}

// UpdateArticleType godoc
// @Summary Update an existing article type
// @Description Update an article type with the provided ID and payload
// @Tags ArticleTypes
// @Accept json
// @Produce json
// @Param articleTypeId path string true "Article Type ID"
// @Param articleType body schemas.UpdateArticleType true "Update article type payload"
// @Success 200 {object} db.ArticleType "Successfully updated article type"
// @Failure 400 {object} e.ErrorResponse "Invalid payload"
// @Failure 404 {object} e.ErrorResponse "Article type not found"
// @Failure 500 {object} e.ErrorResponse "Failed to update article type"
// @Router /article-type/{articleTypeId} [put]
func (cc *ArticleTypeController) UpdateArticleType(ctx *gin.Context) {
	var payload *schemas.UpdateArticleType
	articleTypeId := ctx.Param("articleTypeId")

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, e.ErrorResponse{Code: e.InvalidPayload, Message: "Invalid Payload", Error: err.Error()})
		return
	}

	args := &db.UpdateArticleTypeParams{
		Uuid:          uuid.MustParse(articleTypeId),
		Name:          payload.Name,
		Desc:          payload.Desc,
		IconCodepoint: payload.IconCodepoint,
		Color:         payload.Color,
	}

	articleType, err := cc.db.UpdateArticleType(ctx, *args)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, e.ErrorResponse{Code: e.NotFound, Message: "Article type not found", Error: err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to update ArticleType", Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, articleType)
}

// GetArticleTypeById godoc
// @Summary Get an article type by ID
// @Description Retrieve an article type using its ID
// @Tags ArticleTypes
// @Produce json
// @Param articleTypeId path string true "Article Type ID"
// @Success 200 {object} db.ArticleType "Successfully retrieved article type"
// @Failure 404 {object} e.ErrorResponse "Article type not found"
// @Failure 500 {object} e.ErrorResponse "Failed to retrieve article type"
// @Router /article-type/{articleTypeId} [get]
func (cc *ArticleTypeController) GetArticleTypeById(ctx *gin.Context) {
	articleTypeId := ctx.Param("articleTypeId")

	articleType, err := cc.db.GetArticleTypeById(ctx, uuid.MustParse(articleTypeId))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, e.ErrorResponse{Code: e.NotFound, Message: "Article type not found", Error: err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to retrieve ArticleType", Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, articleType)
}

// GetAllArticleTypes godoc
// @Summary Retrieve all article types
// @Description Get a list of all article types
// @Tags ArticleTypes
// @Produce json
// @Success 200 {array} db.ArticleType "Successfully retrieved all article types"
// @Failure 500 {object} e.ErrorResponse "Failed to retrieve article types"
// @Router /article-type [get]
func (cc *ArticleTypeController) GetAllArticleTypes(ctx *gin.Context) {
	articleTypes, err := cc.db.GetArticleTypes(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to retrieve ArticleTypes", Error: err.Error()})
		return
	}

	if articleTypes == nil {
		articleTypes = []db.ArticleType{}
	}

	ctx.JSON(http.StatusOK, articleTypes)
}

// GetAllArticleTypes godoc
// @Summary Retrieve all article types
// @Description Get a list of all article types
// @Tags ArticleTypes
// @Produce json
// @Success 200 {array} db.ArticleType "Successfully retrieved all article types"
// @Failure 500 {object} e.ErrorResponse "Failed to retrieve article types"
// @Router /article-type [get]
func (cc *ArticleTypeController) GetAllArticleTypesWithArticles(ctx *gin.Context) {
	articleTypes, err := cc.db.GetArticleTypesWithArticles(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to retrieve ArticleTypes with Articles", Error: err.Error()})
		return
	}

	if articleTypes == nil {
		// if empty just return empty array
		articleTypes = []db.GetArticleTypesWithArticlesRow{}
		ctx.JSON(http.StatusOK, articleTypes)
		return
	}

	// at this point the parsing begins

	result := []schemas.ArticleTypeWithArticles{}

	currentAtwa := schemas.ArticleTypeWithArticles{}

	for _, articleType := range articleTypes {
		if articleType.ArticleType.Uuid == currentAtwa.Uuid {
			if articleType.Uuid.Valid {
				currentAtwa.Articles = append(currentAtwa.Articles, db.Article{Uuid: articleType.Uuid.UUID, Name: articleType.Name.String, Desc: articleType.Desc, PurchasePrice: articleType.PurchasePrice.Float64, ResellPrice: articleType.ResellPrice.Float64, ArticleTypeUuid: articleType.Uuid.UUID})
			}

		} else {
			// add currentatwa to result, and make new article type to current atwa
			if currentAtwa.Uuid != uuid.Nil {
				result = append(result, currentAtwa)
			}
			currentAtwa = schemas.ArticleTypeWithArticles{ArticleType: articleType.ArticleType, Articles: []db.Article{}}
			if articleType.Uuid.Valid {
				currentAtwa.Articles = append(currentAtwa.Articles, db.Article{Uuid: articleType.Uuid.UUID, Name: articleType.Name.String, Desc: articleType.Desc, PurchasePrice: articleType.PurchasePrice.Float64, ResellPrice: articleType.ResellPrice.Float64, ArticleTypeUuid: articleType.Uuid.UUID})
			}
		}

	}

	ctx.JSON(http.StatusOK, result)
}

// DeleteArticleTypeById godoc
// @Summary Delete an article type by ID
// @Description Remove an article type using its ID
// @Tags ArticleTypes
// @Produce json
// @Param articleTypeId path string true "Article Type ID"
// @Success 204 "Successfully deleted article type"
// @Failure 404 {object} e.ErrorResponse "Article type not found"
// @Failure 500 {object} e.ErrorResponse "Failed to delete article type"
// @Router /article-type/{articleTypeId} [delete]
func (cc *ArticleTypeController) DeleteArticleTypeById(ctx *gin.Context) {
	articleTypeId := ctx.Param("articleTypeId")

	_, err := cc.db.GetArticleTypeById(ctx, uuid.MustParse(articleTypeId))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, e.ErrorResponse{Code: e.NotFound, Message: "Article type not found", Error: err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to retrieve ArticleType", Error: err.Error()})
		return
	}

	err = cc.db.DeleteArticleType(ctx, uuid.MustParse(articleTypeId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to delete ArticleType", Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
