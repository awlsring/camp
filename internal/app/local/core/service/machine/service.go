package machine

import (
	"github.com/awlsring/camp/internal/app/local/ports/repository"
	"github.com/awlsring/camp/internal/app/local/ports/service"
	"github.com/awlsring/camp/internal/app/local/ports/topic"
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
