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

type TransactionController struct {
	db  *db.Queries
	ctx context.Context
}

func NewTransactionController(db *db.Queries, ctx context.Context) *TransactionController {
	return &TransactionController{db, ctx}
}

// @Summary Create a new transaction
// @Description Create a new transaction with the provided date and price
// @Tags Transactions
// @Accept json
// @Produce json
// @Param payload body schemas.CreateTransaction true "CreateTransaction payload"
// @Success 200 {object} db.Transaction "Transaction data"
// @Failure 400 {object} e.ErrorResponse "Invalid Payload"
// @Router /transaction [post]
func (cc *TransactionController) CreateTransaction(ctx *gin.Context) {
	var payload *schemas.CreateTransaction

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, e.ErrorResponse{Code: e.InvalidPayload, Message: "Payload is invalid", Error: err.Error()})
		return
	}

	args := &db.CreateTransactionParams{
		Date:  payload.Date,
		Price: payload.Price,
	}

	Transaction, err := cc.db.CreateTransaction(ctx, *args)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to create the transaction", Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, Transaction)
}

// @Summary Update a transaction
// @Description Update a transaction with the provided date and price
// @Tags Transactions
// @Accept json
// @Produce json
// @Param payload body schemas.UpdateTransaction true "UpdateTransaction payload"
// @Success 200 {object} db.Transaction "Transaction data"
// @Failure 400 {object} e.ErrorResponse "Invalid Payload"
// @Router /transaction [patch]
func (cc *TransactionController) UpdateTransaction(ctx *gin.Context) {
	var payload *schemas.UpdateTransaction
	TransactionId := ctx.Param("transactionId")

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, e.ErrorResponse{Code: e.InvalidPayload, Message: "Payload is invalid", Error: err.Error()})
		return
	}

	args := &db.UpdateTransactionParams{
		Uuid:  uuid.MustParse(TransactionId),
		Date:  payload.Date,
		Price: payload.Price,
	}

	Transaction, err := cc.db.UpdateTransaction(ctx, *args)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, e.ErrorResponse{Code: e.NotFound, Message: "Failed to find Transaction", Error: err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to retrieve Transaction", Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, Transaction)
}

// @Summary Retrieve a transaction
// @Description Retrieve a transaction by the provided id
// @Tags Transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} db.Transaction "Transaction data"
// @Failure 500 {object} e.ErrorResponse "Internal Server Error"
// @Failure 404 {object} e.ErrorResponse "Transaction not found"
// @Router /transaction/{id} [get]
func (cc *TransactionController) GetTransactionById(ctx *gin.Context) {
	TransactionId := ctx.Param("transactionId")

	Transaction, err := cc.db.GetTransactionById(ctx, uuid.MustParse(TransactionId))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, e.ErrorResponse{Code: e.NotFound, Message: "Transaction not found", Error: err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to retrieve this Transaction", Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, Transaction)
}

// @Summary Retrieve all transactions
// @Description Retrieve a list of all transactions
// @Tags Transactions
// @Accept json
// @Produce json
// @Success 200 {array} db.Transaction "List of transactions"
// @Failure 500 {object} e.ErrorResponse "Internal Server Error"
// @Router /transaction [get]
func (cc *TransactionController) GetAllTransactions(ctx *gin.Context) {

	Transactions, err := cc.db.GetTransactions(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to retrieve Transactions", Error: err.Error()})
		return
	}

	if Transactions == nil {
		Transactions = []db.Transaction{}
	}

	ctx.JSON(http.StatusOK, Transactions)
}

// @Summary Delete a transaction
// @Description Delete a transaction by the provided id
// @Tags Transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 204 "Transaction deleted successfully"
// @Failure 500 {object} e.ErrorResponse "Failed to delete Transaction"
// @Failure 404 {object} e.ErrorResponse "Transaction not found"
// @Router /transaction/{id} [delete]
func (cc *TransactionController) DeleteTransactionById(ctx *gin.Context) {
	TransactionId := ctx.Param("transactionId")

	_, err := cc.db.GetTransactionById(ctx, uuid.MustParse(TransactionId))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, e.ErrorResponse{Code: e.NotFound, Message: "Failed to retrieve Transaction", Error: err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to retrieve Transaction", Error: err.Error()})
		return
	}

	err = cc.db.DeleteTransaction(ctx, uuid.MustParse(TransactionId))
	if err != nil {
		ctx.JSON(http.StatusBadGateway, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to delete Transaction", Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)

}
