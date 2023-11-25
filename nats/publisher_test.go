package nats

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPublisher(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	producer, err := NewPublisher("nats://localhost:4222")
	require.NoError(t, err)

	for {
		require.NoError(t, producer.PublishRaw(context.Background(), "test-topic", []byte("message data")))
		slog.Info("Message is sent")
		time.Sleep(time.Second)
	}
}
