package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	machine_repository "github.com/awlsring/camp/apps/local/internal/adapters/secondary/repository/postgres/machine"
	tag_repository "github.com/awlsring/camp/apps/local/internal/adapters/secondary/repository/postgres/tag"
	"github.com/awlsring/camp/apps/local/internal/adapters/secondary/topic/mqtt/power_state_change"
	"github.com/awlsring/camp/apps/local/internal/config"
	"github.com/awlsring/camp/apps/local/internal/core/service/state_monitor"
	"github.com/awlsring/camp/internal/pkg/agent"
	"github.com/awlsring/camp/internal/pkg/database"
	"github.com/awlsring/camp/internal/pkg/logger"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const TopicName = "camp/local/power/state_change"

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

	mqttOpts := mqtt.NewClientOptions()
	mqttOpts.AddBroker("tcp://localhost:1883")
	mqttOpts.SetClientID("camp-local-test-send")

	client := mqtt.NewClient(mqttOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	defer client.Disconnect(250)

	topic := power_state_change.New(client, TopicName)

	campdClient, err := agent.NewAgentClientCache()
	if err != nil {
		panic(errors.Wrap(err, "campd client"))
	}

	log.Debug().Msg("Initializing Machine Service")
	svc := state_monitor.NewService(machineRepo, topic, campdClient)

	// every 60 seconds, check the state of all machines
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			err := svc.VerifyAndAdjustMachineStates(ctx)
			if err != nil {
				log.Error().Err(err).Msg("failed to verify and adjust machine states")
			}
			select {
			case <-ctx.Done():
				log.Debug().Msg("context done, exiting")
				return
			case <-time.After(60 * time.Second):
				log.Debug().Msg("sleeping for 60 seconds")
			}
		}
	}()

	// wait for shutdown signal
	<-ctx.Done()
	log.Info().Msg("Shutting down")
	wg.Wait()
	log.Info().Msg("Shutdown complete")
}
