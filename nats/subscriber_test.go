package nats

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
)

func TestSubscriber(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	subscriber, err := NewSubscriber("nats://localhost:4222", "queue_name")
	require.NoError(t, err)

	subject := "test-subject"

	err = subscriber.Subscribe(subject, func(ctx context.Context, msg *nats.Msg) error {
		time.Sleep(2 * time.Second)
		fmt.Println("message data", string(msg.Data))
		return nil
	})
	require.NoError(t, err)

	producer, err := NewPublisher("nats://localhost:4222")
	require.NoError(t, err)

	for i := 0; i < 3; i++ {
		require.NoError(t, producer.PublishRaw(context.Background(), subject, []byte("message data")))
	}

	time.Sleep(time.Second)

	time.Sleep(5 * time.Second)
	fmt.Println("shutdown start")
	require.NoError(t, subscriber.Shutdown())
	fmt.Println("shutdown finished")
}
