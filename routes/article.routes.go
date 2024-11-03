package routes

import (
	"github.com/KevinGruber2001/rupay-bar-backend/controllers"
	"github.com/gin-gonic/gin"
)

type ArticleRoutes struct {
    ArticleController controllers.ArticleController
}

func NewRouteArticle(ArticleController controllers.ArticleController) ArticleRoutes {
    return ArticleRoutes{ArticleController}
}

func (cr *ArticleRoutes) ArticleRoute(rg *gin.RouterGroup) {

    router := rg.Group("article")
    router.POST("/", cr.ArticleController.CreateArticle)
    router.GET("/", cr.ArticleController.GetAllArticles)
    router.PATCH("/:articleId", cr.ArticleController.UpdateArticle)
    router.GET("/:articleId", cr.ArticleController.GetArticleById)
    router.DELETE("/:articleId", cr.ArticleController.DeleteArticleById)
}