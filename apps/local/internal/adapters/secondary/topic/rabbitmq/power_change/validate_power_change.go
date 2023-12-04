package power_change_topic

import (
	"encoding/json"
	"fmt"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/apps/local/internal/core/domain/power"
)

type ValidatePowerChangeRequestMessageJson struct {
	Identifier string  `json:"identifier"`
	Target     string  `json:"target"`
	Key        *string `json:"key,omitempty"`
	Attempt    int     `json:"attempt"`
	ChangeType string  `json:"changeType"`
}

func ValidatePowerChangeRequestMessageJsonFromDomain(msg *power.ValidatePowerChangeRequestMessage) *ValidatePowerChangeRequestMessageJson {
	return &ValidatePowerChangeRequestMessageJson{
		Identifier: msg.Identifier.String(),
		Target:     msg.Target,
		Key:        msg.Key,
		Attempt:    msg.Attempt,
		ChangeType: msg.ChangeType.String(),
	}
}

func (m *ValidatePowerChangeRequestMessageJson) ToDomain() (*power.ValidatePowerChangeRequestMessage, error) {
	if m.Key == nil {
		return nil, fmt.Errorf("key is nil, is required for validation")
	}

	changeType, err := power.ChangeTypeFromString(m.ChangeType)
	if err != nil {
		return nil, err
	}

	id, err := machine.IdentifierFromString(m.Identifier)
	if err != nil {
		return nil, err
	}

	switch changeType {
	case power.ChangeTypeReboot:
		return power.NewValidateRebootMessage(id, m.Target, *m.Key), nil
	case power.ChangeTypePowerOff:
		return power.NewValidatePowerOffMessage(id, m.Target, *m.Key), nil
	case power.ChangeTypeWakeOnLan:
		return power.NewValidatePowerOnMessage(id, m.Target, *m.Key), nil
	default:
		return nil, power.ErrInvalidChangeType
	}
}

func (m *ValidatePowerChangeRequestMessageJson) ToJson() ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return b, nil
}
