package power_state_job

import (
	"context"

	"github.com/awlsring/camp/internal/app/local/core/domain/power"
	"github.com/awlsring/camp/internal/app/local/ports/service"
	"github.com/awlsring/camp/internal/pkg/logger"
)

type RequestStateChangeJob struct {
	powerSvc service.PowerState
}

func NewRequestStateChangeJob(powerSvc service.PowerState) *RequestStateChangeJob {
	return &RequestStateChangeJob{
		powerSvc: powerSvc,
	}
}

func (j *RequestStateChangeJob) Execute(ctx context.Context, msg []byte) error {
	log := logger.FromContext(ctx)
	log.Info().Msgf("handling message: %s", string(msg))

	log.Debug().Msgf("converting message to domain")
	message, err := j.requestStateChangeMessageToDomain(ctx, msg)
	if err != nil {
		return err
	}

	log.Debug().Msgf("routing message to correct handler")
	switch message.ChangeType {
	case power.ChangeTypeReboot:
		log.Debug().Msgf("handling reboot message")
		return j.powerSvc.Reboot(ctx, message.Identifier)
	case power.ChangeTypePowerOff:
		log.Debug().Msgf("handling power off message")
		return j.powerSvc.PowerOff(ctx, message.Identifier)
	case power.ChangeTypeWakeOnLan:
		log.Debug().Msgf("handling wake on lan message")
		return j.powerSvc.WakeOnLan(ctx, message.Identifier)
	default:
		log.Error().Msgf("invalid change type: %s", message.ChangeType)
		return power.ErrInvalidChangeType
	}
}
