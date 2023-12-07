package power

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
)

const (
	StartingTimeout  = 600
	StoppingTimeout  = 600
	RebootingTimeout = StartingTimeout + StoppingTimeout
)

var ErrKeyRequired = fmt.Errorf("key is nil, is required for validation")

type RequestStateChangeMessage struct {
	Identifier machine.Identifier
	Target     string
	Key        *string
	ChangeType ChangeType
	Timeout    int64
	Time       time.Time
}

func NewRebootMessage(identifier machine.Identifier, endpoint string, key string) *RequestStateChangeMessage {
	return &RequestStateChangeMessage{
		Identifier: identifier,
		Target:     endpoint,
		Key:        &key,
		ChangeType: ChangeTypeReboot,
		Timeout:    RebootingTimeout,
		Time:       time.Now().UTC().UTC(),
	}
}

func NewPowerOffMessage(identifier machine.Identifier, endpoint string, key string) *RequestStateChangeMessage {
	return &RequestStateChangeMessage{
		Identifier: identifier,
		Target:     endpoint,
		Key:        &key,
		ChangeType: ChangeTypePowerOff,
		Timeout:    StoppingTimeout,
		Time:       time.Now().UTC().UTC(),
	}
}

func NewWakeOnLanMessage(identifier machine.Identifier, mac string) *RequestStateChangeMessage {
	return &RequestStateChangeMessage{
		Identifier: identifier,
		Target:     mac,
		ChangeType: ChangeTypeWakeOnLan,
		Timeout:    StartingTimeout,
		Time:       time.Now().UTC().UTC(),
	}
}

func (m *RequestStateChangeMessage) AsJsonModel() *RequestStateChangeMessageJson {
	return &RequestStateChangeMessageJson{
		Identifier: m.Identifier.String(),
		Target:     m.Target,
		Key:        m.Key,
		ChangeType: m.ChangeType.String(),
		Timeout:    m.Timeout,
		Time:       m.Time,
	}
}

func (m *RequestStateChangeMessage) ToJson() ([]byte, error) {
	jsonMsg := m.AsJsonModel()
	return json.Marshal(jsonMsg)
}

type RequestStateChangeMessageJson struct {
	Identifier string    `json:"identifier"`
	Target     string    `json:"target"`
	Key        *string   `json:"key,omitempty"`
	ChangeType string    `json:"changeType"`
	Timeout    int64     `json:"timeout"`
	Time       time.Time `json:"time"`
}

func (m *RequestStateChangeMessageJson) ToDomain() (*RequestStateChangeMessage, error) {
	changeType, err := ChangeTypeFromString(m.ChangeType)
	if err != nil {
		return nil, err
	}

	id, err := machine.IdentifierFromString(m.Identifier)
	if err != nil {
		return nil, err
	}

	switch changeType {
	case ChangeTypeReboot:
		if m.Key == nil {
			return nil, ErrKeyRequired
		}
		return NewRebootMessage(id, m.Target, *m.Key), nil
	case ChangeTypePowerOff:
		if m.Key == nil {
			return nil, ErrKeyRequired
		}
		return NewPowerOffMessage(id, m.Target, *m.Key), nil
	case ChangeTypeWakeOnLan:
		return NewWakeOnLanMessage(id, m.Target), nil
	default:
		return nil, ErrInvalidChangeType
	}
}
