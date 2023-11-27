package logger

import (
	"context"
	"os"
	"runtime/debug"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var DefaultLogger zerolog.Logger = log.With().Logger()

type LoggingOpt func(zctx zerolog.Context, ctx context.Context) zerolog.Context

func InitContextLogger(ctx context.Context, lvl zerolog.Level, opts ...LoggingOpt) context.Context {
	buildInfo, _ := debug.ReadBuildInfo()

	logger := FromContext(ctx)
	logger = logger.Level(lvl)
	logger = logger.With().
		Caller().
		Timestamp().
		Caller().
		Int("pid", os.Getpid()).
		Str("go_version", buildInfo.GoVersion).
		Logger()

	for _, opt := range opts {
		logger.UpdateContext(func(zctx zerolog.Context) zerolog.Context {
			return opt(zctx, ctx)
		})
	}

	ctx = logger.WithContext(ctx)

	return ctx
}

func FromContext(ctx context.Context) zerolog.Logger {
	logger := log.Ctx(ctx)
	if logger == nil {
		return DefaultLogger
	}
	return *logger
}

func LogLevelFromEnv() zerolog.Level {
	lvlString := os.Getenv("LOG_LEVEL")
	if lvlString == "" {
		return zerolog.InfoLevel
	}
	lvl, err := zerolog.ParseLevel(lvlString)
	if err != nil {
		return zerolog.InfoLevel
	}
	return lvl
}
