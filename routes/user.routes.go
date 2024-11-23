package routes

import (
	"github.com/KevinGruber2001/rupay-bar-backend/controllers"
	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	UserController controllers.UserController
}

func NewRouteUser(UserController controllers.UserController) UserRoutes {
	return UserRoutes{UserController}
}

func (cr *UserRoutes) UserRoute(rg *gin.RouterGroup) {

	router := rg.Group("user")
	router.POST("/", cr.UserController.CreateUser)
	router.GET("/", cr.UserController.GetAllUsers)
	router.PATCH("/:name", cr.UserController.UpdateUser)
	router.GET("/:name", cr.UserController.GetUserByUsername)
	router.DELETE("/:name", cr.UserController.DeleteUserByUsername)
}
