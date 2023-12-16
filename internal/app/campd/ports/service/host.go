package service

import (
	"context"

	"github.com/awlsring/camp/internal/pkg/domain/host"
)

type Host interface {
	Describe(ctx context.Context) (*host.Host, error)
	Uptime(ctx context.Context) (uint64, error)
	BootTime(ctx context.Context) (uint64, error)
}
