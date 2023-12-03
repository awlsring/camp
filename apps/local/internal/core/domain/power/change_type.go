package power

import (
	"errors"
	"strings"
)

type ChangeType int

var ErrInvalidChangeType = errors.New("invalid change type")

const (
	ChangeTypePowerOff ChangeType = iota
	ChangeTypeReboot
	ChangeTypeWakeOnLan
)

func (t ChangeType) String() string {
	switch t {
	case ChangeTypePowerOff:
		return "PowerOff"
	case ChangeTypeReboot:
		return "Reboot"
	case ChangeTypeWakeOnLan:
		return "WakeOnLan"
	default:
		return "Unknown"
	}
}

func ChangeTypeFromString(s string) (ChangeType, error) {
	switch strings.ToLower(s) {
	case "poweroff":
		return ChangeTypePowerOff, nil
	case "reboot":
		return ChangeTypeReboot, nil
	case "wakeonlan":
		return ChangeTypeWakeOnLan, nil
	default:
		return 0, ErrInvalidChangeType
	}
}
