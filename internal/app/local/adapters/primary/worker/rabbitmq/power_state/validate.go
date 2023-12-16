package power_state_job

import (
	"context"
	"fmt"

	"github.com/awlsring/camp/internal/app/local/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/logger"
)

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
