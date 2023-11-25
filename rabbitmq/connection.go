package rabbitmq

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stepanbukhtii/easy-tools/elog"
)

var DefaultRetryInterval = 10 * time.Second

var ErrConnectionClosed = errors.New("connection is closed")

type Connection struct {
	url    string
	closed bool
	conn   *amqp.Connection
}

func NewConnection(url string) (*Connection, error) {
	c := Connection{url: url}

	if err := c.connect(); err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *Connection) Close() error {
	if c.closed {
		return nil
	}
	c.closed = true

	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}

func (c *Connection) Declare() Declare {
	return Declare{connection: c}
}

func (c *Connection) Channel(ctx context.Context) (*amqp.Channel, error) {
	for {
		if c.closed {
			return nil, ErrConnectionClosed
		}

		ch, err := c.conn.Channel()
		if err != nil {
			slog.With(elog.Err(err)).Error("Could not open RabbitMQ channel")

			select {
			case <-ctx.Done():
				return nil, err
			case <-time.After(time.Second):
				continue
			}
		}

		return ch, nil
	}
}

func (c *Connection) connect() error {
	conn, err := amqp.Dial(c.url)
	if err != nil {
		return err
	}
	c.conn = conn

	go func() {
		err := <-conn.NotifyClose(make(chan *amqp.Error, 1))
		if err != nil {
			slog.With(elog.Err(err)).Error("RabbitMQ connection closed. Reconnecting...")

			for {
				if c.closed {
					return
				}

				if err := c.connect(); err != nil {
					slog.With(elog.Err(err)).Error("RabbitMQ reconnection failed. Retrying connection...")
					time.Sleep(DefaultRetryInterval)
					continue
				}

				log.Println("RabbitMQ reconnected successfully")
				return
			}
		}

		slog.Info("RabbitMQ connection closed. Graceful shutdown.")
	}()

	return nil
}
