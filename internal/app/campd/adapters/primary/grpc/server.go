package grpc

import (
	"context"
	"net"

	"github.com/awlsring/camp/internal/app/campd/adapters/primary/grpc/interceptor"

	"github.com/awlsring/camp/internal/pkg/logger"
	campd "github.com/awlsring/camp/pkg/gen/campd_grpc"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

const (
	DefaultNetwork = "tcp"
	DefaultAddress = ":8032"
)

type CampdGrpcServer struct {
	network  string
	address  string
	srv      *grpc.Server
	listener net.Listener
	metrics  *grpcprom.ServerMetrics
}

func (c *CampdGrpcServer) GetMetricsCollector() *grpcprom.ServerMetrics {
	return c.metrics
}

func NewServer(hdl campd.CampdServer, opts ...ServerOpt) (*CampdGrpcServer, error) {
	metrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerCounterOptions(),
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
	)

	grpcOpts := []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.NewLoggingInterceptor(zerolog.DebugLevel)),
		grpc.ChainUnaryInterceptor(
			metrics.UnaryServerInterceptor(),
		),
	}
	srv := grpc.NewServer(grpcOpts...)
	metrics.InitializeMetrics(srv)

	campd.RegisterCampdServer(srv, hdl)
	s := &CampdGrpcServer{
		network: DefaultNetwork,
		address: DefaultAddress,
		srv:     srv,
		metrics: metrics,
	}

	for _, opt := range opts {
		opt(s)
	}

	lis, err := net.Listen(s.network, s.address)
	if err != nil {
		return nil, err
	}
	s.listener = lis

	return s, nil
}

func (s *CampdGrpcServer) Start(ctx context.Context) error {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Starting server...")

	go func() {
		log.Debug().Msgf("server listening at %v", s.listener.Addr())
		if err := s.srv.Serve(s.listener); err != nil {
			log.Error().Err(err).Msg("Server error")
		}
	}()

	go func() {
		<-ctx.Done()
		log.Debug().Msg("Shutting down server...")
		s.srv.GracefulStop()
	}()

	return nil
}
