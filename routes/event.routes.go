package routes

import (
	"github.com/KevinGruber2001/rupay-bar-backend/controllers"
	"github.com/gin-gonic/gin"
)

type EventRoutes struct {
    EventController controllers.EventController
}

func NewRouteEvent(EventController controllers.EventController) EventRoutes {
    return EventRoutes{EventController}
}

func (cr *EventRoutes) EventRoute(rg *gin.RouterGroup) {

    router := rg.Group("event")
    router.POST("/", cr.EventController.CreateEvent)
    router.GET("/", cr.EventController.GetAllEvents)
    router.PATCH("/:eventId", cr.EventController.UpdateEvent)
    router.GET("/:eventId", cr.EventController.GetEventById)
    router.DELETE("/:eventId", cr.EventController.DeleteEventById)
}