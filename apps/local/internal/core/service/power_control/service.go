package power_control

import (
	"github.com/awlsring/camp/internal/pkg/agent"
	"github.com/awlsring/camp/internal/pkg/wol"

	"github.com/awlsring/camp/apps/local/internal/ports/repository"
	"github.com/awlsring/camp/apps/local/internal/ports/service"
)

var _ service.PowerControl = &powerControlService{}

type powerControlService struct {
	wol   wol.Client
	mRepo repository.Machine
	agent agent.Client
}

func NewPowerControlService(wol wol.Client, repo repository.Machine, agent agent.Client) service.PowerControl {
	return &powerControlService{
		wol:   wol,
		mRepo: repo,
		agent: agent,
	}
}
