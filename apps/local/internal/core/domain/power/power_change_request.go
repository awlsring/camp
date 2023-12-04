package power

import "github.com/awlsring/camp/apps/local/internal/core/domain/machine"

const TimeoutFiveMinutes = 300

type PowerChangeRequestMessage struct {
	Identifier machine.Identifier
	Target     string
	Key        *string
	ChangeType ChangeType
	Timeout    int64
}

func NewRebootMessage(identifier machine.Identifier, endpoint string, key string) *PowerChangeRequestMessage {
	return &PowerChangeRequestMessage{
		Identifier: identifier,
		Target:     endpoint,
		Key:        &key,
		ChangeType: ChangeTypeReboot,
		Timeout:    TimeoutFiveMinutes,
	}
}

func NewPowerOffMessage(identifier machine.Identifier, endpoint string, key string) *PowerChangeRequestMessage {
	return &PowerChangeRequestMessage{
		Identifier: identifier,
		Target:     endpoint,
		Key:        &key,
		ChangeType: ChangeTypePowerOff,
		Timeout:    TimeoutFiveMinutes,
	}
}

func NewWakeOnLanMessage(identifier machine.Identifier, mac string) *PowerChangeRequestMessage {
	return &PowerChangeRequestMessage{
		Identifier: identifier,
		Target:     mac,
		ChangeType: ChangeTypeWakeOnLan,
		Timeout:    TimeoutFiveMinutes,
	}
}
