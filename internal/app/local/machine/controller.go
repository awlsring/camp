package machine

import (
	"context"

	"github.com/rs/zerolog/log"
)

type Controller interface {
	RegisterMachine(ctx context.Context, machine *Model) error
	DescribeMachine(ctx context.Context, id string) (*Model, error)
	ListMachines(ctx context.Context, input *GetMachinesFilters) ([]*Model, error)
	AcknowledgeHeartbeat(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status MachineStatus) error
	UpdateMachine(ctx context.Context, machine *Model) error
}

type MachineController struct {
	repo Repo
}

func NewController(repo Repo) Controller {
	return &MachineController{
		repo: repo,
	}
}

func (c *MachineController) RegisterMachine(ctx context.Context, model *Model) error {
	return c.repo.CreateMachine(ctx, model)
}

func (c *MachineController) DescribeMachine(ctx context.Context, identifier string) (*Model, error) {
	log.Debug().Msg("Invoke Controller.DescribeMachine")
	return c.repo.GetMachineById(ctx, identifier)
}

func (c *MachineController) ListMachines(ctx context.Context, filters *GetMachinesFilters) ([]*Model, error) {
	return c.repo.GetMachines(ctx, filters)
}

func (c *MachineController) AcknowledgeHeartbeat(ctx context.Context, id string) error {
	return c.repo.AcknowledgeHeartbeat(ctx, id)
}

func (c *MachineController) UpdateStatus(ctx context.Context, id string, status MachineStatus) error {
	return c.repo.UpdateStatus(ctx, id, status)
}

func (c *MachineController) UpdateMachine(ctx context.Context, model *Model) error {
	return c.repo.UpdateMachine(ctx, model)
}
