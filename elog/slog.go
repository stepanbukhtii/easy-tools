package elog

import (
	"errors"
	"log/slog"
	"os"

	"github.com/stepanbukhtii/easy-tools/config"
	"github.com/stepanbukhtii/easy-tools/errx"
)

func NewSlog(logConfig config.Log, serviceConfig config.Service) *slog.Logger {
	return slog.New(NewCustomHandler(baseHandler(logConfig, serviceConfig), nil))
}

func NewSlogCustomHandle(logConfig config.Log, serviceConfig config.Service, customHandle CustomHandle) *slog.Logger {
	return slog.New(NewCustomHandler(baseHandler(logConfig, serviceConfig), customHandle))
}

func Err(err error) slog.Attr {
	if err == nil {
		return slog.Attr{}
	}
	return slog.Any(ErrorMessage, err)
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
		slog.String(ServiceName, serviceConfig.Name),
		slog.String(ServiceEnvironment, serviceConfig.Environment),
	}

	if serviceConfig.Version != "" {
		defaultAttrs = append(defaultAttrs, slog.String(ServiceVersion, serviceConfig.Version))
	}

	return slog.NewJSONHandler(os.Stdout, opts).WithAttrs(defaultAttrs)
}

func slogReplaceAttr(_ []string, a slog.Attr) slog.Attr {
	if a.Key == slog.MessageKey {
		a.Key = "message"
	}

	if a.Key == ErrorMessage {
		if err, ok := a.Value.Any().(error); ok {
			var errorX errx.Error
			if errors.As(err, &errorX) && errorX.LogData != nil {
				return slog.Group("error",
					slog.String("message", errorX.Error()),
					slog.Any("data", errorX.LogData),
				)
			}
		}
	}

	return a
}
