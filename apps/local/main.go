package main

import (
	"context"
	"flag"
	"net/http"
	"os"

	_ "embed"

	"github.com/awlsring/camp/apps/local/machine"
	"github.com/awlsring/camp/apps/local/service"
	"github.com/awlsring/camp/internal/pkg/server"
	camplocal "github.com/awlsring/camp/packages/camp_local"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

//go:embed swagger/swagger.json
var doc []byte

func main() {
	server.Run(func(ctx context.Context, lg *zap.Logger) error {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		if err := godotenv.Load(); err != nil {
			log.Fatal().Err(err).Msg("Error loading .env file")
		}
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

		dbHost := os.Getenv("DB_HOST")
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")

		dbConfig := machine.RepoConfig{
			Driver:   "postgres",
			Host:     dbHost,
			Port:     5432,
			Username: dbUser,
			Password: dbPassword,
			Database: "camplocal",
			UseSsl:   false,
		}

		machineRepo := machine.NewRepo(dbConfig)
		machineController := machine.NewController(machineRepo)

		srv, err := camplocal.NewServer(service.NewHandler(machineController),
			service.SecurityHandler("a", []string{"a"}),
			camplocal.WithTracerProvider(m.TracerProvider()),
			camplocal.WithMeterProvider(m.MeterProvider()),
			camplocal.WithErrorHandler(server.SmithyErrorHandler),
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
		mux.Handle("/", srv)
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
