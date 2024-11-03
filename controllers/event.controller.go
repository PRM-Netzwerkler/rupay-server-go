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

type EventController struct {
	db  *db.Queries
	ctx context.Context
}

func NewEventController(db *db.Queries, ctx context.Context) *EventController {
	return &EventController{db, ctx}
}

// @Summary Create a new event
// @Description Create a new event with the provided details
// @Tags Events
// @Accept json
// @Produce json
// @Param payload body schemas.CreateEvent true "CreateEvent payload"
// @Success 200 {object} db.Event "Event data"
// @Failure 400 {object} e.ErrorResponse "Invalid Payload"
// @Router /event [post]
func (cc *EventController) CreateEvent(ctx *gin.Context) {
	var payload *schemas.CreateEvent

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, e.ErrorResponse{Code: e.InvalidPayload, Message: "Invalid Payload", Error: err.Error()})
		return
	}

	args := &db.CreateEventParams{
		Name:     payload.Name,
		Desc:     payload.Desc,
		FromDate: payload.FromDate,
		ToDate:   payload.ToDate,
	}

	Event, err := cc.db.CreateEvent(ctx, *args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to create Event", Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, Event)
}

// @Summary Update an event
// @Description Update an event by the provided id and details
// @Tags Events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} db.Event "Event data"
// @Failure 400 {object} e.ErrorResponse "Invalid Payload"
// @Router /event/{id} [patch]
func (cc *EventController) UpdateEvent(ctx *gin.Context) {
	var payload *schemas.UpdateEvent
	EventId := ctx.Param("eventId")

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, e.ErrorResponse{Code: e.InvalidPayload, Message: "Invalid Payload", Error: err.Error()})
		return
	}

	args := &db.UpdateEventParams{
		Uuid:     uuid.MustParse(EventId),
		Name:     payload.Name,
		Desc:     payload.Desc,
		FromDate: payload.FromDate,
		ToDate:   payload.ToDate,
	}

	Event, err := cc.db.UpdateEvent(ctx, *args)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, e.ErrorResponse{Code: e.NotFound, Message: "Event not found", Error: err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to update Event", Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, Event)
}

// @Summary Retrieve an event by ID
// @Description Retrieve an event by the provided ID
// @Tags Events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {object} db.Event "Event data"
// @Failure 404 {object} e.ErrorResponse "Event not found"
// @Router /event/{id} [get]
func (cc *EventController) GetEventById(ctx *gin.Context) {
	EventId := ctx.Param("eventId")

	Event, err := cc.db.GetEventById(ctx, uuid.MustParse(EventId))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, e.ErrorResponse{Code: e.NotFound, Message: "Event not found", Error: err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to retrieve Event", Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, Event)
}

// @Summary Retrieve all events
// @Description Retrieve a list of all events
// @Tags Events
// @Accept json
// @Produce json
// @Success 200 {array} db.Event "List of events"
// @Failure 500 {object} e.ErrorResponse "Internal Server Error"
// @Router /events [get]
func (cc *EventController) GetAllEvents(ctx *gin.Context) {
	Events, err := cc.db.GetEvents(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to retrieve Events", Error: err.Error()})
		return
	}

	if Events == nil {
		Events = []db.Event{}
	}

	ctx.JSON(http.StatusOK, Events)
}

// @Summary Delete an event by ID
// @Description Delete an event by the provided ID
// @Tags Events
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Success 204 "Event deleted successfully"
// @Failure 404 {object} e.ErrorResponse "Event not found"
// @Router /event/{id} [delete]
func (cc *EventController) DeleteEventById(ctx *gin.Context) {
	EventId := ctx.Param("eventId")

	_, err := cc.db.GetEventById(ctx, uuid.MustParse(EventId))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, e.ErrorResponse{Code: e.NotFound, Message: "Event not found", Error: err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to retrieve Event", Error: err.Error()})
		return
	}

	err = cc.db.DeleteEvent(ctx, uuid.MustParse(EventId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to delete Event", Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
