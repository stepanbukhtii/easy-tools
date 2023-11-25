package elog

import (
	"errors"
	"log/slog"
	"os"

	"github.com/stepanbukhtii/easy-tools/config"
	"github.com/stepanbukhtii/easy-tools/errx"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
)

func NewSlogHandler(logConfig config.Log, serviceConfig config.Service) slog.Handler {
	return NewCustomHandler(baseHandler(logConfig, serviceConfig), nil)
}

func NewSlogHandlerCustom(logConfig config.Log, serviceConfig config.Service, customHandler CustomHandler) slog.Handler {
	return NewCustomHandler(baseHandler(logConfig, serviceConfig), customHandler)
}

func Err(err error) slog.Attr {
	if err == nil {
		return slog.Attr{}
	}
	return slog.Any(string(semconv.ErrorTypeKey), err)
}

func baseHandler(logConfig config.Log, serviceConfig config.Service) slog.Handler {
	var level slog.Level
	if err := level.UnmarshalText([]byte(logConfig.Level)); err != nil {
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:       level,
		ReplaceAttr: slogReplaceAttr,
	}

	defaultAttrs := []slog.Attr{
		slog.String(string(semconv.ServiceNameKey), serviceConfig.Name),
		slog.String(string(semconv.DeploymentEnvironmentNameKey), serviceConfig.Environment),
	}

	if serviceConfig.Version != "" {
		defaultAttrs = append(defaultAttrs, slog.String(string(semconv.ServiceVersionKey), serviceConfig.Version))
	}

	return slog.NewJSONHandler(os.Stdout, opts).WithAttrs(defaultAttrs)
}

func slogReplaceAttr(_ []string, a slog.Attr) slog.Attr {
	if a.Key == slog.MessageKey {
		a.Key = "message"
	}

	if a.Key == string(semconv.ErrorTypeKey) {
		if err, ok := a.Value.Any().(error); ok {
			var errorX errx.Error
			if errors.As(err, &errorX) && errorX.LogData != nil {
				return slog.Group("error",
					slog.String("type", errorX.Error()),
					slog.Any("data", errorX.LogData),
				)
			}
		}
	}

	return a
}
