package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/awlsring/camp/internal/app/local/adapters/primary/schedule"
	machine_repository "github.com/awlsring/camp/internal/app/local/adapters/secondary/repository/postgres/machine"
	tag_repository "github.com/awlsring/camp/internal/app/local/adapters/secondary/repository/postgres/tag"
	power_state_topic "github.com/awlsring/camp/internal/app/local/adapters/secondary/topic/rabbitmq/power_state"
	"github.com/awlsring/camp/internal/app/local/config"
	"github.com/awlsring/camp/internal/app/local/core/service/state_monitor"
	"github.com/awlsring/camp/internal/pkg/database"
	"github.com/awlsring/camp/internal/pkg/logger"
	"github.com/jmoiron/sqlx"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	ctx = logger.InitContextLogger(ctx, zerolog.DebugLevel)
	defer cancel()

	log := logger.FromContext(ctx)
	log.Info().Msg("Initializing")

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.toml"
		log.Warn().Msgf("CONFIG_PATH not set, using default config path as %s", configPath)
	}
	cfg, err := config.LoadConfig(ctx, configPath)
	panicOnError(err)

	log.Info().Msg("Connecting to Postgres")
	pgDb, err := sqlx.Connect("postgres", database.CreatePostgresConnectionString(*cfg.Database.Postgres))
	panicOnError(err)
	defer pgDb.Close()

	tagRepo, err := tag_repository.New(pgDb)
	panicOnError(err)

	log.Debug().Msg("Initializing Machine Repo")
	machineRepo, err := machine_repository.New(pgDb, tagRepo)
	panicOnError(err)

	conn, err := amqp.Dial("amqp://rabbit:rabbit@localhost:5672/")
	panicOnError(err)
	defer conn.Close()

	channel, err := conn.Channel()
	panicOnError(err)

	topic := power_state_topic.New(channel)

	log.Debug().Msg("Initializing Machine Service")
	svc := state_monitor.NewService(machineRepo, topic)

	scheduler := schedule.NewScheduler(svc, schedule.WithInterval(60*time.Second))
	scheduler.Start(ctx)
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
