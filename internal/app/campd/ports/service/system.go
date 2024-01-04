package service

import (
	"context"

	"github.com/awlsring/camp/internal/app/campd/core/domain/system"
	"github.com/awlsring/camp/internal/pkg/domain/cpu"
	"github.com/awlsring/camp/internal/pkg/domain/host"
	"github.com/awlsring/camp/internal/pkg/domain/memory"
	"github.com/awlsring/camp/internal/pkg/domain/motherboard"
)

type System interface {
	DescribeSystem(ctx context.Context) (*system.System, error)
	DescribeHost(ctx context.Context) (*host.Host, error)
	DescribeMotherboard(ctx context.Context) (*motherboard.Motherboard, error)
	DescribeBios(ctx context.Context) (*motherboard.Bios, error)
	DescribeCpu(ctx context.Context) ([]*cpu.CPU, error)
	DescribeMemory(ctx context.Context) ([]*memory.Memory, error)
	
}
