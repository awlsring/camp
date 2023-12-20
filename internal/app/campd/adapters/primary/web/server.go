package web

import (
	"context"
	"net/http"
	"time"

	"github.com/awlsring/camp/internal/app/campd/adapters/primary/web/handler/system"
	"github.com/awlsring/camp/internal/app/campd/adapters/primary/web/middleware"
	"github.com/awlsring/camp/internal/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type Server struct {
	address string
	hdl     *system.Handler
}

func NewServer(hdl *system.Handler, opts ...ServerOpt) *Server {
	s := &Server{
		address: ":8080",
		hdl:     hdl,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) Start(ctx context.Context) error {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Loading webserver routes")

	router := chi.NewRouter()
	router.Use(middleware.InitLogging(zerolog.DebugLevel))

	MountAssets(router)
	s.hdl.Mount(router)

	server := &http.Server{
		Addr:    s.address,
		Handler: http.TimeoutHandler(router, 30*time.Second, "request timed out"),
	}

	go func() {
		log.Debug().Msgf("web server listening at %v", s.address)
		if err := server.ListenAndServe(); err != nil {
			log.Error().Err(err).Msg("error listening")
		}
	}()

	go func() {
		<-ctx.Done()
		log.Debug().Msg("Shutting down web server...")
	}()

	return nil
}
