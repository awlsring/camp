package power

import (
	"encoding/json"
	"time"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
)

type StateValidationMessage struct {
	Identifier    machine.Identifier
	ReportedState machine.MachineStatus
	Target        string
	Key           string
	Time          time.Time
}

func NewStateValidationMessage(identifier machine.Identifier, reportedState machine.MachineStatus, target string, key string) *StateValidationMessage {
	now := time.Now().UTC().UTC()
	return &StateValidationMessage{
		Identifier:    identifier,
		ReportedState: reportedState,
		Target:        target,
		Key:           key,
		Time:          now,
	}
}

func (m *StateValidationMessage) AsJsonModel() *StateValidationMessageJson {
	return &StateValidationMessageJson{
		Identifier:    m.Identifier.String(),
		ReportedState: m.ReportedState.String(),
		Target:        m.Target,
		Key:           m.Key,
		Time:          m.Time,
	}
}

func (m *StateValidationMessage) AsJson() ([]byte, error) {
	return json.Marshal(m.AsJsonModel())
}

type StateValidationMessageJson struct {
	Identifier    string    `json:"identifier"`
	ReportedState string    `json:"reportedState"`
	Target        string    `json:"target"`
	Key           string    `json:"key,omitempty"`
	Time          time.Time `json:"time"`
}

func (m *StateValidationMessageJson) ToDomain() (*StateValidationMessage, error) {
	id, err := machine.IdentifierFromString(m.Identifier)
	if err != nil {
		return nil, err
	}

	reportedState, err := machine.MachineStatusFromString(m.ReportedState)
	if err != nil {
		return nil, err
	}

	return &StateValidationMessage{
		Identifier:    id,
		ReportedState: reportedState,
		Target:        m.Target,
		Key:           m.Key,
		Time:          m.Time,
	}, nil
}
