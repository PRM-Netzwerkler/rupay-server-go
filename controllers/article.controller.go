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

type ArticleController struct {
	db  *db.Queries
	ctx context.Context
}

func NewArticleController(db *db.Queries, ctx context.Context) *ArticleController {
	return &ArticleController{db, ctx}
}

// CreateArticle godoc
// @Summary Create a new article
// @Description Create a new article with the provided payload
// @Tags Articles
// @Accept json
// @Produce json
// @Param article body schemas.CreateArticle true "Create article payload"
// @Success 200 {object} db.Article "Successfully created article"
// @Failure 400 {object} e.ErrorResponse "Invalid payload"
// @Failure 500 {object} e.ErrorResponse "Failed to create article"
// @Router /article [post]
func (cc *ArticleController) CreateArticle(ctx *gin.Context) {
	var payload *schemas.CreateArticle

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, e.ErrorResponse{Code: e.InvalidPayload, Message: "Invalid Payload", Error: err.Error()})
		return
	}

	args := &db.CreateArticleParams{
		Name:            payload.Name,
		Desc:            payload.Desc,
		PurchasePrice:   payload.PurchasePrice,
		ResellPrice:     payload.ResellPrice,
		ArticleTypeUuid: payload.ArticleTypeUuid,
	}

	article, err := cc.db.CreateArticle(ctx, *args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to create Article", Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, article)
}

// UpdateArticle godoc
// @Summary Update an existing article
// @Description Update an article with the provided ID and payload
// @Tags Articles
// @Accept json
// @Produce json
// @Param articleId path string true "Article ID"
// @Param article body schemas.UpdateArticle true "Update article payload"
// @Success 200 {object} db.Article "Successfully updated article"
// @Failure 400 {object} e.ErrorResponse "Invalid payload"
// @Failure 404 {object} e.ErrorResponse "Article not found"
// @Failure 500 {object} e.ErrorResponse "Failed to update article"
// @Router /articles/{articleId} [put]
func (cc *ArticleController) UpdateArticle(ctx *gin.Context) {
	var payload *schemas.UpdateArticle
	articleId := ctx.Param("articleId")

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, e.ErrorResponse{Code: e.InvalidPayload, Message: "Invalid Payload", Error: err.Error()})
		return
	}

	args := &db.UpdateArticleParams{
		Uuid:            uuid.MustParse(articleId),
		Name:            payload.Name,
		Desc:            payload.Desc,
		PurchasePrice:   payload.PurchasePrice,
		ResellPrice:     payload.ResellPrice,
		ArticleTypeUuid: payload.ArticleTypeUuid,
	}

	article, err := cc.db.UpdateArticle(ctx, *args)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, e.ErrorResponse{Code: e.NotFound, Message: "Article not found", Error: err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to update Article", Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, article)
}

// GetArticleById godoc
// @Summary Get an article by ID
// @Description Retrieve an article using its ID
// @Tags Articles
// @Produce json
// @Param articleId path string true "Article ID"
// @Success 200 {object} db.Article "Successfully retrieved article"
// @Failure 404 {object} e.ErrorResponse "Article not found"
// @Failure 500 {object} e.ErrorResponse "Failed to retrieve article"
// @Router /articles/{articleId} [get]
func (cc *ArticleController) GetArticleById(ctx *gin.Context) {
	articleId := ctx.Param("articleId")

	article, err := cc.db.GetArticleById(ctx, uuid.MustParse(articleId))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, e.ErrorResponse{Code: e.NotFound, Message: "Article not found", Error: err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to retrieve Article", Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, article)
}

// GetAllArticles godoc
// @Summary Retrieve all articles
// @Description Get a list of all articles
// @Tags Articles
// @Produce json
// @Success 200 {array} db.Article "Successfully retrieved all articles"
// @Failure 500 {object} e.ErrorResponse "Failed to retrieve articles"
// @Router /articles [get]
func (cc *ArticleController) GetAllArticles(ctx *gin.Context) {
	articles, err := cc.db.GetArticles(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to retrieve Articles", Error: err.Error()})
		return
	}

	if articles == nil {
		articles = []db.Article{}
	}

	ctx.JSON(http.StatusOK, articles)
}

// DeleteArticleById godoc
// @Summary Delete an article by ID
// @Description Remove an article using its ID
// @Tags Articles
// @Produce json
// @Param articleId path string true "Article ID"
// @Success 204 "Successfully deleted article"
// @Failure 404 {object} e.ErrorResponse "Article not found"
// @Failure 500 {object} e.ErrorResponse "Failed to delete article"
// @Router /articles/{articleId} [delete]
func (cc *ArticleController) DeleteArticleById(ctx *gin.Context) {
	articleId := ctx.Param("articleId")

	_, err := cc.db.GetArticleById(ctx, uuid.MustParse(articleId))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, e.ErrorResponse{Code: e.NotFound, Message: "Article not found", Error: err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to retrieve Article", Error: err.Error()})
		return
	}

	err = cc.db.DeleteArticle(ctx, uuid.MustParse(articleId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to delete Article", Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
