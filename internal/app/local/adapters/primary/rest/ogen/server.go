package ogen

import (
	"context"
	"net/http"
	"time"

	"github.com/awlsring/camp/internal/app/local/adapters/primary/rest/ogen/auth"
	"github.com/awlsring/camp/internal/app/local/adapters/primary/rest/ogen/smithy/exception"
	"github.com/awlsring/camp/internal/pkg/logger"
	ogen_server "github.com/awlsring/camp/internal/pkg/server/ogen"
	"github.com/awlsring/camp/internal/pkg/server/ogen/middleware"
	camplocal "github.com/awlsring/camp/pkg/gen/local"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Config struct {
	ServiceName    string
	MetricsAddress string
	ApiAddress     string
	ApiKeys        []string
	AgentKeys      []string
}

type CampLocalServer struct {
	Api     *ogen_server.Rest
	Metrics *ogen_server.Metrics
}

func NewCampLocalServer(handler camplocal.Handler, cfg Config) (*CampLocalServer, error) {
	m, err := ogen_server.NewMetrics(cfg.MetricsAddress, cfg.ServiceName)
	if err != nil {
		return nil, errors.Wrap(err, "metrics")
	}

	srv, err := camplocal.NewServer(handler,
		auth.NewSecurityHandler(cfg.ApiKeys, cfg.AgentKeys),
		camplocal.WithTracerProvider(m.TracerProvider()),
		camplocal.WithMeterProvider(m.MeterProvider()),
		camplocal.WithNotFound(exception.UnknownOperationHandler),
		camplocal.WithErrorHandler(exception.ResponseHandlerWithLogger(zerolog.DebugLevel)),
		camplocal.WithMiddleware(middleware.LoggingMiddleware(zerolog.DebugLevel)),
	)
	if err != nil {
		return nil, errors.Wrap(err, "server init")
	}

	s := &http.Server{
		Addr:    cfg.ApiAddress,
		Handler: srv,
	}
	r := ogen_server.NewRest(s)

	return &CampLocalServer{
		Api:     r,
		Metrics: m,
	}, nil
}

func (s *CampLocalServer) Start(ctx context.Context) error {
	log := logger.FromContext(ctx)
	log.Info().Msg("Starting Camp Local Server")
	go func() {
		log.Info().Msg("Starting Api server")
		s.Api.Start(ctx)
	}()
	go func() {
		log.Info().Msg("Starting Metrics server")
		s.Metrics.Start(ctx)
	}()

	// Wait for the context to be, then signal stop to both servers.
	<-ctx.Done()
	log.Info().Msg("Server context done, stopping servers")

	return s.Stop(context.Background())
}

func (s *CampLocalServer) Stop(ctx context.Context) error {
	log := logger.FromContext(ctx)

	log.Info().Msg("Stopping Camp Local Server")
	context, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	go func() {
		log.Info().Msg("Sending stop to Api server")
		s.Api.Stop(context)
	}()
	go func() {
		log.Info().Msg("Sending stop to Metrics server")
		s.Metrics.Stop(context)
	}()

	<-context.Done()

	log.Info().Msg("Server stopped")
	return nil
}
