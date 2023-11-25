package nats

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stepanbukhtii/easy-tools/elog"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
	"go.opentelemetry.io/otel/trace"
)

const tracerSubscriberName = "nats-subscriber"

type HandlerFunc func(ctx context.Context, msg *nats.Msg) error

type Subscriber struct {
	queue         string
	connection    *nats.Conn
	subscriptions []*nats.Subscription
	wg            sync.WaitGroup
	mu            sync.RWMutex
	isClosed      bool
	stopChan      chan struct{}
}

func NewSubscriber(url, queue string) (*Subscriber, error) {
	connection, err := NewConnection(url)
	if err != nil {
		return nil, err
	}
	return &Subscriber{
		connection: connection,
		queue:      queue,
	}, nil
}

func NewSubscriberConnection(connection *nats.Conn, queue string) *Subscriber {
	return &Subscriber{
		connection: connection,
		queue:      queue,
	}
}

func (p *Subscriber) Subscribe(subject string, handler HandlerFunc) error {
	subscription, err := p.connection.Subscribe(subject, p.handle(handler))
	if err != nil {
		return err
	}

	p.subscriptions = append(p.subscriptions, subscription)

	return nil
}

func (p *Subscriber) QueueSubscribe(subject string, handler HandlerFunc) error {
	subscription, err := p.connection.QueueSubscribe(subject, p.queue, p.handle(handler))
	if err != nil {
		return err
	}

	p.subscriptions = append(p.subscriptions, subscription)

	return nil
}

func (p *Subscriber) Shutdown() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.isClosed {
		return nil
	}
	p.isClosed = true

	for _, subscription := range p.subscriptions {
		if err := subscription.Drain(); err != nil {
			return err
		}
	}

	p.wg.Wait()

	p.connection.Close()

	return nil
}

func (p *Subscriber) handle(handler HandlerFunc) nats.MsgHandler {
	return func(msg *nats.Msg) {
		p.mu.RLock()
		defer p.mu.RUnlock()

		parentCtx := otel.GetTextMapPropagator().Extract(context.Background(), HeadersCarrier(msg.Header))

		spanName := fmt.Sprintf("%s process", msg.Subject)
		ctx, span := otel.Tracer(tracerSubscriberName).Start(parentCtx, spanName, trace.WithSpanKind(trace.SpanKindConsumer))
		defer span.End()

		logger := slog.With(
			slog.String(string(semconv.MessagingSystemKey), "nats"),
			slog.String(string(semconv.MessagingDestinationNameKey), msg.Subject),
		)

		if len(msg.Data) > 0 {
			logger = logger.With(slog.String(elog.MessagingMessageBodyContent, string(msg.Data)))
		}

		start := time.Now()
		err := handler(ctx, msg)
		duration := time.Now().Sub(start)

		logger = logger.With(slog.Duration(elog.MessagingOperationDuration, duration))
		if err != nil {
			logger.With(elog.Err(err)).ErrorContext(ctx, "nats message processing failed")
			return
		}

		logger.InfoContext(ctx, "nats message processed")
	}
}
