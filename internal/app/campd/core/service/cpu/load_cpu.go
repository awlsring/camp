package cpu

import (
	"context"
	"fmt"
	"runtime"

	"github.com/awlsring/camp/internal/pkg/domain/cpu"
	"github.com/awlsring/camp/internal/pkg/logger"
	"github.com/awlsring/camp/internal/pkg/sys"
	"github.com/jaypipes/ghw"
	pscpu "github.com/shirou/gopsutil/v3/cpu"
)

func loadStandard(ctx context.Context, arch cpu.Architecture) (*cpu.CPU, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Getting CPU via standard process")

	i, err := ghw.CPU()
	if err != nil {
		log.Error().Err(err).Msg("failed to get CPU info")
		return nil, err
	}

	vendor := ""
	model := ""

	procs := []*cpu.Processor{}
	for _, p := range i.Processors {
		vendor = p.Vendor
		model = p.Model

		cores := []*cpu.Core{}
		for _, c := range p.Cores {
			cores = append(cores, cpu.NewCore(
				c.ID,
				c.NumThreads,
			))
		}

		procs = append(procs, cpu.NewProcessor(
			p.ID,
			p.NumCores,
			p.NumThreads,
			p.Vendor,
			p.Model,
			cores,
		))
	}

	cpu := cpu.NewCPU(i.TotalCores, i.TotalThreads, arch, vendor, model, procs)

	return cpu, nil
}

func loadCPUDarwin(ctx context.Context, arch cpu.Architecture) (*cpu.CPU, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Getting CPU for apple")

	c, err := pscpu.Info()
	if err != nil {
		log.Error().Err(err).Msg("failed to get CPU info")
		return nil, err
	}

	if len(c) != 1 {
		err := fmt.Errorf("expected 1 CPU, got %d", len(c))
		return nil, err
	}

	cpuResult := c[0]
	vendor := "Apple"
	return &cpu.CPU{
		TotalCores:   uint32(cpuResult.Cores),
		Vendor:       &vendor,
		Architecture: arch,
		Model:        &cpuResult.ModelName,
	}, nil
}

func loadCPU(ctx context.Context) (*cpu.CPU, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Initializing CPU service")

	arch := cpu.ArchitectureFromString(runtime.GOARCH)

	if sys.IsMacOS() {
		log.Warn().Msg("CPU service not fully implemented for macOS, returning partial results")
		return loadCPUDarwin(ctx, arch)
	}

	return loadStandard(ctx, arch)
}
