package routes

import (
	"github.com/KevinGruber2001/rupay-bar-backend/controllers"
	"github.com/gin-gonic/gin"
)

type TransactionRoutes struct {
    TransactionController controllers.TransactionController
}

func NewRouteTransaction(TransactionController controllers.TransactionController) TransactionRoutes {
    return TransactionRoutes{TransactionController}
}

func (cr *TransactionRoutes) TransactionRoute(rg *gin.RouterGroup) {

    router := rg.Group("transaction")
    router.POST("/", cr.TransactionController.CreateTransaction)
    router.GET("/", cr.TransactionController.GetAllTransactions)
    router.PATCH("/:transactionId", cr.TransactionController.UpdateTransaction)
    router.GET("/:transactionId", cr.TransactionController.GetTransactionById)
    router.DELETE("/:transactionId", cr.TransactionController.DeleteTransactionById)
}