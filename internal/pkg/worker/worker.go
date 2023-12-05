package worker

import "context"

type Worker interface {
	Name() string
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
