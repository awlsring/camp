package state_monitor

import (
	"github.com/awlsring/camp/apps/local/internal/ports/repository"
	"github.com/awlsring/camp/apps/local/internal/ports/service"
	"github.com/awlsring/camp/apps/local/internal/ports/topic"
	"github.com/awlsring/camp/internal/pkg/agent"
)

type Service struct {
	repo             repository.Machine
	stateChangeTopic topic.PowerStateChange
	campdClient      agent.Client
}

func NewService(repo repository.Machine, stateChangeTopic topic.PowerStateChange, campd agent.Client) service.StateMonitor {
	return &Service{
		repo:             repo,
		stateChangeTopic: stateChangeTopic,
		campdClient:      campd,
	}
}
