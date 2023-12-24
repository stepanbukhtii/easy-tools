package easylog

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stepanbukhtii/easy-tools/config"
	"github.com/stepanbukhtii/easy-tools/easycontext"
	"os"
)

type Hook struct {
	serviceName string
	environment string
	version     string
}

func NewHook(c config.Application) *Hook {
	return &Hook{
		serviceName: c.ServiceName,
		environment: c.Environment,
		version:     c.Version,
	}
}

func (h *Hook) Run(e *zerolog.Event, level zerolog.Level, message string) {
	e.Str(ServiceEnvironment, h.environment)
	e.Str(ServiceName, h.serviceName)

	if h.version != "" {
		e.Str(ServiceVersion, h.version)
	}
}

func Init(logConfig config.Log, appConfig config.Application) {
	h := NewHook(appConfig)

	level, err := zerolog.ParseLevel(logConfig.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}

	log.Logger = zerolog.New(os.Stdout).Hook(h).With().Timestamp().Logger().Level(level)
}

func Debug(ctx context.Context) *zerolog.Event {
	return addLogsFields(ctx, log.Debug())
}

func Info(ctx context.Context) *zerolog.Event {
	return addLogsFields(ctx, log.Info())
}

func Warn(ctx context.Context) *zerolog.Event {
	return addLogsFields(ctx, log.Warn())
}

func Error(ctx context.Context) *zerolog.Event {
	return addLogsFields(ctx, log.Error())
}

func Fatal() *zerolog.Event {
	return log.Fatal()
}

func addLogsFields(ctx context.Context, e *zerolog.Event) *zerolog.Event {
	if traceID, ok := easycontext.GetTraceID(ctx); ok {
		e.Str(TraceID, traceID)
	}
	if fields, ok := easycontext.GetLogFields(ctx); ok {
		for key, value := range fields {
			e.Str(key, value)
		}
	}
	return e
}
