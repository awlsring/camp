package power

import "github.com/awlsring/camp/apps/local/internal/core/domain/machine"

type PowerChangeRequestMessage struct {
	Identifier machine.Identifier
	Target     string
	Key        *string
	ChangeType ChangeType
}

func NewRebootMessage(identifier machine.Identifier, endpoint string, key string) *PowerChangeRequestMessage {
	return &PowerChangeRequestMessage{
		Identifier: identifier,
		Target:     endpoint,
		Key:        &key,
		ChangeType: ChangeTypeReboot,
	}
}

func NewPowerOffMessage(identifier machine.Identifier, endpoint string, key string) *PowerChangeRequestMessage {
	return &PowerChangeRequestMessage{
		Identifier: identifier,
		Target:     endpoint,
		Key:        &key,
		ChangeType: ChangeTypePowerOff,
	}
}

func NewWakeOnLanMessage(identifier machine.Identifier, mac string) *PowerChangeRequestMessage {
	return &PowerChangeRequestMessage{
		Identifier: identifier,
		Target:     mac,
		ChangeType: ChangeTypeWakeOnLan,
	}
}
