package machine

import "context"

type Controller interface {
	RegisterMachine(ctx context.Context, input *RegisterMachineInput) (*Model, error)
}

type RegisterMachineInput struct {
	Model *Model
}

type MachineController struct {
	repo Repo
}

func NewController() Controller {
	return &MachineController{}
}

func (c *MachineController) RegisterMachine(ctx context.Context, input *RegisterMachineInput) (*Model, error) {
	err := c.repo.CreateMachine(ctx, input.Model)
	return input.Model, err
}
