package power

import (
	"encoding/json"
	"time"

	"github.com/awlsring/camp/internal/app/local/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/domain/power"
)

type StateChangeMessage struct {
	Identifier machine.Identifier `json:"identifier"`
	Was        power.StatusCode   `json:"was"`
	Now        power.StatusCode   `json:"now"`
	Planned    bool               `json:"planned"`
	Time       time.Time          `json:"time"`
}

func NewStateChangeMessage(identifier machine.Identifier, was power.StatusCode, now power.StatusCode, planned bool) *StateChangeMessage {
	return &StateChangeMessage{
		Identifier: identifier,
		Was:        was,
		Now:        now,
		Planned:    planned,
		Time:       time.Now().UTC(),
	}
}

func (m *StateChangeMessage) AsJsonModel() *StateChangeMessageJson {
	return &StateChangeMessageJson{
		Identifier: m.Identifier.String(),
		Was:        m.Was.String(),
		Now:        m.Now.String(),
		Planned:    m.Planned,
		Time:       m.Time,
	}
}

func (m *StateChangeMessage) ToJson() ([]byte, error) {
	jsonMsg := m.AsJsonModel()
	return json.Marshal(jsonMsg)
}

type StateChangeMessageJson struct {
	Identifier string    `json:"identifier"`
	Was        string    `json:"was"`
	Now        string    `json:"now"`
	Planned    bool      `json:"planned"`
	Time       time.Time `json:"time"`
}

func (m *StateChangeMessageJson) ToDomain() (*StateChangeMessage, error) {
	identifier, err := machine.IdentifierFromString(m.Identifier)
	if err != nil {
		return nil, err
	}

	was, err := power.StatusCodeFromString(m.Was)
	if err != nil {
		return nil, err
	}

	now, err := power.StatusCodeFromString(m.Now)
	if err != nil {
		return nil, err
	}

	return &StateChangeMessage{
		Identifier: identifier,
		Was:        was,
		Now:        now,
		Planned:    m.Planned,
		Time:       m.Time,
	}, nil
}
