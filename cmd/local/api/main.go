package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	_ "embed"

	"github.com/awlsring/camp/internal/app/local/adapters/primary/rest/ogen"
	"github.com/awlsring/camp/internal/app/local/adapters/primary/rest/ogen/handler"
	machine_repository "github.com/awlsring/camp/internal/app/local/adapters/secondary/repository/postgres/machine"
	tag_repository "github.com/awlsring/camp/internal/app/local/adapters/secondary/repository/postgres/tag"
	"github.com/awlsring/camp/internal/app/local/config"
	"github.com/awlsring/camp/internal/app/local/core/service/machine"
	"github.com/awlsring/camp/internal/pkg/database"
	"github.com/awlsring/camp/internal/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const ServiceName = "CampLocal"

func main() {
	ctx := logger.InitContextLogger(context.Background(), zerolog.DebugLevel)
	log := logger.FromContext(ctx)
	log.Info().Msg("Initializing")

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.toml"
		log.Warn().Msgf("CONFIG_PATH not set, using default config path as %s", configPath)
	}
	cfg, err := config.LoadConfig(ctx, configPath)
	if err != nil {
		panic(errors.Wrap(err, "config"))
	}

	log.Info().Msg("Connecting to Postgres")
	pgDb, err := sqlx.Connect("postgres", database.CreatePostgresConnectionString(*cfg.Database.Postgres))
	if err != nil {
		panic(errors.Wrap(err, "postgres"))
	}
	defer pgDb.Close()

	tagRepo, err := tag_repository.New(pgDb)
	if err != nil {
		panic(errors.Wrap(err, "tag repo"))
	}

	log.Debug().Msg("Initializing Machine Repo")
	machineRepo, err := machine_repository.New(pgDb, tagRepo)
	if err != nil {
		panic(errors.Wrap(err, "machine repo"))
	}
	log.Debug().Msg("Initializing Machine Service")
	mSvc := machine.NewMachineService(machineRepo)

	log.Debug().Msg("Initializing Handler")
	hdl := handler.NewHandler(mSvc)

	log.Debug().Msg("Initializing Server")
	srv, err := ogen.NewCampLocalServer(hdl, ogen.Config{
		ServiceName:    ServiceName,
		MetricsAddress: cfg.Metrics.Address,
		ApiAddress:     cfg.Server.Address,
		ApiKeys:        cfg.Server.ApiKeys,
		AgentKeys:      cfg.Server.AgentKeys,
	})
	if err != nil {
		panic(errors.Wrap(err, "server init"))
	}

	log.Info().Msg("Starting Camp Local Server")
	srv.Start(ctx)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	// Start the HTTP server in a goroutine
	go func() {
		srv.Start(ctx)
	}()

	// Wait for an interrupt signal
	<-stopChan
	log.Info().Msg("Shutting down server")

	// Create a context with a timeout to give outstanding requests a chance to finish
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := srv.Stop(ctx); err != nil {
		panic(err)
	}
}
