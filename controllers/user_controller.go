package controllers

import (
	"context"
	"database/sql"
	"net/http"

	db "github.com/KevinGruber2001/rupay-bar-backend/db/sqlc"
	e "github.com/KevinGruber2001/rupay-bar-backend/errors"
	"github.com/KevinGruber2001/rupay-bar-backend/schemas"
	"github.com/gin-gonic/gin"
	"github.com/guregu/null/v5"
)

// UserController handles user-related requests
type UserController struct {
	db  *db.Queries
	ctx context.Context
}

// NewUserController creates a new UserController
func NewUserController(db *db.Queries, ctx context.Context) *UserController {
	return &UserController{db, ctx}
}

func (cc *UserController) CreateUser(ctx *gin.Context) {
	var payload *schemas.CreateUser

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, e.ErrorResponse{Code: e.InvalidPayload, Message: "Invalid Payload", Error: err.Error()})
		return
	}

	args := &db.CreateUserParams{
		Name: payload.Name,
		Code: payload.Code,
	}

	user, err := cc.db.CreateUser(ctx, *args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to create User", Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (cc *UserController) UpdateUser(ctx *gin.Context) {
	var payload *schemas.UpdateUser
	name := ctx.Param("username")

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, e.ErrorResponse{Code: e.InvalidPayload, Message: "Invalid Payload", Error: err.Error()})
		return
	}

	args := &db.UpdateUserParams{
		Name: null.StringFrom(name),
		Code: payload.Code,
	}

	user, err := cc.db.UpdateUser(ctx, *args)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, e.ErrorResponse{Code: e.NotFound, Message: "User not found", Error: err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to update User", Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (cc *UserController) GetUserByUsername(ctx *gin.Context) {
	name := ctx.Param("name")

	user, err := cc.db.GetUserById(ctx, name)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, e.ErrorResponse{Code: e.NotFound, Message: "User not found", Error: err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to retrieve User", Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (cc *UserController) GetAllUsers(ctx *gin.Context) {
	users, err := cc.db.GetUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to retrieve Users", Error: err.Error()})
		return
	}

	if users == nil {
		users = []db.Resident{}
	}

	ctx.JSON(http.StatusOK, users)
}

func (cc *UserController) DeleteUserByUsername(ctx *gin.Context) {
	name := ctx.Param("name")

	_, err := cc.db.GetUserById(ctx, name)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, e.ErrorResponse{Code: e.NotFound, Message: "User not found", Error: err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to retrieve User", Error: err.Error()})
		return
	}

	err = cc.db.DeleteUser(ctx, name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, e.ErrorResponse{Code: e.InternalServerError, Message: "Failed to delete User", Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
