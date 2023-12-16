package machine

import "github.com/awlsring/camp/internal/pkg/domain/network"

type PowerCapabilities struct {
	WakeOnLan PowerCapabilityWakeOnLan
	PowerOff  PowerCapabilityPowerOff
	Reboot    PowerCapabilityReboot
}

func NewPowerCapabilities(wakeOnLan *PowerCapabilityWakeOnLan, powerOff *PowerCapabilityPowerOff, reboot *PowerCapabilityReboot) *PowerCapabilities {
	rebootCap := PowerCapabilityReboot{
		Enabled: false,
	}
	if reboot != nil {
		rebootCap = *reboot
	}

	powerOffCap := PowerCapabilityPowerOff{
		Enabled: false,
	}
	if powerOff != nil {
		powerOffCap = *powerOff
	}

	wakeOnLanCap := PowerCapabilityWakeOnLan{
		Enabled: false,
	}
	if wakeOnLan != nil {
		wakeOnLanCap = *wakeOnLan
	}

	return &PowerCapabilities{
		WakeOnLan: wakeOnLanCap,
		PowerOff:  powerOffCap,
		Reboot:    rebootCap,
	}
}

type PowerCapabilityReboot struct {
	Enabled bool
}

func NewPowerCapabilityReboot(enabled bool) *PowerCapabilityReboot {
	return &PowerCapabilityReboot{
		Enabled: enabled,
	}
}

type PowerCapabilityPowerOff struct {
	Enabled bool
}

func NewPowerCapabilityPowerOff(enabled bool) *PowerCapabilityPowerOff {
	return &PowerCapabilityPowerOff{
		Enabled: enabled,
	}
}

type PowerCapabilityWakeOnLan struct {
	Enabled    bool
	MacAddress *network.MacAddress
}

func NewPowerCapabilityWakeOnLan(enabled bool, macAddress *network.MacAddress) *PowerCapabilityWakeOnLan {
	return &PowerCapabilityWakeOnLan{
		Enabled:    enabled,
		MacAddress: macAddress,
	}
}
