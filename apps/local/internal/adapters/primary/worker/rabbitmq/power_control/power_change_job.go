package rabbitmq_power_control

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/apps/local/internal/core/domain/power"
	"github.com/awlsring/camp/apps/local/internal/ports/service"
	"github.com/awlsring/camp/internal/pkg/logger"
)

type PowerChangeRequestJob struct {
	powerSvc service.PowerControl
}

func NewPowerChangeRequestJob(powerSvc service.PowerControl) *PowerChangeRequestJob {
	return &PowerChangeRequestJob{
		powerSvc: powerSvc,
	}
}

func (p *PowerChangeRequestJob) Execute(ctx context.Context, msg []byte) error {
	log := logger.FromContext(ctx)
	log.Info().Msgf("handling message: %s", string(msg))

	log.Debug().Msgf("unmarshalling message")
	var message power.RequestStateChangeMessageJson
	err := json.Unmarshal(msg, &message)
	if err != nil {
		log.Error().Err(err).Msgf("failed to unmarshal message")
		return err
	}

	log.Debug().Msgf("converting message to domain")
	domainMessage, err := message.ToDomain()
	if err != nil {
		log.Error().Err(err).Msgf("failed to convert message to domain")
		return err
	}

	log.Debug().Msgf("handling message")
	switch domainMessage.ChangeType {
	case power.ChangeTypeReboot:
		log.Debug().Msgf("handling reboot message")
		return p.reboot(ctx, domainMessage)
	case power.ChangeTypePowerOff:
		log.Debug().Msgf("handling power off message")
		return p.powerOff(ctx, domainMessage)
	case power.ChangeTypeWakeOnLan:
		log.Debug().Msgf("handling wake on lan message")
		return p.wakeOnLan(ctx, domainMessage)
	default:
		log.Error().Msgf("invalid change type: %s", domainMessage.ChangeType)
		return power.ErrInvalidChangeType
	}
}

func validateAgentParameters(ctx context.Context, endpoint string, key *string) (machine.MachineEndpoint, machine.AgentKey, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("validating key is set")
	if key == nil {
		return "", "", fmt.Errorf("key is nil, is required for validation")
	}
	log.Debug().Msgf("validating key is valid")
	k, err := machine.AgentKeyFromString(*key)
	if err != nil {
		log.Error().Err(err).Msgf("key is invalid")
		return "", "", err
	}

	log.Debug().Msgf("validating endpoint is valid")
	e, err := machine.MachineEndpointFromString(endpoint)
	if err != nil {
		log.Error().Err(err).Msgf("endpoint is invalid")
		return "", "", err
	}

	return e, k, nil
}

func (p *PowerChangeRequestJob) pollTillChanged(ctx context.Context, deadline time.Time, id machine.Identifier, endpoint machine.MachineEndpoint, key machine.AgentKey, state bool) error {
	log := logger.FromContext(ctx)

	for {
		if time.Now().After(deadline) {
			return fmt.Errorf("machine %s failed to reach state within timeout", id.String())
		}

		log.Debug().Msgf("validating machine %s is powere", id.String())
		on, err := p.powerSvc.CheckMachinePower(ctx, id, endpoint, key)
		if err != nil {
			return err
		}
		if on == state {
			log.Debug().Msgf("machine %s has reached state", id.String())
			break
		}

		log.Debug().Msgf("machine %s is not at desired state, waiting 5 seconds", id.String())
		time.Sleep(5 * time.Second)
	}

	return nil
}

func (p *PowerChangeRequestJob) powerOff(ctx context.Context, msg *power.RequestStateChangeMessage) error {
	log := logger.FromContext(ctx)
	log.Info().Msgf("powering off machine %s", msg.Identifier.String())
	log.Debug().Msgf("validating message")

	log.Debug().Msg("validating key and endpoint")
	endpoint, key, err := validateAgentParameters(ctx, msg.Target, msg.Key)
	if err != nil {
		return err
	}

	log.Debug().Msgf("Sending powering off request")
	err = p.powerSvc.PowerOff(ctx, msg.Identifier, endpoint, key)
	if err != nil {
		return err
	}
	log.Debug().Msgf("power off request sent to machine %s", msg.Identifier.String())

	deadline := time.Now().Add(time.Duration(msg.Timeout) * time.Second)
	log.Debug().Msgf("Waiting timeout duration for machine %s to power off", msg.Identifier.String())

	err = p.pollTillChanged(ctx, deadline, msg.Identifier, endpoint, key, false)
	if err != nil {
		return err
	}

	log.Debug().Msgf("machine %s is powered off", msg.Identifier.String())
	return err
}

func (p *PowerChangeRequestJob) reboot(ctx context.Context, msg *power.RequestStateChangeMessage) error {
	log := logger.FromContext(ctx)
	log.Info().Msgf("rebooting machine %s", msg.Identifier.String())
	log.Debug().Msgf("validating message")

	log.Debug().Msg("validating key and endpoint")
	endpoint, key, err := validateAgentParameters(ctx, msg.Target, msg.Key)
	if err != nil {
		return err
	}

	log.Debug().Msgf("Sending reboot request")
	err = p.powerSvc.Reboot(ctx, msg.Identifier, endpoint, key)

	deadline := time.Now().Add(time.Duration(msg.Timeout) * time.Second)
	log.Debug().Msgf("Waiting timeout duration for machine %s to power off", msg.Identifier.String())
	err = p.pollTillChanged(ctx, deadline, msg.Identifier, endpoint, key, false)
	if err != nil {
		return err
	}

	log.Debug().Msgf("Waiting timeout duration for machine %s to power on", msg.Identifier.String())
	err = p.pollTillChanged(ctx, deadline, msg.Identifier, endpoint, key, true)
	if err != nil {
		return err
	}

	log.Debug().Msgf("reboot request sent to machine %s", msg.Identifier.String())
	return err
}

func (p *PowerChangeRequestJob) wakeOnLan(ctx context.Context, msg *power.RequestStateChangeMessage) error {
	log := logger.FromContext(ctx)
	log.Info().Msgf("waking up machine %s", msg.Identifier.String())
	log.Debug().Msgf("validating message")

	log.Debug().Msgf("validating endpoint is valid")
	mac, err := machine.MacAddressFromString(msg.Target)
	if err != nil {
		return err
	}

	log.Debug().Msgf("Sending wake on lan request")
	err = p.powerSvc.WakeOnLan(ctx, msg.Identifier, mac)

	// deadline := time.Now().Add(time.Duration(msg.Timeout) * time.Second)
	// log.Debug().Msgf("Waiting timeout duration for machine %s to power off", msg.Identifier.String())

	// log.Debug().Msgf("Waiting timeout duration for machine %s to power on", msg.Identifier.String())
	// err = p.pollTillChanged(ctx, deadline, msg.Identifier, endpoint, key, true)
	// if err != nil {
	// 	return err
	// }

	log.Debug().Msgf("wake on lan request sent to machine %s", msg.Identifier.String())
	return err
}
