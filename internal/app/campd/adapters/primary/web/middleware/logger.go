package middleware

import (
	"net/http"

	"github.com/awlsring/camp/internal/pkg/logger"
	"github.com/rs/zerolog"
)

func InitLogging(level zerolog.Level) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := logger.InitContextLogger(r.Context(), level)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
