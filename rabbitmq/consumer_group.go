package rabbitmq

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stepanbukhtii/easy-tools/elog"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
	"go.opentelemetry.io/otel/trace"
)

const tracerConsumerName = "rabbitmq-consumer"

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
	url               string
	globalConsumerTag string
	connection        *Connection
	wg                sync.WaitGroup
	mu                sync.Mutex
	isClosed          bool
	stopChan          chan struct{}
	consumers         []ConsumerOptions
}

func NewConsumerGroup(url string, consumerTag string) *ConsumerGroup {
	return &ConsumerGroup{
		url:               url,
		globalConsumerTag: consumerTag,
		stopChan:          make(chan struct{}),
	}
}

func NewConsumerWithConnection(connection *Connection, consumerTag string) *ConsumerGroup {
	return &ConsumerGroup{
		connection:        connection,
		globalConsumerTag: consumerTag,
		stopChan:          make(chan struct{}),
	}
}

func (c *ConsumerGroup) Shutdown() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isClosed {
		return nil
	}
	c.isClosed = true

	close(c.stopChan)

	c.wg.Wait()

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

func (c *ConsumerGroup) AddWithOption(consumerOptions ConsumerOptions) {
	if consumerOptions.PrefetchCount <= 0 {
		consumerOptions.PrefetchCount = DefaultPrefetchCount
	}
	c.consumers = append(c.consumers, consumerOptions)
}

func (c *ConsumerGroup) Consume(ctx context.Context) error {
	if c.url != "" || c.connection == nil {
		connection, err := NewConnection(c.url)
		if err != nil {
			return err
		}
		c.connection = connection
	}

	for i := range c.consumers {
		go c.listen(ctx, c.consumers[i])
	}

	<-c.stopChan

	return nil
}

func (c *ConsumerGroup) listen(ctx context.Context, consumer ConsumerOptions) {
	logger := slog.With(slog.String(string(semconv.MessagingDestinationNameKey), consumer.Queue))

	for {
		select {
		case <-c.stopChan:
			logger.InfoContext(ctx, "rabbitmq consumer closed")
			return
		default:
			if err := c.consume(ctx, consumer); err != nil {
				logger.With(elog.Err(err)).ErrorContext(ctx, "rabbitmq consumer failed, recovering...")

				select {
				case <-time.After(DefaultRetryInterval):
					continue
				case <-c.stopChan:
					logger.InfoContext(ctx, "rabbitmq consumer closed")
					return
				}
			}

			logger.InfoContext(ctx, "rabbitmq consumer closed")
			return
		}
	}
}

func (c *ConsumerGroup) consume(ctx context.Context, consumer ConsumerOptions) error {
	channel, err := c.connection.Channel(ctx)
	if err != nil {
		return err
	}

	if err = channel.Qos(consumer.PrefetchCount, 0, false); err != nil {
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

	for {
		select {
		case <-c.stopChan:
			return nil
		case err := <-channel.NotifyClose(make(chan *amqp.Error, 1)):
			return err
		case msg, ok := <-deliveries:
			if ok {
				go c.handle(ctx, consumer, msg)
			}
		}
	}
}

func (c *ConsumerGroup) handle(baseCtx context.Context, consumer ConsumerOptions, msg amqp.Delivery) {
	c.wg.Add(1)
	defer c.wg.Done()

	parentCtx := otel.GetTextMapPropagator().Extract(baseCtx, HeadersCarrier(msg.Headers))

	spanName := fmt.Sprintf("%s process", consumer.Queue)
	ctx, span := otel.Tracer(tracerConsumerName).Start(parentCtx, spanName, trace.WithSpanKind(trace.SpanKindConsumer))
	defer span.End()

	destination := consumer.Queue
	if msg.Exchange != "" {
		destination = fmt.Sprintf("%s:%s", msg.Exchange, destination)
	}

	logger := slog.With(
		slog.String(string(semconv.MessagingSystemKey), "rabbitmq"),
		slog.String(string(semconv.MessagingDestinationNameKey), destination),
		slog.String(string(semconv.MessagingConsumerGroupNameKey), consumer.ConsumerTag),
	)

	if msg.RoutingKey != "" {
		logger = logger.With(slog.String(string(semconv.MessagingRabbitMQDestinationRoutingKeyKey), msg.RoutingKey))
	}

	if len(msg.Body) > 0 {
		logger = logger.With(slog.String(elog.MessagingMessageBodyContent, string(msg.Body)))
	}

	start := time.Now()
	err := consumer.Handler(ctx, msg)
	duration := time.Now().Sub(start)

	logger = logger.With(slog.Duration(elog.MessagingOperationDuration, duration))

	if err != nil {
		if err := msg.Nack(false, true); err != nil {
			logger.With(elog.Err(err)).ErrorContext(ctx, "rabbitmq message nack failed")
			return
		}
		logger.With(elog.Err(err)).ErrorContext(ctx, "rabbitmq message processing failed")
		return
	}

	if err := msg.Ack(false); err != nil {
		logger.With(elog.Err(err)).ErrorContext(ctx, "rabbitmq message ack failed")
		return
	}

	logger.InfoContext(ctx, "rabbitmq message processed")
}
