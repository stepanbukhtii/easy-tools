package recovery

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/stepanbukhtii/easy-tools/econtext"
	"github.com/stretchr/testify/require"
)

func TestRecover(t *testing.T) {
	var logWriter bytes.Buffer
	slog.SetDefault(slog.New(slog.NewJSONHandler(&logWriter, nil)))
	panicError := errors.New("panic recovered")

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer Recover()

		panic(panicError)
	}()

	wg.Wait()
	require.GreaterOrEqual(t, logWriter.Len(), 842)

	Go(func() { panic(panicError) })

	time.Sleep(time.Second)

	require.GreaterOrEqual(t, logWriter.Len(), 1826)
}

func TestRecoverContext(t *testing.T) {
	var logWriter bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logWriter, nil))
	ctx, cancel := context.WithCancel(econtext.WithLogger(context.Background(), logger))
	cancel()
	panicError := errors.New("panic recovered")

	var wg sync.WaitGroup
	wg.Add(1)

	go func(ctx context.Context) {
		defer wg.Done()
		defer RecoverContext(ctx)

		time.Sleep(500 * time.Millisecond)

		require.NoError(t, ctx.Err())

		panic(panicError)
	}(context.WithoutCancel(ctx))

	wg.Wait()
	require.GreaterOrEqual(t, logWriter.Len(), 954)

	GoContext(ctx, func(ctx context.Context) {
		time.Sleep(500 * time.Millisecond)

		require.NoError(t, ctx.Err())

		panic(panicError)
	})

	require.Equal(t, ctx.Err(), context.Canceled)

	time.Sleep(time.Second)

	require.GreaterOrEqual(t, logWriter.Len(), 2094)
}
