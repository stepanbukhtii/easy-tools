package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stepanbukhtii/easy-tools/econtext"
	"github.com/stepanbukhtii/easy-tools/elog"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
	"go.opentelemetry.io/otel/trace"
)

const tracerPublisherName = "rabbitmq-publisher"

type Publisher struct {
	connection    *Connection
	channel       *amqp.Channel
	deliveryMode  uint8
	wg            sync.WaitGroup
	mu            sync.RWMutex
	reconnectChan chan struct{}
	stopChan      chan struct{}
	connected     bool
}

func NewPublisher(url string) (*Publisher, error) {
	connection, err := NewConnection(url)
	if err != nil {
		return nil, err
	}

	return NewPublisherWithConnection(connection), nil
}

func NewPublisherWithConnection(connection *Connection) *Publisher {
	p := Publisher{
		connection:    connection,
		deliveryMode:  amqp.Persistent,
		reconnectChan: make(chan struct{}),
		stopChan:      make(chan struct{}),
	}

	go p.handleReconnect()

	return &p
}

func (p *Publisher) Close() error {
	close(p.stopChan)
	p.wg.Wait()
	return p.connection.Close()
}

func (p *Publisher) Publish(ctx context.Context, exchange, routingKey string, value any) error {
	body, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return p.PublishRaw(ctx, exchange, routingKey, body)
}

func (p *Publisher) PublishRaw(ctx context.Context, exchange, routingKey string, body []byte) error {
	return p.publish(ctx, exchange, routingKey, body)
}

func (p *Publisher) PublishQueue(ctx context.Context, queue string, value any) error {
	body, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return p.PublishQueueRaw(ctx, queue, body)
}

func (p *Publisher) PublishQueueRaw(ctx context.Context, queue string, body []byte) error {
	return p.publish(ctx, "", queue, body)
}

func (p *Publisher) PublishEvent(ctx context.Context, event Event) error {
	return p.Publish(ctx, event.Exchange(), event.RoutingKey(), event)
}

func (p *Publisher) publish(ctx context.Context, exchange, routingKey string, body []byte) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	spanName := fmt.Sprintf("%s publish", routingKey)
	if exchange != "" {
		spanName = fmt.Sprintf("%s.%s publish", exchange, routingKey)
	}

	spanCtx, span := otel.Tracer(tracerPublisherName).Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	msg := amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: p.deliveryMode,
		MessageId:    uuid.NewString(),
		Timestamp:    time.Now(),
		Body:         body,
		Headers:      make(amqp.Table),
	}

	otel.GetTextMapPropagator().Inject(spanCtx, HeadersCarrier(msg.Headers))

	destination := routingKey
	if exchange != "" {
		destination = fmt.Sprintf("%s:%s", exchange, destination)
	}

	logger := econtext.Logger(spanCtx).With(
		slog.String(string(semconv.MessagingDestinationNameKey), destination),
		slog.String(elog.MessagingMessageBodyContent, string(body)),
	)

	channel, err := p.getChannel(ctx)
	if err != nil {
		return err
	}

	if err := channel.Publish(exchange, routingKey, false, false, msg); err != nil {
		logger.With(elog.Err(err)).ErrorContext(spanCtx, "rabbitmq message publish failed")
		return err
	}

	logger.InfoContext(spanCtx, "rabbitmq message published")

	return nil
}

func (p *Publisher) getChannel(ctx context.Context) (*amqp.Channel, error) {
	p.mu.RLock()
	connected := p.connected
	p.mu.RUnlock()

	if connected {
		return p.channel, nil
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-p.reconnectChan:
		return p.channel, nil
	case <-p.stopChan:
		return nil, ErrConnectionClosed
	}
}

func (p *Publisher) handleReconnect() {
	ctx := context.Background()

	for {
		p.mu.Lock()
		p.connected = false
		p.mu.Unlock()

		channel, err := p.connection.Channel(ctx)
		if err != nil {
			slog.With(elog.Err(err)).ErrorContext(ctx, "rabbitmq publisher closed")
			return
		}

		p.mu.Lock()
		p.channel = channel
		p.connected = true
		p.mu.Unlock()

		close(p.reconnectChan)
		p.reconnectChan = make(chan struct{})

		select {
		case <-p.stopChan:
			slog.InfoContext(ctx, "rabbitmq publisher closed")
			return
		case err := <-channel.NotifyClose(make(chan *amqp.Error, 1)):
			if err != nil {
				slog.With(elog.Err(err)).ErrorContext(ctx, "rabbitmq publish channel closed, recovering...")
				continue
			}

			slog.InfoContext(ctx, "rabbitmq publisher closed")
			return
		}
	}
}
