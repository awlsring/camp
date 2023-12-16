package handler

import (
	"context"

	campd "github.com/awlsring/camp/pkg/gen/campd_grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) Reboot(context.Context, *emptypb.Empty) (*campd.RebootOutput, error) {
	panic("implement me")
}
