package power_change_topic

import (
	"encoding/json"
	"fmt"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/apps/local/internal/core/domain/power"
)

type PowerChangeRequestMessageJson struct {
	Identifier string  `json:"identifier"`
	Target     string  `json:"target"`
	Key        *string `json:"key,omitempty"`
	ChangeType string  `json:"changeType"`
	Timeout    int64   `json:"timeout"`
}

func PowerChangeRequestMessageJsonFromDomain(msg *power.PowerChangeRequestMessage) *PowerChangeRequestMessageJson {
	return &PowerChangeRequestMessageJson{
		Identifier: msg.Identifier.String(),
		Target:     msg.Target,
		Key:        msg.Key,
		ChangeType: msg.ChangeType.String(),
		Timeout:    msg.Timeout,
	}
}

func (m *PowerChangeRequestMessageJson) ToDomain() (*power.PowerChangeRequestMessage, error) {
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
		if m.Key == nil {
			return nil, fmt.Errorf("key is nil, is required for validation")
		}
		return power.NewRebootMessage(id, m.Target, *m.Key), nil
	case power.ChangeTypePowerOff:
		if m.Key == nil {
			return nil, fmt.Errorf("key is nil, is required for validation")
		}
		return power.NewPowerOffMessage(id, m.Target, *m.Key), nil
	case power.ChangeTypeWakeOnLan:
		return power.NewWakeOnLanMessage(id, m.Target), nil
	default:
		return nil, power.ErrInvalidChangeType
	}
}

func (m *PowerChangeRequestMessageJson) ToJson() ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return b, nil
}
