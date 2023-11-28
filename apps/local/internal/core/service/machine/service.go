package machine

import (
	"github.com/awlsring/camp/apps/local/internal/ports/repository"
	"github.com/awlsring/camp/apps/local/internal/ports/service"
)

var _ service.Machine = &machineService{}

type machineService struct {
	repo repository.Machine
}

func NewMachineService(repo repository.Machine) service.Machine {
	return &machineService{
		repo: repo,
	}
}
