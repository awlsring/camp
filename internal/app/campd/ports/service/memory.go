package service

import (
	"context"

	"github.com/awlsring/camp/internal/pkg/domain/memory"
)

type Memory interface {
	Description(ctx context.Context) (*memory.Memory, error)
	Utilization(ctx context.Context) (*memory.Utilization, error)
}
