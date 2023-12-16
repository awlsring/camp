package cpu

import "github.com/awlsring/camp/internal/pkg/values"

type CPU struct {
	TotalCores   uint32
	TotalThreads uint32
	Architecture Architecture
	Vendor       *string
	Model        *string
	Processors   []*Processor
}

func NewCPU(totalCores uint32, totalThreads uint32, architecture Architecture, vendor string, model string, processors []*Processor) *CPU {
	return &CPU{
		TotalCores:   totalCores,
		TotalThreads: totalThreads,
		Architecture: architecture,
		Vendor:       values.ParseOptional(vendor),
		Model:        values.ParseOptional(model),
		Processors:   processors,
	}
}
