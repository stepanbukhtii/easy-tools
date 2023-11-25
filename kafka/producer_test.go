package kafka

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPublisher(t *testing.T) {
	producer, err := NewProducer("localhost:9092")
	require.NoError(t, err)

	data := []byte("message data")
	for {
		require.NoError(t, producer.Produce(context.Background(), "test-topic", data))
		fmt.Println("Message is sent")
		time.Sleep(500 * time.Millisecond)
	}
}
