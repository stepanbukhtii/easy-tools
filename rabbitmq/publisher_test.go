package rabbitmq

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPublisher(t *testing.T) {
	connection, err := NewConnection("amqp://guest:guest@localhost:5672/")
	publisher, err := NewPublisher(connection)
	require.NoError(t, err)

	data := []byte("message data")
	for {
		require.NoError(t, publisher.PublishQueue("test", data))
		fmt.Println("Message is sent")
		time.Sleep(5 * time.Second)
	}
}
