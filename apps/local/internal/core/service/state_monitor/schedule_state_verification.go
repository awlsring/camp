package state_monitor

import (
	"context"
	"sync"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
	"github.com/awlsring/camp/apps/local/internal/core/domain/power"
	"github.com/awlsring/camp/internal/pkg/logger"
)

func (s *Service) ScheduleStateVerificationJobs(ctx context.Context) error {
	log := logger.FromContext(ctx)
	log.Debug().Msg("scheduling state verification jobs")

	log.Debug().Msg("listing machines")
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
			s.scheduleStateVerificationJob(ctx, machine)
		}(m)
	}
	log.Debug().Msg("waiting for state verification jobs to send")
	wg.Wait()
	log.Debug().Msg("state verification jobs sent")
	return nil
}

func (s *Service) scheduleStateVerificationJob(ctx context.Context, m *machine.Machine) error {
	log := logger.FromContext(ctx)
	log.Debug().Msgf("scheduling state verification job for machine %s", m.Identifier)

	log.Debug().Msgf("forming state verification job for machine %s", m.Identifier)
	job := power.NewStateValidationMessage(m.Identifier, m.Status, m.AgentEndpoint, m.AgentApiKey)

	log.Debug().Msgf("publishing state verification job for machine %s", m.Identifier)
	err := s.stateJobTopic.SendStateValidationMessage(ctx, job)
	if err != nil {
		log.Error().Err(err).Msgf("failed to publish state verification job for machine %s", m.Identifier)
		return err
	}

	log.Debug().Msgf("successfully scheduled state verification job for machine %s", m.Identifier)
	return nil
}
