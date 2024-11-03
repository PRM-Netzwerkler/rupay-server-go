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

type ArticleTransactionController struct {
	db  *db.Queries
	ctx context.Context
}

func NewArticleTransactionController(db *db.Queries, ctx context.Context) *ArticleTransactionController {
	return &ArticleTransactionController{db, ctx}
}

// @Summary Create a new article transaction
// @Description Create a new article transaction with the provided payload
// @Tags ArticleTransactions
// @Accept json
// @Produce json
// @Param payload body db.ArticleTransaction true "CreateArticleTransaction payload"
// @Success 200 {object} db.ArticleTransaction "ArticleTransaction data"
// @Failure 400 {object} e.ErrorResponse "Invalid Payload"
// @Router /article-transaction [post]
func (cc *ArticleTransactionController) CreateArticleTransaction(ctx *gin.Context) {
	var payload *schemas.CreateArticleTransaction

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, e.ErrorResponse{Code: e.InvalidPayload, Message: "Payload is invalid", Error: err.Error()})
		return
	}

	args := &db.CreateArticleTransactionParams{
		ArticleUuid:     payload.ArticleUuid,
		TransactionUuid: payload.TransactionUuid,
		Amount:          payload.Amount,
		Price:           payload.Price,
	}

	articleTransaction, err := cc.db.CreateArticleTransaction(ctx, *args)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to create the ArticleTransaction", Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, articleTransaction)
}

// @Summary Update an article transaction
// @Description Update an article transaction with the provided payload
// @Tags ArticleTransactions
// @Accept json
// @Produce json
// @Param articleTransactionId path string true "Article Transaction ID"
// @Param payload body schemas.UpdateArticleTransaction true "UpdateArticleTransaction payload"
// @Success 200 {object} db.ArticleTransaction "ArticleTransaction data"
// @Failure 400 {object} e.ErrorResponse "Invalid Payload"
// @Failure 404 {object} e.ErrorResponse "Article transaction not found"
// @Router /article-transaction/{articleTransactionId} [put]
func (cc *ArticleTransactionController) UpdateArticleTransaction(ctx *gin.Context) {
	var payload *schemas.UpdateArticleTransaction
	articleTransactionId := ctx.Param("articleTransactionId")

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, e.ErrorResponse{Code: e.InvalidPayload, Message: "Payload is invalid", Error: err.Error()})
		return
	}

	args := &db.UpdateArticleTransactionParams{
		Uuid:            uuid.MustParse(articleTransactionId),
		ArticleUuid:     payload.ArticleUuid,
		TransactionUuid: payload.TransactionUuid,
		Amount:          payload.Amount,
		Price:           payload.Price,
	}

	articleTransaction, err := cc.db.UpdateArticleTransaction(ctx, *args)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, e.ErrorResponse{Code: e.NotFound, Message: "Failed to find ArticleTransaction", Error: err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to update ArticleTransaction", Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, articleTransaction)
}

// @Summary Retrieve an article transaction
// @Description Retrieve an article transaction by the provided ID
// @Tags ArticleTransactions
// @Accept json
// @Produce json
// @Param articleTransactionId path string true "Article Transaction ID"
// @Success 200 {object} db.ArticleTransaction "ArticleTransaction data"
// @Failure 404 {object} e.ErrorResponse "Article transaction not found"
// @Failure 500 {object} e.ErrorResponse "Internal Server Error"
// @Router /article-transaction/{articleTransactionId} [get]
func (cc *ArticleTransactionController) GetArticleTransactionById(ctx *gin.Context) {
	articleTransactionId := ctx.Param("articleTransactionId")

	articleTransaction, err := cc.db.GetArticleTransactionById(ctx, uuid.MustParse(articleTransactionId))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, e.ErrorResponse{Code: e.NotFound, Message: "ArticleTransaction not found", Error: err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to retrieve this ArticleTransaction", Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, articleTransaction)
}

// @Summary Retrieve all article transactions
// @Description Retrieve a list of all article transactions
// @Tags ArticleTransactions
// @Accept json
// @Produce json
// @Success 200 {array} db.ArticleTransaction "List of article transactions"
// @Failure 500 {object} e.ErrorResponse "Internal Server Error"
// @Router /article-transaction [get]
func (cc *ArticleTransactionController) GetAllArticleTransactions(ctx *gin.Context) {
	articleTransactions, err := cc.db.GetArticleTransactions(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to retrieve ArticleTransactions", Error: err.Error()})
		return
	}

	if articleTransactions == nil {
		articleTransactions = []db.ArticleTransaction{}
	}

	ctx.JSON(http.StatusOK, articleTransactions)
}

func (cc *ArticleTransactionController) GetAllArticleTransactionsGroupedByArticle(ctx *gin.Context) {
	articleTransactions, err := cc.db.GetArticleTransactionsGroupedByArticle(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to retrieve ArticleTransactions", Error: err.Error()})
		return
	}

	if articleTransactions == nil {
		articleTransactions = []db.GetArticleTransactionsGroupedByArticleRow{}
	}

	ctx.JSON(http.StatusOK, articleTransactions)
}

// @Summary Delete an article transaction
// @Description Delete an article transaction by the provided ID
// @Tags ArticleTransactions
// @Accept json
// @Produce json
// @Param articleTransactionId path string true "Article Transaction ID"
// @Success 204 "ArticleTransaction deleted successfully"
// @Failure 404 {object} e.ErrorResponse "Article transaction not found"
// @Failure 500 {object} e.ErrorResponse "Failed to delete ArticleTransaction"
// @Router /article-transaction/{articleTransactionId} [delete]
func (cc *ArticleTransactionController) DeleteArticleTransactionById(ctx *gin.Context) {
	articleTransactionId := ctx.Param("articleTransactionId")

	_, err := cc.db.GetArticleTransactionById(ctx, uuid.MustParse(articleTransactionId))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, e.ErrorResponse{Code: e.NotFound, Message: "Failed to retrieve ArticleTransaction", Error: err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to retrieve ArticleTransaction", Error: err.Error()})
		return
	}

	err = cc.db.DeleteArticleTransaction(ctx, uuid.MustParse(articleTransactionId))
	if err != nil {
		ctx.JSON(http.StatusBadGateway, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to delete ArticleTransaction", Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
