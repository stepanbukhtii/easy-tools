package rabbitmq

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConsumerGroup(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	ctx := context.Background()
	queueName := "test"

	connection, err := NewConnection("amqp://guest:guest@localhost:5672/")
	assert.NoError(t, err)

	require.NoError(t, connection.Declare().Queue(queueName))

	publisher := NewPublisherWithConnection(connection)

	consumerGroup := NewConsumerWithConnection(connection, "consumer_tag")

	consumerGroup.Add(queueName, func(ctx context.Context, msg amqp.Delivery) error {
		time.Sleep(1 * time.Second)
		return nil
	})

	for i := 0; i < 10; i++ {
		require.NoError(t, publisher.PublishQueue(ctx, queueName, "test message"))
	}

	go func() {
		consumerGroup.Consume(ctx)
	}()

	time.Sleep(30 * time.Second)

	slog.Info("shutdown start")
	require.NoError(t, consumerGroup.Shutdown())
	slog.Info("shutdown finished")
}
