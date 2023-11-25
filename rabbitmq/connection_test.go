package rabbitmq

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestConnection(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	conn, err := NewConnection("amqp://guest:guest@localhost:5672/")
	require.NoError(t, err)

	go func() {
		time.Sleep(30 * time.Second)
		require.NoError(t, conn.Close())
	}()

	for {
		channel, err := conn.Channel(context.Background())
		if err != nil {
			slog.With(slog.Any("error", err)).Error("could not open RabbitMQ channel")
			break
		}

		slog.With(slog.Bool("is_closed", channel.IsClosed())).Info("channel has been opened")
		time.Sleep(1700 * time.Millisecond)
		require.NoError(t, channel.Close())
	}

	time.Sleep(2 * time.Second)
}
