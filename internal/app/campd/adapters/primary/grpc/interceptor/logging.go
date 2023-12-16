package interceptor

import (
	"context"

	"github.com/awlsring/camp/internal/pkg/logger"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

func NewLoggingInterceptor(level zerolog.Level) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = logger.InitContextLogger(ctx, level)
		return handler(ctx, req)
	}
}
