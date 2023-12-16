package cpu

import "github.com/awlsring/camp/internal/pkg/values"

type Processor struct {
	Id          int
	CoreCount   uint32
	ThreadCount uint32
	Model       *string
	Vendor      *string
	Cores       []*Core
}

func NewProcessor(id int, coreCount uint32, threadCount uint32, vendor, model string, cores []*Core) *Processor {
	return &Processor{
		Id:          id,
		CoreCount:   coreCount,
		ThreadCount: threadCount,
		Model:       values.ParseOptional(model),
		Vendor:      values.ParseOptional(vendor),
	}
}
