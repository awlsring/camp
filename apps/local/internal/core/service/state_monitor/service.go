package state_monitor

import (
	"github.com/awlsring/camp/apps/local/internal/ports/repository"
	"github.com/awlsring/camp/apps/local/internal/ports/service"
	"github.com/awlsring/camp/apps/local/internal/ports/topic"
)

type Service struct {
	repo          repository.Machine
	stateJobTopic topic.PowerStateJob
}

func NewService(repo repository.Machine, topic topic.PowerStateJob) service.StateMonitor {
	return &Service{
		repo: repo,
	}
}
