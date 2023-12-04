package power

import "github.com/awlsring/camp/apps/local/internal/core/domain/machine"

type ValidatePowerChangeRequestMessage struct {
	Identifier machine.Identifier
	Target     string
	Key        *string
	Attempt    int
	ChangeType ChangeType
}

func NewValidateRebootMessage(identifier machine.Identifier, endpoint string, key string) *ValidatePowerChangeRequestMessage {
	return &ValidatePowerChangeRequestMessage{
		Identifier: identifier,
		Target:     endpoint,
		Key:        &key,
		ChangeType: ChangeTypeReboot,
	}
}

func NewValidatePowerOffMessage(identifier machine.Identifier, endpoint string, key string) *ValidatePowerChangeRequestMessage {
	return &ValidatePowerChangeRequestMessage{
		Identifier: identifier,
		Target:     endpoint,
		Key:        &key,
		ChangeType: ChangeTypePowerOff,
	}
}

func NewValidatePowerOnMessage(identifier machine.Identifier, endpoint string, key string) *ValidatePowerChangeRequestMessage {
	return &ValidatePowerChangeRequestMessage{
		Identifier: identifier,
		Target:     endpoint,
		Key:        &key,
		ChangeType: ChangeTypeWakeOnLan,
	}
}
