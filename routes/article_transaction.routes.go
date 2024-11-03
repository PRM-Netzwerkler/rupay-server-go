package routes

import (
	"github.com/KevinGruber2001/rupay-bar-backend/controllers"
	"github.com/gin-gonic/gin"
)

type ArticleTransactionRoutes struct {
	ArticleTransactionController controllers.ArticleTransactionController
}

func NewRouteArticleTransaction(ArticleTransactionController controllers.ArticleTransactionController) ArticleTransactionRoutes {
	return ArticleTransactionRoutes{ArticleTransactionController}
}

func (cr *ArticleTransactionRoutes) ArticleTransactionRoute(rg *gin.RouterGroup) {

	router := rg.Group("article-transaction")
	router.POST("/", cr.ArticleTransactionController.CreateArticleTransaction)
	router.GET("/", cr.ArticleTransactionController.GetAllArticleTransactions)
	router.GET("/grouped-by-article", cr.ArticleTransactionController.GetAllArticleTransactionsGroupedByArticle)
	router.PATCH("/:articleTransactionId", cr.ArticleTransactionController.UpdateArticleTransaction)
	router.GET("/:articleTransactionId", cr.ArticleTransactionController.GetArticleTransactionById)
	router.DELETE("/:articleTransactionId", cr.ArticleTransactionController.DeleteArticleTransactionById)
}
