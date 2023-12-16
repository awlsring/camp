package grpc

import (
	"context"
	"net"

	"github.com/awlsring/camp/internal/pkg/logger"
	campd "github.com/awlsring/camp/pkg/gen/campd_grpc"
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
}

func NewServer(hdl campd.CampdServer, opts ...ServerOpt) (*CampdGrpcServer, error) {
	srv := grpc.NewServer()
	campd.RegisterCampdServer(srv, hdl)
	s := &CampdGrpcServer{
		network: DefaultNetwork,
		address: DefaultAddress,
		srv:     srv,
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
