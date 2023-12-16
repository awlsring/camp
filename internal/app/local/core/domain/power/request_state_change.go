package power

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/awlsring/camp/internal/app/local/core/domain/machine"
)

const (
	StartingTimeout  = 600
	StoppingTimeout  = 600
	RebootingTimeout = StartingTimeout + StoppingTimeout
)

var ErrKeyRequired = fmt.Errorf("key is nil, is required for validation")

type RequestStateChangeMessage struct {
	Identifier machine.Identifier
	ChangeType ChangeType
	Time       time.Time
}

func NewRebootMessage(identifier machine.Identifier, endpoint string, key string) *RequestStateChangeMessage {
	now := time.Now().UTC()
	return &RequestStateChangeMessage{
		Identifier: identifier,
		ChangeType: ChangeTypeReboot,
		Time:       now,
	}
}

func NewPowerOffMessage(identifier machine.Identifier, endpoint string, key string) *RequestStateChangeMessage {
	now := time.Now().UTC()
	return &RequestStateChangeMessage{
		Identifier: identifier,
		ChangeType: ChangeTypePowerOff,
		Time:       now,
	}
}

func NewWakeOnLanMessage(identifier machine.Identifier, mac string) *RequestStateChangeMessage {
	now := time.Now().UTC()
	return &RequestStateChangeMessage{
		Identifier: identifier,
		ChangeType: ChangeTypeWakeOnLan,
		Time:       now,
	}
}

func (m *RequestStateChangeMessage) AsJsonModel() *RequestStateChangeMessageJson {
	return &RequestStateChangeMessageJson{
		Identifier: m.Identifier.String(),
		ChangeType: m.ChangeType.String(),
		Time:       m.Time,
	}
}

func (m *RequestStateChangeMessage) ToJson() ([]byte, error) {
	jsonMsg := m.AsJsonModel()
	return json.Marshal(jsonMsg)
}

type RequestStateChangeMessageJson struct {
	Identifier string    `json:"identifier"`
	ChangeType string    `json:"changeType"`
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

	return &RequestStateChangeMessage{
		Identifier: id,
		ChangeType: changeType,
		Time:       m.Time,
	}, nil
}
