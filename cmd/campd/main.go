package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/awlsring/camp/internal/app/campd/adapters/primary/grpc"
	"github.com/awlsring/camp/internal/app/campd/adapters/primary/grpc/handler"
	"github.com/awlsring/camp/internal/app/campd/core/service/board"
	"github.com/awlsring/camp/internal/app/campd/core/service/cpu"
	"github.com/awlsring/camp/internal/app/campd/core/service/host"
	"github.com/awlsring/camp/internal/app/campd/core/service/memory"
	"github.com/awlsring/camp/internal/app/campd/core/service/network"
	"github.com/awlsring/camp/internal/app/campd/core/service/storage"
	"github.com/awlsring/camp/internal/pkg/logger"
	"github.com/rs/zerolog"
)

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = logger.InitContextLogger(ctx, zerolog.DebugLevel)
	log := logger.FromContext(ctx)
	log.Info().Msg("Initializing")

	log.Info().Msg("Initializing gRPC server")
	cpuSvc, err := cpu.InitService(ctx)
	panicOnErr(err)

	hostSvc, err := host.InitService(ctx)
	panicOnErr(err)

	memSvc, err := memory.InitService(ctx)
	panicOnErr(err)

	moboSvc, err := board.InitService(ctx)
	panicOnErr(err)

	netSvc, err := network.InitService(ctx)
	panicOnErr(err)

	storageSvc, err := storage.InitService(ctx)
	panicOnErr(err)

	hdl := handler.New(cpuSvc, hostSvc, memSvc, moboSvc, netSvc, storageSvc)
	srv, err := grpc.NewServer(hdl)
	panicOnErr(err)

	log.Info().Msg("Starting server")
	err = srv.Start(ctx)
	panicOnErr(err)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	log.Info().Msg("Shutting down server")
	cancel()

	log.Info().Msg("Waiting for server to shutdown")
	<-ctx.Done()
}
