package web

import (
	"context"
	"net/http"

	"github.com/awlsring/camp/internal/app/campd/adapters/primary/web/components"
	"github.com/awlsring/camp/internal/app/campd/adapters/primary/web/handler"
	"github.com/awlsring/camp/internal/pkg/logger"
)

type Server struct {
	address string
	hdl     *handler.SystemHandler
}

func NewServer(hdl *handler.SystemHandler, opts ...ServerOpt) *Server {
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sum, err := s.hdl.GetSystem(ctx)
		if err != nil {
			log.Error().Err(err).Msg("error getting system information")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Debug().Msg("Handling request")
		components.Summary("m-123", sum).Render(ctx, w)
	})

	go func() {
		log.Debug().Msgf("web server listening at %v", s.address)
		if err := http.ListenAndServe(s.address, nil); err != nil {
			log.Error().Err(err).Msg("error listening")
		}
	}()

	go func() {
		<-ctx.Done()
		log.Debug().Msg("Shutting down web server...")
	}()

	return nil
}
