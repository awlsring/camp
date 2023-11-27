package main

import (
	"context"
	"flag"
	"net/http"
	"os"

	_ "embed"

	"github.com/awlsring/camp/apps/local/machine/controller"
	"github.com/awlsring/camp/apps/local/machine/repo"
	"github.com/awlsring/camp/apps/local/service"
	ogen_server "github.com/awlsring/camp/internal/pkg/server/ogen"
	camplocal "github.com/awlsring/camp/packages/camp_local"
	"github.com/jmoiron/sqlx"
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
	ogen_server.Run(func(ctx context.Context, lg *zap.Logger) error {
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

		m, err := ogen_server.NewMetrics(arg.MetricsAddr, "CampLocal")
		if err != nil {
			return errors.Wrap(err, "metrics")
		}

		dbHost := os.Getenv("DB_HOST")
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")

		dbConfig := repo.RepoConfig{
			Driver:   "postgres",
			Host:     dbHost,
			Port:     5432,
			Username: dbUser,
			Password: dbPassword,
			Database: "camplocal",
			UseSsl:   false,
		}

		pgDb, err := sqlx.Connect("postgres", repo.CreatePostgresConnectionString(dbConfig))
		if err != nil {
			return errors.Wrap(err, "postgres")
		}
		defer pgDb.Close()

		machineRepo, err := repo.NewPqRepo(pgDb)
		if err != nil {
			return errors.Wrap(err, "machine repo")
		}
		machineController := controller.NewController(machineRepo)

		srv, err := camplocal.NewServer(service.NewHandler(machineController),
			service.SecurityHandler("a", []string{"a"}),
			camplocal.WithTracerProvider(m.TracerProvider()),
			camplocal.WithMeterProvider(m.MeterProvider()),
			camplocal.WithErrorHandler(ogen_server.SmithyErrorHandler),
			camplocal.WithMiddleware(),
		)
		if err != nil {
			return errors.Wrap(err, "server init")
		}

		mux := http.NewServeMux()
		mux.Handle("/", srv)
		mux.Handle("/swagger/", ogen_server.SwaggerUIHandler())
		mux.Handle("/swagger/swagger.json", ogen_server.SwaggerAPIv1Handler(doc))
		httpServer := http.Server{
			Addr:    arg.Addr,
			Handler: mux,
		}
		g, ctx := errgroup.WithContext(ctx)
		g.Go(func() error {
			return m.Start(ctx)
		})
		g.Go(func() error {
			<-ctx.Done()
			return httpServer.Shutdown(ctx)
		})
		g.Go(func() error {
			defer lg.Info("Server stopped")
			if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				return errors.Wrap(err, "http")
			}
			return nil
		})

		return g.Wait()
	})
}
