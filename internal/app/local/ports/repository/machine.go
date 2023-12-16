package repository

import (
	"context"
	"errors"

	"github.com/awlsring/camp/internal/app/local/core/domain/machine"
	"github.com/awlsring/camp/internal/pkg/domain/power"
	"github.com/awlsring/camp/internal/pkg/domain/tag"
)

var (
	ErrMachineNotFound = errors.New("machine does not exist")
	ErrInternalFailure = errors.New("internal error")
)

type ListMachinesFilters struct{}

type Machine interface {
	Get(ctx context.Context, id machine.Identifier) (*machine.Machine, error)
	List(ctx context.Context, filters *ListMachinesFilters) ([]*machine.Machine, error)
	Add(ctx context.Context, m *machine.Machine) error
	AddTags(ctx context.Context, id machine.Identifier, tags []*tag.Tag) error
	RemoveTag(ctx context.Context, id machine.Identifier, tag tag.TagKey) error
	Delete(ctx context.Context, id machine.Identifier) error
	Update(ctx context.Context, m *machine.Machine) error
	UpdateHeartbeat(ctx context.Context, id machine.Identifier) error
	UpdateStatus(ctx context.Context, id machine.Identifier, status power.StatusCode) error
}
