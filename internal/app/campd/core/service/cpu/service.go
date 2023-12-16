package cpu

import (
	"context"

	"github.com/awlsring/camp/internal/app/campd/ports/service"
	"github.com/awlsring/camp/internal/pkg/domain/cpu"
)

type Service struct {
	cpu *cpu.CPU
}

func InitService(ctx context.Context) (service.CPU, error) {
	cpu, err := loadCPU(ctx)
	if err != nil {
		return nil, err
	}

	return &Service{
		cpu: cpu,
	}, nil
}
