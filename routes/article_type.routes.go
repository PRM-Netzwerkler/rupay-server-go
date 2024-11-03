package routes

import (
	"github.com/KevinGruber2001/rupay-bar-backend/controllers"
	"github.com/gin-gonic/gin"
)

type ArticleTypeRoutes struct {
	ArticleTypeController controllers.ArticleTypeController
}

func NewRouteArticleType(ArticleTypeController controllers.ArticleTypeController) ArticleTypeRoutes {
	return ArticleTypeRoutes{ArticleTypeController}
}

func (cr *ArticleTypeRoutes) ArticleTypeRoute(rg *gin.RouterGroup) {

	router := rg.Group("article-type")
	router.POST("/", cr.ArticleTypeController.CreateArticleType)
	router.GET("/", cr.ArticleTypeController.GetAllArticleTypes)
	router.GET("/article", cr.ArticleTypeController.GetAllArticleTypesWithArticles)
	router.PATCH("/:articleTypeId", cr.ArticleTypeController.UpdateArticleType)
	router.GET("/:articleTypeId", cr.ArticleTypeController.GetArticleTypeById)
	router.DELETE("/:articleTypeId", cr.ArticleTypeController.DeleteArticleTypeById)
}
