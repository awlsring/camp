package model

import (
	"github.com/awlsring/camp/internal/pkg/domain/cpu"
	campd "github.com/awlsring/camp/pkg/gen/campd_grpc"
	grpcmodel "github.com/awlsring/camp/pkg/grpc_model"
)

func NewArchitecture(in cpu.Architecture) campd.Architecture {
	switch in {
	case cpu.ArchitectureArm64:
		return campd.Architecture_ARM_64
	case cpu.ArchitectureArmV7:
		return campd.Architecture_ARM_V7
	case cpu.ArchitectureRiscV64:
		return campd.Architecture_RISCV_64
	case cpu.Architecturex86:
		return campd.Architecture_X86
	default:
		return campd.Architecture_ARCHITECTURE_UNKNOWN
	}
}

func CpuFromDomain(in *cpu.CPU) *campd.CpuSummary {
	processors := make([]*campd.ProcessorSummary, len(in.Processors))
	for i, p := range in.Processors {
		processors[i] = ProcessorFromDomain(p)
	}
	return &campd.CpuSummary{
		TotalCores:   int32(in.TotalCores),
		TotalThreads: int32(in.TotalThreads),
		Processors:   processors,
		Architecture: NewArchitecture(in.Architecture),
		Vendor:       grpcmodel.NewStringValue(in.Vendor),
		Model:        grpcmodel.NewStringValue(in.Model),
	}
}

func ProcessorFromDomain(in *cpu.Processor) *campd.ProcessorSummary {
	cores := make([]*campd.CoreSummary, len(in.Cores))
	for i, c := range in.Cores {
		cores[i] = CoreFromDomain(c)
	}

	return &campd.ProcessorSummary{
		Identifier:  int32(in.Id),
		CoreCount:   int32(in.CoreCount),
		ThreadCount: int32(in.ThreadCount),
		Vendor:      grpcmodel.NewStringValue(in.Vendor),
		Model:       grpcmodel.NewStringValue(in.Model),
		Cores:       cores,
	}
}

func CoreFromDomain(in *cpu.Core) *campd.CoreSummary {
	return &campd.CoreSummary{
		Identifier:  int32(in.Id),
		ThreadCount: int32(in.Threads),
	}
}
