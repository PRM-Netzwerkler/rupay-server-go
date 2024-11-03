package schemas

import (
	"time"

	"github.com/guregu/null/v5"
)

type CreateTransaction struct {
	Date  time.Time `json:"date" binding:"required" example:"2024-01-24T00:00:00Z"`
	Price float64   `json:"price" binding:"required"`
}

type UpdateTransaction struct {
	Date  null.Time  `json:"date" example:"2024-01-24T00:00:00Z"`
	Price null.Float `json:"price"`
}
