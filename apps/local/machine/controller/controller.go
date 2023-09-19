package controller

import (
	"context"

	"github.com/awlsring/camp/apps/local/machine"
	"github.com/awlsring/camp/apps/local/machine/repo"
	"github.com/rs/zerolog/log"
)

type Controller interface {
	RegisterMachine(ctx context.Context, machine *machine.Model) error
	DescribeMachine(ctx context.Context, id string) (*machine.Model, error)
	ListMachines(ctx context.Context, input *repo.GetMachinesFilters) ([]*machine.Model, error)
	AcknowledgeHeartbeat(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status machine.MachineStatus) error
	UpdateMachine(ctx context.Context, machine *machine.Model) error
}

type MachineController struct {
	repo repo.Repo
}

func NewController(repo repo.Repo) Controller {
	return &MachineController{
		repo: repo,
	}
}

func (c *MachineController) RegisterMachine(ctx context.Context, model *machine.Model) error {
	log.Debug().Msg("Invoke Controller.RegisterMachine")
	return c.repo.CreateMachine(ctx, model)
}

func (c *MachineController) DescribeMachine(ctx context.Context, identifier string) (*machine.Model, error) {
	log.Debug().Msg("Invoke Controller.DescribeMachine")
	return c.repo.GetMachineById(ctx, identifier)
}

func (c *MachineController) ListMachines(ctx context.Context, filters *repo.GetMachinesFilters) ([]*machine.Model, error) {
	log.Debug().Msg("Invoke Controller.ListMachines")
	return c.repo.GetMachines(ctx, filters)
}

func (c *MachineController) AcknowledgeHeartbeat(ctx context.Context, id string) error {
	log.Debug().Msg("Invoke Controller.AcknowledgeHeartbeat")
	return c.repo.AcknowledgeHeartbeat(ctx, id)
}

func (c *MachineController) UpdateStatus(ctx context.Context, id string, status machine.MachineStatus) error {
	log.Debug().Msg("Invoke Controller.UpdateStatus")
	return c.repo.UpdateStatus(ctx, id, status)
}

func (c *MachineController) UpdateMachine(ctx context.Context, model *machine.Model) error {
	log.Debug().Msg("Invoke Controller.UpdateMachine")
	return c.repo.UpdateMachine(ctx, model)
}
