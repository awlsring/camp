package power

import (
	"encoding/json"
	"time"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
)

type StateValidationMessage struct {
	Identifier    machine.Identifier
	ReportedState machine.MachineStatus
	Time          time.Time
}

func NewStateValidationMessage(identifier machine.Identifier, reportedState machine.MachineStatus, endpoint machine.MachineEndpoint, key machine.AgentKey) *StateValidationMessage {
	now := time.Now().UTC()
	return &StateValidationMessage{
		Identifier:    identifier,
		ReportedState: reportedState,
		Time:          now,
	}
}

func (m *StateValidationMessage) AsJsonModel() *StateValidationMessageJson {
	return &StateValidationMessageJson{
		Identifier:    m.Identifier.String(),
		ReportedState: m.ReportedState.String(),
		Time:          m.Time,
	}
}

func (m *StateValidationMessage) AsJson() ([]byte, error) {
	return json.Marshal(m.AsJsonModel())
}

type StateValidationMessageJson struct {
	Identifier    string    `json:"identifier"`
	ReportedState string    `json:"reportedState"`
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
		Time:          m.Time,
	}, nil
}
