package power_state

import (
	"github.com/awlsring/camp/apps/local/internal/ports/repository"
	"github.com/awlsring/camp/apps/local/internal/ports/service"
	"github.com/awlsring/camp/apps/local/internal/ports/topic"
	"github.com/awlsring/camp/internal/pkg/agent"
	"github.com/awlsring/camp/internal/pkg/wol"
)

type Service struct {
	repo  repository.Machine
	campd agent.Client
	wol   wol.Client

	changeTopic topic.PowerStateChange
}

func New(repo repository.Machine, campd agent.Client, wol wol.Client, changeTopic topic.PowerStateChange) service.PowerState {
	return &Service{
		repo:        repo,
		campd:       campd,
		wol:         wol,
		changeTopic: changeTopic,
	}
}
