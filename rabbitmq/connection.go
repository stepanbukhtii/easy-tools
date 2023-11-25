package rabbitmq

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stepanbukhtii/easy-tools/elog"
)

var DefaultRetryInterval = 10 * time.Second

var ErrConnectionClosed = errors.New("connection is closed")

type Connection struct {
	url           string
	closed        bool
	connected     bool
	mu            sync.RWMutex
	reconnectChan chan struct{}
	stopChan      chan struct{}
	conn          *amqp.Connection
}

func NewConnection(url string) (*Connection, error) {
	c := Connection{
		url:           url,
		reconnectChan: make(chan struct{}),
		stopChan:      make(chan struct{}),
	}

	if err := c.connect(); err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *Connection) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}
	c.closed = true
	c.connected = false

	close(c.stopChan)

	if c.conn != nil && !c.conn.IsClosed() {
		return c.conn.Close()
	}

	return nil
}

func (c *Connection) Declare() Declare {
	return Declare{connection: c}
}

func (c *Connection) Channel(ctx context.Context) (*amqp.Channel, error) {
	c.mu.RLock()
	closed := c.closed
	connected := c.connected
	c.mu.RUnlock()

	if closed {
		return nil, ErrConnectionClosed
	}

	if connected {
		return c.conn.Channel()
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-c.reconnectChan:
		return c.conn.Channel()
	case <-c.stopChan:
		return nil, ErrConnectionClosed
	}
}

func (c *Connection) connect() error {
	conn, err := amqp.Dial(c.url)
	if err != nil {
		return err
	}
	c.conn = conn

	c.mu.Lock()
	c.connected = true
	c.mu.Unlock()

	go func() {
		err := <-conn.NotifyClose(make(chan *amqp.Error, 1))
		if err == nil {
			slog.Info("rabbitmq connection closed")
			return
		}

		c.mu.Lock()
		c.connected = false
		c.mu.Unlock()

		slog.With(elog.Err(err)).Error("rabbitmq connection closed, reconnecting...")

		for {
			if c.closed {
				return
			}

			if err := c.connect(); err != nil {
				slog.With(elog.Err(err)).Error("rabbitmq reconnection failed, retrying connection...")

				select {
				case <-time.After(DefaultRetryInterval):
					continue
				case <-c.stopChan:
					return
				}
			}

			close(c.reconnectChan)
			c.reconnectChan = make(chan struct{})

			slog.Info("rabbitmq reconnected")

			return
		}
	}()

	return nil
}
