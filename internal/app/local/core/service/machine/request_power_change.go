package machine

import (
	"context"
	"fmt"

	"github.com/awlsring/camp/internal/app/local/core/domain/machine"
	pwr "github.com/awlsring/camp/internal/app/local/core/domain/power"
	"github.com/awlsring/camp/internal/pkg/domain/power"
	"github.com/awlsring/camp/internal/pkg/errors"
	"github.com/awlsring/camp/internal/pkg/logger"
)

var (
	ErrUnsupportedCapability = fmt.Errorf("unsupported capability")
	ErrInvalidStatus         = fmt.Errorf("invalid machine status")
)

func (s *machineService) RequestPowerChange(ctx context.Context, id machine.Identifier, changeType pwr.ChangeType) error {
	log := logger.FromContext(ctx)
	log.Info().Msgf("Initiating power change for machine %s", id)

	log.Debug().Msg("Getting machine from repository")
	m, err := s.repo.Get(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get machine from repository")
		return err
	}

	log.Debug().Msg("Validating capability is enabled and building message")
	var msg *pwr.RequestStateChangeMessage
	switch changeType {
	case pwr.ChangeTypeWakeOnLan:
		msg, err = formWakeOnLanMessage(ctx, m)
	case pwr.ChangeTypePowerOff:
		msg, err = formPowerOffMessage(ctx, m)
	case pwr.ChangeTypeReboot:
		msg, err = formRebootMessage(ctx, m)
	default:
		log.Error().Msg("Unsupported power change type")
		return errors.New(errors.ErrValidation, fmt.Errorf("unsupported power change type %s", changeType))
	}
	if err != nil {
		return err
	}

	log.Debug().Msg("Posting power change request to topic")
	err = s.stateChangeRequestTopic.SendRequestStateChangeMessage(ctx, msg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to post power change request to topic")
		return err
	}

	log.Debug().Msg("Successfully posted power change request to topic")
	return nil
}

func formWakeOnLanMessage(ctx context.Context, m *machine.Machine) (*pwr.RequestStateChangeMessage, error) {
	log := logger.FromContext(ctx)

	log.Debug().Msg("Validating WOL capability is enabled")
	if !m.PowerCapabilities.WakeOnLan.Enabled {
		log.Error().Msg("Machine does not support wake on lan")
		return nil, errors.New(ErrUnsupportedCapability, fmt.Errorf("machine %s does not support wake on lan", m.Identifier))
	}
	if m.PowerCapabilities.WakeOnLan.MacAddress == nil {
		log.Error().Msg("Machine does not have a mac address to send wake on lan to")
		return nil, errors.New(errors.ErrInternalFailure, fmt.Errorf("machine %s does not have a mac address to send wake on lan to", m.Identifier))
	}

	log.Debug().Msg("Validating machine is Stopped or Unknown")
	if m.Status.Status != power.StatusCodeStopped && m.Status.Status != power.StatusCodeUnknown {
		log.Error().Msg("Machine is not stopped and cannot be woken up")
		return nil, errors.New(ErrInvalidStatus, fmt.Errorf("machine %s is not stopped and cannot be woken up", m.Identifier))
	}

	log.Debug().Msg("Building message")
	mac := string(*m.PowerCapabilities.WakeOnLan.MacAddress)
	msg := pwr.NewWakeOnLanMessage(m.Identifier, mac)

	log.Debug().Msg("Successfully built message")
	return msg, nil
}

func formPowerOffMessage(ctx context.Context, m *machine.Machine) (*pwr.RequestStateChangeMessage, error) {
	log := logger.FromContext(ctx)

	log.Debug().Msg("Validating capability is enabled")
	if !m.PowerCapabilities.PowerOff.Enabled {
		log.Error().Msg("Machine does not support power off")
		return nil, errors.New(ErrUnsupportedCapability, fmt.Errorf("machine %s does not support power off", m.Identifier))
	}

	log.Debug().Msg("Validating machine is running")
	if m.Status.Status != power.StatusCodeRunning {
		log.Error().Msg("Machine is not running and cannot be powered off")
		return nil, errors.New(ErrInvalidStatus, fmt.Errorf("machine %s is not running and cannot be powered off", m.Identifier))
	}

	log.Debug().Msg("Building message")
	msg := pwr.NewPowerOffMessage(m.Identifier, m.AgentEndpoint.String(), m.AgentApiKey.String())

	log.Debug().Msg("Successfully built message")
	return msg, nil
}

func formRebootMessage(ctx context.Context, m *machine.Machine) (*pwr.RequestStateChangeMessage, error) {
	log := logger.FromContext(ctx)

	log.Debug().Msg("Validating capability is enabled")
	if !m.PowerCapabilities.Reboot.Enabled {
		log.Error().Msg("Machine does not support reboot")
		return nil, errors.New(ErrUnsupportedCapability, fmt.Errorf("machine %s does not support reboot", m.Identifier))
	}

	log.Debug().Msg("Validating machine is running")
	if m.Status.Status != power.StatusCodeRunning {
		log.Error().Msg("Machine is not running and cannot be rebooted")
		return nil, errors.New(ErrInvalidStatus, fmt.Errorf("machine %s is not running and cannot be rebooted", m.Identifier))
	}

	log.Debug().Msg("Building message")
	msg := pwr.NewRebootMessage(m.Identifier, m.AgentEndpoint.String(), m.AgentApiKey.String())

	log.Debug().Msg("Successfully built message")
	return msg, nil
}
