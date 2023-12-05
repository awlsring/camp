package power

import (
	"encoding/json"
	"time"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
)

type StateChangeMessage struct {
	Identifier machine.Identifier    `json:"identifier"`
	Was        machine.MachineStatus `json:"was"`
	Now        machine.MachineStatus `json:"now"`
	Planned    bool                  `json:"planned"`
	Time       time.Time             `json:"time"`
}

func NewStateChangeMessage(identifier machine.Identifier, was machine.MachineStatus, now machine.MachineStatus, planned bool) *StateChangeMessage {
	return &StateChangeMessage{
		Identifier: identifier,
		Was:        was,
		Now:        now,
		Planned:    planned,
		Time:       time.Now(),
	}
}

func (m *StateChangeMessage) AsJsonModel() *StateChangeMessageJson {
	return &StateChangeMessageJson{
		Identifier: m.Identifier.String(),
		Was:        m.Was.String(),
		Now:        m.Now.String(),
		Planned:    m.Planned,
		Time:       m.Time.Unix(),
	}
}

func (m *StateChangeMessage) ToJson() ([]byte, error) {
	jsonMsg := m.AsJsonModel()
	return json.Marshal(jsonMsg)
}

type StateChangeMessageJson struct {
	Identifier string `json:"identifier"`
	Was        string `json:"was"`
	Now        string `json:"now"`
	Planned    bool   `json:"planned"`
	Time       int64  `json:"time"`
}

func (m *StateChangeMessageJson) ToDomain() (*StateChangeMessage, error) {
	identifier, err := machine.IdentifierFromString(m.Identifier)
	if err != nil {
		return nil, err
	}

	was, err := machine.MachineStatusFromString(m.Was)
	if err != nil {
		return nil, err
	}

	now, err := machine.MachineStatusFromString(m.Now)
	if err != nil {
		return nil, err
	}

	return &StateChangeMessage{
		Identifier: identifier,
		Was:        was,
		Now:        now,
		Planned:    m.Planned,
		Time:       time.Unix(m.Time, 0),
	}, nil
}
