package model

import (
	"github.com/awlsring/camp/internal/pkg/iface"
	campd "github.com/awlsring/camp/pkg/gen/campd_grpc"
)

func NewTimestamp[T iface.Number](v T) *campd.Timestamp {
	return &campd.Timestamp{
		Value: int64(v),
	}
}
