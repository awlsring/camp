package machine

import (
	"github.com/awlsring/camp/apps/local/internal/ports/repository"
	"github.com/awlsring/camp/apps/local/internal/ports/service"
	"github.com/awlsring/camp/apps/local/internal/ports/topic"
)

var _ service.Machine = &machineService{}

type machineService struct {
	stateChangeTopic        topic.PowerStateChange
	stateChangeRequestTopic topic.PowerStateJob
	repo                    repository.Machine
}

func NewMachineService(repo repository.Machine) service.Machine {
	return &machineService{
		repo: repo,
	}
}
