package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-faster/jx"
	ht "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/ogen-go/ogen/validate"
)

func ErrorCode(err error) (code int) {
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
		code = ogenErr.Code()
	}

	return code
}

func SmithyErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	code := ErrorCode(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	e := jx.GetEncoder()
	e.ObjStart()
	e.FieldStart("message")
	e.StrEscape(err.Error())
	e.ObjEnd()

	_, _ = w.Write(e.Bytes())
}
