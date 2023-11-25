package recovery

import (
	"context"
	"log/slog"
	"runtime/debug"

	"github.com/stepanbukhtii/easy-tools/econtext"
	"github.com/stepanbukhtii/easy-tools/elog"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
)

func Recover() {
	r := recover()
	if r == nil {
		return
	}

	if err, ok := r.(error); ok {
		slog.With(
			elog.Err(err),
			slog.String(string(semconv.ExceptionStacktraceKey), string(debug.Stack())),
		).Error("panic recovered")
		return
	}

	slog.With(slog.String(string(semconv.ExceptionStacktraceKey), string(debug.Stack()))).Error("panic recovered")
}

func RecoverContext(ctx context.Context) {
	r := recover()
	if r == nil {
		return
	}

	logger := econtext.Logger(ctx)
	if !logger.Enabled(ctx, slog.LevelError) {
		logger = slog.Default()
	}
	logger = logger.With(slog.Any(string(semconv.ExceptionStacktraceKey), string(debug.Stack())))

	if err, ok := r.(error); ok {
		logger.With(elog.Err(err)).ErrorContext(ctx, "panic recovered")
		return
	}

	logger.ErrorContext(ctx, "panic recovered")
}

func Go(f func()) {
	go func() {
		defer Recover()
		f()
	}()
}

func GoContext(ctx context.Context, f func(context.Context)) {
	go func(ctx context.Context) {
		defer RecoverContext(ctx)
		f(ctx)
	}(context.WithoutCancel(ctx))
}
