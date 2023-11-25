package easylog

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/stepanbukhtii/easy-tools/config"
	"github.com/stepanbukhtii/easy-tools/easycontext"
)

type Hook struct {
	serviceName string
	environment string
	version     string
}

func NewHook(c config.Service) *Hook {
	return &Hook{
		serviceName: c.Name,
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

func Init(logConfig config.Log, serviceConfig config.Service) {
	h := NewHook(serviceConfig)

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
	if traceID := easycontext.TraceID(ctx); traceID != "" {
		e.Str(TraceID, traceID)
	}
	if fields := easycontext.LogFields(ctx); fields != nil {
		for key, value := range fields {
			e.Str(key, value)
		}
	}
	return e
}
