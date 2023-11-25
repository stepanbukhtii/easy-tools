package rabbitmq

import (
	"context"
	"log/slog"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stepanbukhtii/easy-tools/elog"
)

var DefaultPrefetchCount = 10

type HandlerFunc func(ctx context.Context, msg amqp.Delivery) error

type ConsumerOptions struct {
	Queue         string
	ConsumerTag   string
	PrefetchCount int
	Handler       HandlerFunc
	Args          amqp.Table
}

type ConsumerGroup struct {
	connection        *Connection
	globalConsumerTag string
	consumers         []ConsumerOptions
}

func NewConsumer(connection *Connection, consumerTag string) *ConsumerGroup {
	return &ConsumerGroup{
		connection:        connection,
		globalConsumerTag: consumerTag,
	}
}

func (c *ConsumerGroup) Shutdown() error {
	return c.connection.Close()
}

func (c *ConsumerGroup) Add(queue string, handler HandlerFunc) {
	consumerOptions := ConsumerOptions{
		Queue:         queue,
		PrefetchCount: DefaultPrefetchCount,
		Handler:       handler,
	}
	c.consumers = append(c.consumers, consumerOptions)
}

func (c *ConsumerGroup) AddOption(consumerOptions ConsumerOptions) {
	if consumerOptions.PrefetchCount <= 0 {
		consumerOptions.PrefetchCount = DefaultPrefetchCount
	}
	c.consumers = append(c.consumers, consumerOptions)
}

func (c *ConsumerGroup) Consume(ctx context.Context) error {
	var wg sync.WaitGroup

	for _, consumer := range c.consumers {
		if err := c.consume(ctx, &wg, consumer); err != nil {
			return err
		}
	}

	wg.Wait()

	return nil
}

func (c *ConsumerGroup) consume(ctx context.Context, wg *sync.WaitGroup, consumer ConsumerOptions) error {
	channel, err := c.connection.Channel(ctx)
	if err != nil {
		return err
	}

	if err := channel.Qos(consumer.PrefetchCount, 0, false); err != nil {
		_ = channel.Close()
		return err
	}

	if consumer.ConsumerTag == "" {
		consumer.ConsumerTag = c.globalConsumerTag
	}

	deliveries, err := channel.Consume(consumer.Queue, consumer.ConsumerTag, false, false, false, false, consumer.Args)
	if err != nil {
		_ = channel.Close()
		return err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case err := <-channel.NotifyClose(make(chan *amqp.Error)):
				if err != nil {
					slog.With(elog.Err(err), slog.Any("queue", consumer.Queue)).Error("Channel closed. Recovering consumer")

					for {
						if err := c.consume(ctx, wg, consumer); err != nil {
							slog.With(elog.Err(err)).Error("RabbitMQ recovering consumer failed. Retrying recovering...")
							time.Sleep(DefaultRetryInterval)
							continue
						}

						slog.Info("RabbitMQ consumer recovering successfully")
						return
					}
				}
			case msg := <-deliveries:
				c.handle(consumer, msg)
			}
		}
	}()

	return nil
}

func (c *ConsumerGroup) handle(consumer ConsumerOptions, msg amqp.Delivery) {
	logger := slog.With(
		slog.String(elog.RabbitMQQueueName, consumer.Queue),
		slog.String(elog.RabbitMQConsumerTag, consumer.ConsumerTag),
	)

	start := time.Now()
	err := consumer.Handler(context.Background(), msg)
	end := time.Now()
	latency := end.Sub(start)

	logger = logger.With(
		slog.Time(elog.EventStart, start),
		slog.Time(elog.EventEnd, end),
		slog.Duration(elog.EventDuration, latency),
	)

	if err != nil {
		logger.With(elog.Err(err)).Error("Failed to process message.")
		return
	}

	logger.Info("Successfully processed message.")
}
