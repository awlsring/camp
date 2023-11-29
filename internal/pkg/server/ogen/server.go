package ogen_server

import (
	"context"
	"time"
)

type OgenServer struct {
	Rest    *Rest
	Metrics *Metrics
}

type OgenServerConfig struct {
	Address         string
	MetricsAddress  string
	ApplicationName string
}

func NewOgenServer(r *Rest, m *Metrics) *OgenServer {
	return &OgenServer{
		Rest:    r,
		Metrics: m,
	}
}

func (s *OgenServer) Start(ctx context.Context) error {
	go func() {
		s.Rest.Start(ctx)
	}()
	go func() {
		s.Metrics.Start(ctx)
	}()

	// Wait for the context to be, then signal stop to both servers.
	<-ctx.Done()

	return s.Stop(context.Background())
}

func (s *OgenServer) Stop(ctx context.Context) error {
	context, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	go func() {
		s.Rest.Stop(context)
	}()
	go func() {
		s.Metrics.Stop(context)
	}()

	<-context.Done()

	return nil
}
