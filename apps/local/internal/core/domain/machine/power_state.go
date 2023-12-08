package machine

import (
	"time"
)

type PowerState struct {
	State     MachineStatus
	UpdatedAt time.Time
}

func NewPowerState(state MachineStatus, updatedAt time.Time) *PowerState {
	return &PowerState{
		State:     state,
		UpdatedAt: updatedAt,
	}
}
