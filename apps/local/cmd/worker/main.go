package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	power_state_job "github.com/awlsring/camp/apps/local/internal/adapters/primary/worker/rabbitmq/power_state"
	machine_repository "github.com/awlsring/camp/apps/local/internal/adapters/secondary/repository/postgres/machine"
	tag_repository "github.com/awlsring/camp/apps/local/internal/adapters/secondary/repository/postgres/tag"
	"github.com/awlsring/camp/apps/local/internal/adapters/secondary/topic/mqtt/power_state_change"
	"github.com/awlsring/camp/apps/local/internal/config"
	"github.com/awlsring/camp/apps/local/internal/core/service/power_state"
	"github.com/awlsring/camp/internal/pkg/campd"
	"github.com/awlsring/camp/internal/pkg/database"
	"github.com/awlsring/camp/internal/pkg/logger"
	"github.com/awlsring/camp/internal/pkg/wol"
	"github.com/awlsring/camp/internal/pkg/worker"
	"github.com/awlsring/camp/internal/pkg/worker/rabbitmq"
	"github.com/awlsring/camp/internal/pkg/worker/rabbitmq/exchange"
	"github.com/awlsring/camp/internal/pkg/worker/rabbitmq/queue"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

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

	mqttOpts := mqtt.NewClientOptions()
	mqttOpts.AddBroker("tcp://localhost:1883")
	mqttOpts.SetClientID("worker")

	log.Debug().Msg("connecting to mqtt broker")
	client := mqtt.NewClient(mqttOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Error().Err(token.Error()).Msg("failed to connect to mqtt broker")
		panic(token.Error())
	}
	defer client.Disconnect(30)
	stateChangeTopic := power_state_change.New(client, "camp/local/power/state_change")

	campd := campd.NewCacheClient()

	powerSvc := power_state.NewService(machineRepo, campd, wakeOnLan, stateChangeTopic)
	requestChangeJob := power_state_job.NewRequestStateChangeJob(powerSvc)

	validateChangeJob := power_state_job.NewValidateStateChangeJob(powerSvc)

	exchange := exchange.NewDefinition("power_change", exchange.ExchangeTypeTopic)

	powerChangeJobDef := &rabbitmq.JobDefinition{
		Queue:          queue.NewDefinition("power_change_request", "power_change.request", "power_change"),
		Exchange:       exchange,
		ConcurrentJobs: 10,
		Job:            requestChangeJob,
	}
	requestWorker := rabbitmq.NewWorker(channel, powerChangeJobDef)

	validationJobDef := &rabbitmq.JobDefinition{
		Queue:          queue.NewDefinition("power_change_validate", "power_change.validate", "power_change"),
		Exchange:       exchange,
		ConcurrentJobs: 10,
		Job:            validateChangeJob,
	}
	validaitonWorker := rabbitmq.NewWorker(channel, validationJobDef)

	manager := worker.NewWorkerManager([]worker.Worker{requestWorker, validaitonWorker})
	panicOnError(manager.Start(ctx))
}
