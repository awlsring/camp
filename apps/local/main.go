package main

import (
	"context"
	"flag"
	"net/http"

	_ "embed"

	camplocal "github.com/awlsring/camp/generated/camp_local"
	camplocalapi "github.com/awlsring/camp/internal/app/local"
	"github.com/awlsring/camp/internal/pkg/server"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

//go:embed swagger/swagger.json
var doc []byte

func main() {
	server.Run(func(ctx context.Context, lg *zap.Logger) error {
		var arg struct {
			Addr        string
			MetricsAddr string
		}
		flag.StringVar(&arg.Addr, "addr", "127.0.0.1:8080", "listen address")
		flag.StringVar(&arg.MetricsAddr, "metrics.addr", "127.0.0.1:9090", "metrics listen address")
		flag.Parse()

		lg.Info("Initializing",
			zap.String("http.addr", arg.Addr),
			zap.String("metrics.addr", arg.MetricsAddr),
		)

		m, err := server.NewMetrics(lg, server.Config{
			Addr: arg.MetricsAddr,
			Name: "api",
		})
		if err != nil {
			return errors.Wrap(err, "metrics")
		}

		srv, err := camplocal.NewServer(camplocalapi.NewHandler(),
			camplocalapi.SecurityHandler(server.NewApiKeyAuth()),
			camplocal.WithTracerProvider(m.TracerProvider()),
			camplocal.WithMeterProvider(m.MeterProvider()),
			camplocal.WithMiddleware(),
		)
		if err != nil {
			return errors.Wrap(err, "server init")
		}
		httpServer := http.Server{
			Addr:    arg.Addr,
			Handler: srv,
		}

		mux := http.NewServeMux()
		mux.Handle("/", http.StripPrefix("/api/v1", srv))
		mux.Handle("/swagger/", server.SwaggerUIHandler())
		mux.Handle("/swagger/swagger.json", server.SwaggerAPIv1Handler(doc))
		g, ctx := errgroup.WithContext(ctx)
		g.Go(func() error {
			return m.Run(ctx)
		})
		g.Go(func() error {
			<-ctx.Done()
			return httpServer.Shutdown(ctx)
		})
		g.Go(func() error {
			defer lg.Info("Server stopped")
			// if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// 	return errors.Wrap(err, "http")
			// }
			if err := http.ListenAndServe(arg.Addr, mux); err != nil && err != http.ErrServerClosed {
				return errors.Wrap(err, "http")
			}
			return nil
		})

		return g.Wait()
	})
}
