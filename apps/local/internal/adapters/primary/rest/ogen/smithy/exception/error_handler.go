package exception

import (
	"context"
	"errors"
	"net/http"

	"github.com/awlsring/camp/internal/pkg/logger"
	ht "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/ogen-go/ogen/validate"
)

const (
	SmithyErrorTypeHeader = "X-Amzn-Errortype"
)

func ErrorCode(ctx context.Context, err error) (code int) {
	log := logger.FromContext(ctx)
	code = http.StatusInternalServerError

	var (
		ctError *validate.InvalidContentTypeError
		ogenErr ogenerrors.Error
	)
	switch {
	case errors.Is(err, ht.ErrNotImplemented):
		code = http.StatusNotImplemented
	case errors.As(err, &ctError):
		code = http.StatusUnsupportedMediaType
	case errors.As(err, &ogenErr):
		log.Info().Msgf("op: %s", ogenErr.OperationName())
		log.Info().Msgf("opid: %s", ogenErr.OperationID())
		log.Info().Msgf("err: %s", ogenErr.Error())
		log.Info().Msgf("code: %d", ogenErr.Code())
		code = ogenErr.Code()
	}

	return code
}

func ResponseHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) { //
	log := logger.FromContext(ctx)
	log.Debug().Err(err).Msg("Handling error as smithy exception")

	log.Debug().Msg("Getting exception from error")
	exception := FromError(err)

	log.Debug().Msg("Building response")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set(SmithyErrorTypeHeader, exception.Type().String())
	w.WriteHeader(exception.Code())
	w.Write(exception.AsJsonMessage())
}
