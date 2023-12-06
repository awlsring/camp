package state_monitor

import (
	"context"
	"sync"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/apps/local/internal/core/domain/power"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *Service) VerifyAndAdjustMachineStates(ctx context.Context) error {
	log := logger.FromContext(ctx)

	machines, err := s.repo.List(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to list machines")
		return err
	}

	var wg sync.WaitGroup
	for _, m := range machines {
		wg.Add(1)
		go func(machine *machine.Machine) {
			defer wg.Done()
			s.verifyAndAdjustMachineState(ctx, machine)
		}(m)
	}

	wg.Wait()
	return nil
}

func (s *Service) verifyAndAdjustMachineState(ctx context.Context, m *machine.Machine) error {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("verifying machine %s state", m.Identifier)

	currentState := m.Status
	log.Debug().Msgf("machine %s current state is %s", m.Identifier, currentState.String())

	var actualState machine.MachineStatus
	var err error
	log.Debug().Msgf("routing to state verification for state %s", m.Status.String())
	switch currentState {
	case machine.MachineStatusPending, machine.MachineStatusStopping, machine.MachineStatusRebooting:
		log.Debug().Msgf("machine %s is in a transitional state", m.Identifier)
		actualState, err = s.verifyTransitionalState(ctx, m)
	case machine.MachineStatusRunning, machine.MachineStatusStopped:
		log.Debug().Msgf("machine %s is in a final state", m.Identifier)
		actualState, err = s.verifyFinalState(ctx, m)
	default:
		log.Debug().Msgf("machine %s is in an unknown state", m.Identifier)
		actualState, err = s.fixUnknownState(ctx, m)
	}
	log.Debug().Msgf("machine %s actual state is %s", m.Identifier, actualState.String())

	if actualState != currentState {
		log.Debug().Msgf("machine %s state is %s, should be %s. Send unplanned state message", m.Identifier, actualState.String(), currentState.String())
		err := s.sendUnplannedStateChangeMessage(ctx, m, actualState)
		if err != nil {
			log.Error().Err(err).Msgf("failed to send unplanned state change message for machine %s", m.Identifier)
			return err
		}
	}

	log.Debug().Msgf("done verifying machine %s state", m.Identifier)
	return err
}

func (s *Service) verifyTransitionalState(ctx context.Context, m *machine.Machine) (machine.MachineStatus, error) {
	// add more details to status, maybe make an event database to help dictate action
	panic("not implemented")
}

func (s *Service) verifyFinalState(ctx context.Context, m *machine.Machine) (machine.MachineStatus, error) {
	log := logger.FromContext(ctx)

	shouldConnect := m.Status == machine.MachineStatusRunning
	log.Debug().Msgf("machine %s state is %s, connectivity should return %t", m.Identifier, m.Status.String(), shouldConnect)

	log.Debug().Msgf("checking connectivity for machine %s", m.Identifier)
	ok, err := s.campdClient.CheckMachineConnectivity(ctx, m.Identifier.String(), m.AgentEndpoint.String(), m.AgentApiKey.String())
	if err != nil {
		log.Error().Err(err).Msgf("failed to check connectivity for machine %s", m.Identifier)
		return machine.MachineStatusUnknown, err
	}

	log.Debug().Msgf("machine %s connectivity check returned %t", m.Identifier, ok)
	if ok == shouldConnect {
		return m.Status, nil
	}

	var actualState machine.MachineStatus
	if ok {
		actualState = machine.MachineStatusRunning
	} else {
		actualState = machine.MachineStatusStopped
	}

	log.Debug().Msgf("machine %s state is %s, should be %s", m.Identifier, actualState.String(), m.Status.String())
	return actualState, nil
}

func (s *Service) fixUnknownState(ctx context.Context, m *machine.Machine) (machine.MachineStatus, error) {
	log := logger.FromContext(ctx)

	log.Debug().Msgf("checking connectivity for machine %s", m.Identifier)
	ok, err := s.campdClient.CheckMachineConnectivity(ctx, m.Identifier.String(), m.AgentEndpoint.String(), m.AgentApiKey.String())
	if err != nil {
		log.Error().Err(err).Msgf("failed to check connectivity for machine %s", m.Identifier)
		return machine.MachineStatusUnknown, err
	}

	log.Debug().Msgf("machine %s connectivity check returned %t", m.Identifier, ok)
	if ok {
		return machine.MachineStatusRunning, nil
	}

	return machine.MachineStatusStopped, nil
}

func (s *Service) sendUnplannedStateChangeMessage(ctx context.Context, machine *machine.Machine, state machine.MachineStatus) error {
	log := logger.FromContext(ctx)

	log.Debug().Msgf("Creating state change message")
	msg := power.NewStateChangeMessage(machine.Identifier, machine.Status, state, false)

	log.Debug().Msgf("Sending state change message")
	err := s.stateChangeTopic.SendStateChangeMessage(ctx, msg)
	if err != nil {
		log.Error().Err(err).Msgf("failed to send state change message for machine %s", machine.Identifier)
		return err
	}

	return nil
}
