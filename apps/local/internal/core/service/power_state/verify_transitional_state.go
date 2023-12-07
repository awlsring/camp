package power_state

import (
	"context"
	"time"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/logger"
)

const (
	PendingTimeout  = 1 * time.Minute
	StartingTimeout = 5 * time.Minute
	StoppingTimeout = 5 * time.Minute
	RebootTimeout   = StartingTimeout + StoppingTimeout
)

func (s *Service) VerifyTransitionalState(ctx context.Context, id machine.Identifier, started, deadline time.Time, state machine.MachineStatus, endpoint machine.MachineEndpoint, token machine.AgentKey) error {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("validating transitional state %s", state.String())

	now := time.Now().UTC()

	switch state {
	case machine.MachineStatusPending:
		log.Debug().Msg("machine is pending, checking timeout")
		timeout := started.Add(PendingTimeout)
		if !now.After(timeout) {
			log.Debug().Msg("machine is within deadline, ignoring")
			return nil
		}
	case machine.MachineStatusStopping, machine.MachineStatusRebooting, machine.MachineStatusStarting:
		log.Debug().Msg("machine is transitioning, checking timeout")
		if !now.After(deadline) {
			log.Debug().Msg("machine is still within deadline, ignoring")
			return nil
		}
	}

	return s.VerifyState(ctx, id, state, endpoint, token)
}
