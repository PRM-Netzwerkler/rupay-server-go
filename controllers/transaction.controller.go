package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	db "github.com/KevinGruber2001/rupay-bar-backend/db/sqlc"
	e "github.com/KevinGruber2001/rupay-bar-backend/errors"
	"github.com/KevinGruber2001/rupay-bar-backend/schemas"
	"github.com/KevinGruber2001/rupay-bar-backend/util"
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

func (cc *TransactionController) GetSavaPageUser(ctx *gin.Context) {
	// Call Arduino to start reading and wait for the message
	code, err := cc.waitForMqttMessage()

	if err != nil {
		ctx.JSON(http.StatusBadGateway, e.ErrorResponse{Code: e.InternalServerError, Message: "Timeout waiting for reader", Error: err.Error()})
		return
	}

	// Retrieve user by code
	user, err := cc.getUserByCode(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Code does not exist", Error: err.Error()})
		return
	}

	// Get user balance from SavaPage
	balance, err := cc.getSavaPageBalance(user)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to reach Savapage", Error: err.Error()})
		return
	}

	// Return the result
	res := schemas.SavaUser{Name: user.Name, Balance: balance}
	ctx.JSON(http.StatusOK, res)
}

// waitForMqttMessage publishes the "read" message and waits for the response
func (cc *TransactionController) waitForMqttMessage() (string, error) {
	mqtt := util.GetClient()
	mqtt.Publish("read", "bar1/read")

	code, err := mqtt.WaitForMessage("bar1/client", 60*time.Second)
	if err != nil {
		return "", fmt.Errorf("Error waiting for message: %v", err)
	}

	return code, nil
}

// getUserByCode retrieves the user from the database by the provided code
func (cc *TransactionController) getUserByCode(ctx *gin.Context, code string) (db.Resident, error) {
	user, err := cc.db.GetUserByCode(ctx, code)
	if err != nil {
		return db.Resident{}, fmt.Errorf("Error retrieving user: %v", err)
	}

	return user, nil
}

// getSavaPageBalance makes an HTTP request to SavaPage to retrieve the user's balance
func (cc *TransactionController) getSavaPageBalance(user db.Resident) (float64, error) {
	config, err := util.LoadConfig("../.")
	if err != nil {
		return 0, fmt.Errorf("Failed to load config: %v", err)
	}

	url := config.SavaPageUrl + fmt.Sprintf("financial/account/balance?type=USER&name=%s", user.Name)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("Failed to create HTTP request: %v", err)
	}

	req.SetBasicAuth(config.SavaPageAdmin, config.SavaPagePassword)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("Failed to make HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Failed to reach SavaPage, status code: %d", resp.StatusCode)
	}

	// Read and parse the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("Error reading response body: %v", err)
	}

	var response schemas.SavaResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return 0, fmt.Errorf("Error decoding JSON response: %v", err)
	}

	if !response.Success {
		return 0, fmt.Errorf("SavaPage returned unsuccessful response")
	}

	if response.Result.Valid {
		return response.Result.Float64, nil
	}

	return 0, nil
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

	UserName := ctx.Param("username")

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, e.ErrorResponse{Code: e.InvalidPayload, Message: "Payload is invalid", Error: err.Error()})
		return
	}

	args := &db.CreateTransactionParams{
		Date:  payload.Date,
		Price: payload.Price,
	}

	// TODO change balance in savapage

	// start savapage request

	config, err := util.LoadConfig("../.")

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to Reach SavaPage", Error: err.Error()})
		return
	}

	url := config.SavaPageUrl + fmt.Sprintf("financial/account/balance?type=USER&name=%s&amount=-%.2f&adjust=true&details=rupay_transaction", UserName, payload.Price)

	print(url)

	req, err := http.NewRequest("POST", url, nil)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to Reach SavaPage", Error: err.Error()})
		return
	}
	println(config.SavaPageAdmin)
	println(config.SavaPagePassword)

	req.SetBasicAuth(config.SavaPageAdmin, config.SavaPagePassword)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to Reach SavaPage", Error: err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		println(resp.StatusCode)
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to Reach SavaPage", Error: "Error"})
		return
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
