package machine

import (
	"errors"
	"strings"

	camperror "github.com/awlsring/camp/internal/pkg/errors"
)

var (
	ErrInvalidStatus = errors.New("invalid machine status")
)

type MachineStatus int64

const (
	MachineStatusRunning MachineStatus = iota
	MachineStatusStopped
	MachineStatusStopping
	MachineStatusStarting
	MachineStatusRestarting
	MachineStatusUnknown
)

func MachineStatusFromString(v string) (MachineStatus, error) {
	switch strings.ToLower(v) {
	case "running":
		return MachineStatusRunning, nil
	case "stopped":
		return MachineStatusStopped, nil
	case "stopping":
		return MachineStatusStopping, nil
	case "starting":
		return MachineStatusStarting, nil
	case "restarting":
		return MachineStatusRestarting, nil
	default:
		return MachineStatusUnknown, camperror.New(camperror.ErrValidation, ErrInvalidStatus)
	}
}

func (m MachineStatus) String() string {
	switch m {
	case MachineStatusRunning:
		return "Running"
	case MachineStatusStopped:
		return "Stopped"
	case MachineStatusStopping:
		return "Stopping"
	case MachineStatusStarting:
		return "Starting"
	case MachineStatusRestarting:
		return "Restarting"
	default:
		return "Unknown"
	}
}