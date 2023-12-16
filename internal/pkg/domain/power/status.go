package power

import (
	"time"
)

type Status struct {
	Status    StatusCode
	UpdatedAt time.Time
}

func NewStatus(status StatusCode, updatedAt time.Time) *Status {
	return &Status{
		Status:    status,
		UpdatedAt: updatedAt,
	}
}
