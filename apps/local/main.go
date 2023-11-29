package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	_ "embed"

	"github.com/awlsring/camp/apps/local/internal/adapters/primary/rest/ogen"
	"github.com/awlsring/camp/apps/local/internal/adapters/primary/rest/ogen/handler"
	"github.com/awlsring/camp/apps/local/internal/adapters/secondary/repository/machine/postgres"
	"github.com/awlsring/camp/apps/local/internal/core/service/machine"
	"github.com/awlsring/camp/internal/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

//go:embed swagger/swagger.json
var doc []byte

func main() {
	ctx := logger.InitContextLogger(context.Background(), zerolog.DebugLevel)
	log := logger.FromContext(ctx)
	log.Info().Msg("Initializing")

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	dbConfig := postgres.RepoConfig{
		Driver:   "postgres",
		Host:     dbHost,
		Port:     5432,
		Username: dbUser,
		Password: dbPassword,
		Database: "camplocal",
		UseSsl:   false,
	}

	log.Info().Msg("Connecting to Postgres")
	pgDb, err := sqlx.Connect("postgres", postgres.CreatePostgresConnectionString(dbConfig))
	if err != nil {
		panic(errors.Wrap(err, "postgres"))
	}
	defer pgDb.Close()

	log.Debug().Msg("Initializing Machine Repo")
	machineRepo, err := postgres.NewPqRepo(pgDb)
	if err != nil {
		panic(errors.Wrap(err, "machine repo"))
	}
	log.Debug().Msg("Initializing Machine Service")
	mSvc := machine.NewMachineService(machineRepo)

	log.Debug().Msg("Initializing Handler")
	hdl := handler.NewHandler(mSvc)

	log.Debug().Msg("Initializing Server")
	srv, err := ogen.NewCampLocalServer(hdl, ogen.Config{
		ServiceName:    "CampLocal",
		MetricsAddress: "127.0.0.1:8032",
		ApiAddress:     "127.0.0.1:7032",
		ApiKeys:        []string{"a"},
		AgentKeys:      []string{"a"},
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
