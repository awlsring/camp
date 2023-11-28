package repository

import (
	"context"
	"errors"

	"github.com/awlsring/camp/apps/local/internal/core/domain/machine"
)

var (
	ErrMachineNotFound = errors.New("user does not exist")
)

type ListMachinesFilters struct{}

type Machine interface {
	Get(ctx context.Context, id machine.Identifier) (*machine.Machine, error)
	List(ctx context.Context, filters *ListMachinesFilters) ([]*machine.Machine, error)
	Add(ctx context.Context, m *machine.Machine) error
	Delete(ctx context.Context, id machine.Identifier) error
	Update(ctx context.Context, m *machine.Machine) error
	UpdateHeartbeat(ctx context.Context, id machine.Identifier) error
	UpdateStatus(ctx context.Context, id machine.Identifier, status machine.MachineStatus) error
}
