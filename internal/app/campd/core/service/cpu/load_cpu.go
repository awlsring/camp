package cpu

import (
	"context"
	"runtime"

	"github.com/awlsring/camp/internal/pkg/domain/cpu"
	"github.com/awlsring/camp/internal/pkg/logger"
	"github.com/jaypipes/ghw"
)

func loadCPU(ctx context.Context) (*cpu.CPU, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Initializing CPU service")

	arch := cpu.ArchitectureFromString(runtime.GOARCH)

	i, err := ghw.CPU()
	if err != nil {
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
