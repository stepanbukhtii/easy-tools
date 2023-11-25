package rabbitmq

import (
	"context"
	"fmt"
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

	connection, err := NewConnection("amqp://guest:guest@localhost:5672/")
	assert.NoError(t, err)

	queueName := "test"

	require.NoError(t, connection.Declare().Queue(queueName))

	consumerGroup := NewConsumer(connection, "")

	time.Sleep(time.Second)

	consumerGroup.Add(queueName, func(ctx context.Context, msg amqp.Delivery) error {
		fmt.Println("text text", string(msg.Body))
		msg.Ack(false)
		return nil
	})

	require.NoError(t, consumerGroup.Consume(context.Background()))

	time.Sleep(time.Second)
}
