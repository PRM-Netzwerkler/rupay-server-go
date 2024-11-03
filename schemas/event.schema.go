package schemas

import (
	"time"

	"github.com/guregu/null/v5"
)

type CreateEvent struct {
	Name     string      `json:"name" binding:"required"`
	Desc     null.String `json:"desc"`
	FromDate time.Time   `json:"from_date" binding:"required"`
	ToDate   time.Time   `json:"to_date" binding:"required"`
}

type UpdateEvent struct {
	Name     null.String `json:"name"`
	Desc     null.String `json:"desc"`
	FromDate null.Time   `json:"from_date"`
	ToDate   null.Time   `json:"to_date"`
}
