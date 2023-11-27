package ogen

import (
	"github.com/awlsring/camp/apps/local/service"
	ogen_server "github.com/awlsring/camp/internal/pkg/server/ogen"
	camplocal "github.com/awlsring/camp/packages/camp_local"
	"github.com/pkg/errors"
)

type Config struct {
	ServiceName    string
	MetricsAddress string
}

type CampLocalServer struct {
	MetricsServer *ogen_server.Metrics
	Server        *camplocal.Server
}

func NewCampLocalServer(handler camplocal.Handler, cfg Config) (*CampLocalServer, error) {
	m, err := ogen_server.NewMetrics(cfg.MetricsAddress, cfg.ServiceName)
	if err != nil {
		return nil, errors.Wrap(err, "metrics")
	}

	srv, err := camplocal.NewServer(handler,
		service.SecurityHandler("a", []string{"a"}),
		camplocal.WithTracerProvider(m.TracerProvider()),
		camplocal.WithMeterProvider(m.MeterProvider()),
		camplocal.WithErrorHandler(ogen_server.SmithyErrorHandler),
		camplocal.WithMiddleware(),
	)
	if err != nil {
		return nil, errors.Wrap(err, "server init")
	}

	return &CampLocalServer{
		MetricsServer: m,
		Server:        srv,
	}, nil
}

func (s *CampLocalServer) Run() error {
	return nil
}
