package middleware

import (
	"github.com/awlsring/camp/internal/pkg/logger"
	"github.com/ogen-go/ogen/middleware"
	"github.com/rs/zerolog"
)

func LoggingMiddleware(lvl zerolog.Level) middleware.Middleware {
	return func(
		req middleware.Request,
		next func(req middleware.Request) (middleware.Response, error),
	) (middleware.Response, error) {
		req.Context = logger.InitContextLogger(req.Context, lvl)

		log := logger.FromContext(req.Context)

		log.Info().Msg("Handling request")
		resp, err := next(req)
		if err != nil {
			log.Error().Err(err)
		}
		return resp, err
	}
}
