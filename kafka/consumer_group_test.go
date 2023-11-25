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

	consumerGroup.Add("test-topic", func(record *kgo.Record) error {
		time.Sleep(time.Second)
		fmt.Println("text text", string(record.Value))
		return nil
	})

	go func() {
		time.Sleep(30 * time.Second)
		fmt.Println("Shutdown")
		consumerGroup.Shutdown()
	}()

	require.NoError(t, consumerGroup.Consume(context.Background()))

	time.Sleep(time.Second)
}
