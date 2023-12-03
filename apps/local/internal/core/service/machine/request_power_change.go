package machine

import (
	"context"
	"fmt"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/apps/local/internal/core/domain/power"
	"github.com/awlsring/camp/internal/pkg/errors"
	"github.com/awlsring/camp/internal/pkg/logger"
)

var ErrUnsupportedCapability = fmt.Errorf("unsupported capability")

func (s *machineService) RequestPowerChange(ctx context.Context, id machine.Identifier, changeType power.ChangeType) error {
	log := logger.FromContext(ctx)
	log.Info().Msgf("Initiating power change for machine %s", id)

	log.Debug().Msg("Getting machine from repository")
	m, err := s.repo.Get(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get machine from repository")
		return err
	}

	log.Debug().Msg("Validating capability is enabled and building message")
	var msg *power.PowerChangeRequestMessage
	switch changeType {
	case power.ChangeTypeWakeOnLan:
		if !m.PowerCapabilities.WakeOnLan.Enabled {
			log.Error().Msg("Machine does not support wake on lan")
			return errors.New(ErrUnsupportedCapability, fmt.Errorf("machine %s does not support wake on lan", id))
		}
		if m.PowerCapabilities.WakeOnLan.MacAddress == nil {
			log.Error().Msg("Machine does not have a mac address to send wake on lan to")
			return errors.New(errors.ErrInternalFailure, fmt.Errorf("machine %s does not have a mac address to send wake on lan to", id))
		}
		mac := string(*m.PowerCapabilities.WakeOnLan.MacAddress)
		msg = power.NewWakeOnLanMessage(id, mac)
	case power.ChangeTypePowerOff:
		if !m.PowerCapabilities.PowerOff.Enabled {
			log.Error().Msg("Machine does not support power off")
			return errors.New(ErrUnsupportedCapability, fmt.Errorf("machine %s does not support power off", id))
		}
		msg = power.NewPowerOffMessage(id, m.AgentEndpoint.String(), m.AgentApiKey.String())
	case power.ChangeTypeReboot:
		if !m.PowerCapabilities.Reboot.Enabled {
			log.Error().Msg("Machine does not support reboot")
			return errors.New(ErrUnsupportedCapability, fmt.Errorf("machine %s does not support reboot", id))
		}
		msg = power.NewRebootMessage(id, m.AgentEndpoint.String(), m.AgentApiKey.String())
	}

	log.Debug().Msg("Posting power change request to topic")
	err = s.powerTopic.SendPowerChangeRequest(ctx, msg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to post power change request to topic")
		return err
	}

	log.Debug().Msg("Successfully posted power change request to topic")
	return nil
}
