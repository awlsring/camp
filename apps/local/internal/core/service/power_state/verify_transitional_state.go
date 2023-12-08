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

func (s *Service) VerifyTransitionalState(ctx context.Context, id machine.Identifier, state machine.MachineStatus) error {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("validating transitional state %s", state.String())

	log.Debug().Msg("getting machine entry")
	m, err := s.repo.Get(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get machine entry")
		return err
	}

	now := time.Now().UTC()
	var timeout time.Duration
	switch state {
	case machine.MachineStatusPending:
		log.Debug().Msg("machine is pending, checking last update")
		timeout = PendingTimeout
	case machine.MachineStatusStopping:
		log.Debug().Msg("machine is stopping, checking last update")
		timeout = StoppingTimeout
	case machine.MachineStatusStarting:
		log.Debug().Msg("machine is starting, checking last update")
		timeout = StartingTimeout
	case machine.MachineStatusRebooting:
		log.Debug().Msg("machine is rebooting, checking last update")
		timeout = RebootTimeout
	default:
		log.Debug().Msg("machine is not in transitional state, ignoring")
		return nil
	}

	if !m.UpdatedAt.Add(timeout).Before(now) {
		log.Debug().Msg("machine is still within timeout, ignoring")
		return nil
	}

	log.Debug().Msg("machine is outside of timeout, checking state")
	return s.verifyState(ctx, m, state)
}
