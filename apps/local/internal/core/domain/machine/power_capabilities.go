package machine

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
	MacAddress *MacAddress
}

func NewPowerCapabilityWakeOnLan(enabled bool, macAddress *MacAddress) *PowerCapabilityWakeOnLan {
	return &PowerCapabilityWakeOnLan{
		Enabled:    enabled,
		MacAddress: macAddress,
	}
}
