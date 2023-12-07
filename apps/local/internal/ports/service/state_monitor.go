package service

import "context"

type StateMonitor interface {
	ScheduleStateVerificationJobs(ctx context.Context) error
}
