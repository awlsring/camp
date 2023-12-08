package power_state_job

import (
	"context"
	"fmt"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/apps/local/internal/ports/service"
	"github.com/awlsring/camp/internal/pkg/logger"
)

type ValidateStateChangeJob struct {
	powerSvc service.PowerState
}

func NewValidateStateChangeJob(powerSvc service.PowerState) *ValidateStateChangeJob {
	return &ValidateStateChangeJob{
		powerSvc: powerSvc,
	}
}

func (j *ValidateStateChangeJob) Execute(ctx context.Context, msg []byte) error {
	log := logger.FromContext(ctx)
	log.Info().Msgf("handling message: %s", string(msg))

	log.Debug().Msgf("converting message to domain")
	message, err := j.validateStateChangeMessageToDomain(ctx, msg)
	if err != nil {
		return err
	}

	reportedState := message.ReportedState
	log.Debug().Msgf("reported state is %s. routing message to correct handler", message.ReportedState)
	switch message.ReportedState {
	case machine.MachineStatusPending, machine.MachineStatusStopping, machine.MachineStatusRebooting:
		log.Debug().Msg("handling transitional state")
		return j.powerSvc.VerifyTransitionalState(ctx, message.Identifier, reportedState)
	case machine.MachineStatusRunning, machine.MachineStatusStopped:
		log.Debug().Msg("handling final state")
		return j.powerSvc.VerifyFinalState(ctx, message.Identifier, reportedState)
	case machine.MachineStatusUnknown:
		log.Debug().Msg("handling unknown state")
		return j.powerSvc.ReconcileUnknownState(ctx, message.Identifier)
	default:
		log.Warn().Msgf("unknown state: %s", message.ReportedState)
		return fmt.Errorf("unknown state: %s", message.ReportedState)
	}
}
