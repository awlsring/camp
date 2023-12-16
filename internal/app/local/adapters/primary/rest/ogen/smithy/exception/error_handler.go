package exception

import (
	"context"
	"net/http"

	"github.com/awlsring/camp/internal/pkg/logger"
	"github.com/rs/zerolog"
)

const (
	SmithyErrorTypeHeader = "X-Amzn-Errortype"
)

func ResponseHandlerWithLogger(lvl zerolog.Level) func(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
		ctx = logger.InitContextLogger(ctx, lvl)
		ResponseHandler(ctx, w, r, err)
	}
}

func ResponseHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) { //
	log := logger.FromContext(ctx)
	log.Debug().Err(err).Msg("Handling error as smithy exception")

	log.Debug().Msg("Getting exception from error")
	exception := FromError(err)
	log.Debug().Interface("exception", exception).Msgf("got exception")

	log.Debug().Msg("Building response")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set(SmithyErrorTypeHeader, exception.Type().String())
	w.WriteHeader(exception.Code())
	w.Write(exception.AsJsonMessage())
}
