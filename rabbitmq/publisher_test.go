package rabbitmq

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPublisher(t *testing.T) {
	publisher, err := NewPublisher("amqp://guest:guest@localhost:5672/")
	require.NoError(t, err)

	go func() {
		time.Sleep(20 * time.Second)
		publisher.Close()
	}()

	data := []byte("message data")
	for {
		require.NoError(t, publisher.PublishQueue(context.Background(), "test", data))
		fmt.Println("Message is sent")
		time.Sleep(2 * time.Second)
	}
}
