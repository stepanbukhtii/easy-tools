package recovery

import (
	"context"
	"log/slog"
	"runtime"

	"github.com/stepanbukhtii/easy-tools/econtext"
	"github.com/stepanbukhtii/easy-tools/elog"
)

func Recover() {
	r := recover()
	if r == nil {
		return
	}

	buf := make([]byte, 1<<14) // 16Kb buff
	stackLen := runtime.Stack(buf, false)

	if err, ok := r.(error); ok {
		slog.With(slog.Any(elog.ErrorStackTrace, buf[:stackLen])).Error("panic recovered", elog.Err(err))
		return
	}

	slog.With(slog.Any(elog.ErrorStackTrace, buf[:stackLen])).Error("panic recovered")
}

func RecoverContext(ctx context.Context) {
	r := recover()
	if r == nil {
		return
	}

	buf := make([]byte, 1<<14) // 16Kb buff
	stackLen := runtime.Stack(buf, false)

	logger := econtext.Logger(ctx)
	if !logger.Enabled(ctx, slog.LevelError) {
		logger = slog.Default()
	}
	logger = logger.With(slog.Any(elog.ErrorStackTrace, buf[:stackLen]))

	if err, ok := r.(error); ok {
		logger.Error("panic recovered", elog.Err(err))
		return
	}

	logger.Error("panic recovered")
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
