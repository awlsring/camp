package service

import (
	"context"

	"github.com/awlsring/camp/internal/pkg/domain/power"
)

type Reporter interface {
	Heartbeat(ctx context.Context) error
	Register(ctx context.Context) error
	ReportStatus(ctx context.Context, status power.StatusCode) error
	ReportSystemInformation(ctx context.Context) error
}
