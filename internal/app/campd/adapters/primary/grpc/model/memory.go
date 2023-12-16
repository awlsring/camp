package model

import (
	"github.com/awlsring/camp/internal/pkg/domain/memory"
	campd "github.com/awlsring/camp/pkg/gen/campd_grpc"
)

func MemoryFromDomain(in *memory.Memory) *campd.MemorySummary {
	return &campd.MemorySummary{
		Total: int64(in.Total),
	}
}
