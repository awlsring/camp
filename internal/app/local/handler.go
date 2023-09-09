package camplocalapi

import (
	"context"

	camplocal "github.com/awlsring/camp/generated/camp_local"
)

// Compile-time check for Handler.
var _ camplocal.Handler = (*Handler)(nil)

type Handler struct{}

func NewHandler() camplocal.Handler {
	return Handler{}
}

func (h Handler) Health(ctx context.Context) (camplocal.HealthRes, error) {
	return &camplocal.HealthResponseContent{
		Success: true,
	}, nil
}
