package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	rabbitmq_power_control "github.com/awlsring/camp/apps/local/internal/adapters/primary/worker/rabbitmq/power_control"
	machine_repository "github.com/awlsring/camp/apps/local/internal/adapters/secondary/repository/postgres/machine"
	tag_repository "github.com/awlsring/camp/apps/local/internal/adapters/secondary/repository/postgres/tag"
	"github.com/awlsring/camp/apps/local/internal/config"
	"github.com/awlsring/camp/apps/local/internal/core/service/power_control"
	"github.com/awlsring/camp/internal/pkg/database"
	"github.com/awlsring/camp/internal/pkg/logger"
	"github.com/awlsring/camp/internal/pkg/rabbitmq_worker"
	"github.com/awlsring/camp/internal/pkg/rabbitmq_worker/exchange"
	"github.com/awlsring/camp/internal/pkg/rabbitmq_worker/job"
	"github.com/awlsring/camp/internal/pkg/rabbitmq_worker/queue"
	"github.com/awlsring/camp/internal/pkg/wol"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Job struct {
}

func (j *Job) Execute(ctx context.Context, body []byte) error {
	log := logger.FromContext(ctx)
	log.Info().Msgf("received message: %s", string(body))
	return nil
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	ctx = logger.InitContextLogger(ctx, zerolog.DebugLevel)

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

	conn, err := amqp.Dial("amqp://rabbit:rabbit@localhost:5672/")
	panicOnError(err)
	defer conn.Close()

	channel, err := conn.Channel()
	panicOnError(err)

	udpConn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: net.IPv4(255, 255, 255, 255), Port: 9})
	panicOnError(err)
	defer udpConn.Close()

	wakeOnLan := wol.New(udpConn)
	tagRepo, err := tag_repository.New(pgDb)
	panicOnError(err)
	machineRepo, err := machine_repository.New(pgDb, tagRepo)
	panicOnError(err)

	powerSvc := power_control.NewPowerControlService(wakeOnLan, machineRepo, nil)
	powerChangeJob := rabbitmq_power_control.NewPowerChangeRequestJob(powerSvc)

	exchange := exchange.NewDefinition("power_change", exchange.ExchangeTypeTopic)

	jobs := []*job.Definition{
		{
			Queue:          queue.NewDefinition("power_change_request", "power_change.request", "power_change"),
			Exchange:       exchange,
			ConcurrentJobs: 10,
			Job:            powerChangeJob,
		},
		{
			Queue:          queue.NewDefinition("power_change_validate", "power_change.validate", "power_change"),
			Exchange:       exchange,
			ConcurrentJobs: 10,
			Job:            powerChangeJob,
		},
	}

	manager := rabbitmq_worker.NewWorkerManager(channel, jobs)
	panicOnError(manager.Start(ctx))
}
