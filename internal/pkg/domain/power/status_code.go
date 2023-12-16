package power

import (
	"errors"
	"strings"
)

var (
	ErrInvalidStatus = errors.New("invalid machine status")
)

type StatusCode int64

const (
	StatusCodeRunning StatusCode = iota
	StatusCodeStopped
	StatusCodeStopping
	StatusCodeStarting
	StatusCodeRebooting
	StatusCodePending
	StatusCodeUnknown
)

func StatusCodeFromString(v string) (StatusCode, error) {
	switch strings.ToLower(v) {
	case "running":
		return StatusCodeRunning, nil
	case "stopped":
		return StatusCodeStopped, nil
	case "stopping":
		return StatusCodeStopping, nil
	case "starting":
		return StatusCodeStarting, nil
	case "restarting", "rebooting":
		return StatusCodeRebooting, nil
	case "pending":
		return StatusCodePending, nil
	default:
		return StatusCodeUnknown, ErrInvalidStatus
	}
}

func (m StatusCode) String() string {
	switch m {
	case StatusCodeRunning:
		return "Running"
	case StatusCodeStopped:
		return "Stopped"
	case StatusCodeStopping:
		return "Stopping"
	case StatusCodeStarting:
		return "Starting"
	case StatusCodeRebooting:
		return "Rebooting"
	case StatusCodePending:
		return "Pending"
	default:
		return "Unknown"
	}
}
