package memory

import (
	"context"

	"github.com/awlsring/camp/internal/app/campd/ports/service"
	"github.com/shirou/gopsutil/mem"
)

type Service struct {
	total uint64
}

func InitService(ctx context.Context) (service.Memory, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	return &Service{
		total: v.Total,
	}, nil
}
