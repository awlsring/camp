package machine

import "fmt"

type PowerCapabilityType int

var ErrInvalidPowerCapabilityType = fmt.Errorf("invalid power capability type")

const (
	PowerCapabilityTypeWakeOnLan PowerCapabilityType = iota
	PowerCapabilityTypePowerOff
	PowerCapabilityTypeReboot
	PowerCapabilityTypeUnknown
)

func (c PowerCapabilityType) String() string {
	switch c {
	case PowerCapabilityTypeWakeOnLan:
		return "WakeOnLan"
	case PowerCapabilityTypePowerOff:
		return "PowerOff"
	case PowerCapabilityTypeReboot:
		return "Reboot"
	default:
		return "Unknown"
	}
}

func PowerCapabilityTypeFromString(s string) (PowerCapabilityType, error) {
	switch s {
	case "WakeOnLan":
		return PowerCapabilityTypeWakeOnLan, nil
	case "PowerOff":
		return PowerCapabilityTypePowerOff, nil
	case "Reboot":
		return PowerCapabilityTypeReboot, nil
	default:
		return PowerCapabilityTypeUnknown, ErrInvalidPowerCapabilityType
	}
}
