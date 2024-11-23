package schemas

import (
	"github.com/guregu/null/v5"
)

type CreateUser struct {
	Name string `json:"name" binding:"required"`
	Code string `json:"code" binding:"required"`
}

type UpdateUser struct {
	Code null.String `json:"code"`
}
