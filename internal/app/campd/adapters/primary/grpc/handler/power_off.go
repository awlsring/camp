package handler

import (
	"context"

	campd "github.com/awlsring/camp/pkg/gen/campd_grpc"
)

func (h *Handler) PowerOff(context.Context, *campd.PowerOffInput) (*campd.PowerOffOutput, error) {
	panic("implement me")
}
