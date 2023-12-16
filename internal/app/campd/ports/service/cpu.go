package service

import (
	"context"

	"github.com/awlsring/camp/internal/pkg/domain/cpu"
)

type CPU interface {
	Description(ctx context.Context) (*cpu.CPU, error)
	Utilization(ctx context.Context) ([]*cpu.Utilization, error)
}
