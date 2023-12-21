package metrics

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	DefaultAddress = ":9090"
)

type SeverOpt func(*Server)

func WithCollectors(collectors ...prometheus.Collector) SeverOpt {
	return func(s *Server) {
		s.collectors = append(s.collectors, collectors...)
	}
}

func WithHandlerOpts(opts promhttp.HandlerOpts) SeverOpt {
	return func(s *Server) {
		s.handlerOpts = opts
	}
}

func WithAddress(address string) SeverOpt {
	return func(s *Server) {
		s.address = address
	}
}

func WithRegistry(registry *prometheus.Registry) SeverOpt {
	return func(s *Server) {
		s.registry = registry
	}
}

type Server struct {
	address     string
	collectors  []prometheus.Collector
	registry    *prometheus.Registry
	handlerOpts promhttp.HandlerOpts
}

func NewServer(opts ...SeverOpt) *Server {
	s := &Server{
		address:     DefaultAddress,
		collectors:  []prometheus.Collector{},
		handlerOpts: promhttp.HandlerOpts{},
		registry:    prometheus.NewRegistry(),
	}
	for _, opt := range opts {
		opt(s)
	}

	s.registry.MustRegister(s.collectors...)

	return s
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(
		s.registry,
		s.handlerOpts,
	))

	srv := &http.Server{
		Addr:    s.address,
		Handler: mux,
	}

	go func() {
		srv.ListenAndServe()
	}()

	go func() {
		<-ctx.Done()
		srv.Shutdown(ctx)
	}()

	return nil
}
