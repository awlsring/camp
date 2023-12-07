package schedule

import (
	"context"
	"sync"
	"time"

	"github.com/awlsring/camp/apps/local/internal/ports/service"
	"github.com/awlsring/camp/internal/pkg/logger"
)

const DefaultInterval = 60 * time.Second

type SchedulerOpt func(*Scheduler)

func WithInterval(interval time.Duration) SchedulerOpt {
	return func(s *Scheduler) {
		s.interval = interval
	}
}

type Scheduler struct {
	svc      service.StateMonitor
	interval time.Duration
}

func NewScheduler(svc service.StateMonitor, opts ...SchedulerOpt) *Scheduler {
	s := &Scheduler{
		svc:      svc,
		interval: DefaultInterval,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Scheduler) Start(ctx context.Context) error {
	log := logger.FromContext(ctx)
	log.Info().Msg("Starting scheduler")
	log.Info().Msgf("Interval is set to %s", s.interval.String())

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			err := s.svc.ScheduleStateVerificationJobs(ctx)
			if err != nil {
				log.Error().Err(err).Msg("failed to verify and adjust machine states")
			}
			select {
			case <-ctx.Done():
				log.Debug().Msg("context done, exiting")
				return
			case <-time.After(60 * time.Second):
				log.Debug().Msg("sleeping for 60 seconds")
			}
		}
	}()

	// wait for shutdown signal
	<-ctx.Done()
	log.Info().Msg("Shutting down")
	wg.Wait()
	log.Info().Msg("Shutdown complete")

	return nil
}
