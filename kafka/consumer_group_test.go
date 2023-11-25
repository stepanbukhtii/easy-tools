package kafka

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/twmb/franz-go/pkg/kgo"
)

func TestConsumerGroup(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	consumerGroup := NewConsumerGroup("groupName", "localhost:9092")

	consumerGroup.Add("test-topic", func(_ context.Context, record *kgo.Record) error {
		time.Sleep(2 * time.Second)
		fmt.Println("text text", string(record.Value))
		return nil
	})

	producer, err := NewProducer("localhost:9092")
	require.NoError(t, err)

	for i := 0; i < 3; i++ {
		require.NoError(t, producer.Produce(context.Background(), "test-topic", []byte("message data")))
	}

	time.Sleep(time.Second)

	go func() {
		require.NoError(t, consumerGroup.Consume(context.Background()))
	}()

	time.Sleep(5 * time.Second)
	fmt.Println("Shutdown start")
	consumerGroup.Shutdown()
	fmt.Println("Shutdown finished")
}
