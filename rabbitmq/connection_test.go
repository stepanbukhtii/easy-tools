package rabbitmq

import (
	"context"
	"fmt"
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

	time.Sleep(time.Second)

	go func() {
		time.Sleep(10 * time.Second)
		conn.Close()
	}()

	for {
		channel, err := conn.Channel(context.Background())
		if err != nil {
			fmt.Println("channel err:", err)
			break
		}

		fmt.Println(channel)
		time.Sleep(2 * time.Second)
	}

	time.Sleep(time.Second)
}
