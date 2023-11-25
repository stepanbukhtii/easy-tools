package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const tracerPublisherName = "nats-publisher"

type Publisher struct {
	connection *nats.Conn
}

func NewPublisher(url string) (*Publisher, error) {
	connection, err := NewConnection(url)
	if err != nil {
		return nil, err
	}
	return &Publisher{connection: connection}, nil
}

func NewPublisherConnection(connection *nats.Conn) *Publisher {
	return &Publisher{connection: connection}
}

func (p *Publisher) Publish(ctx context.Context, subject string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return p.PublishRaw(ctx, subject, data)
}

func (p *Publisher) PublishRaw(ctx context.Context, subject string, data []byte) error {
	spanName := fmt.Sprintf("%s publish", subject)
	ctx, span := otel.Tracer(tracerPublisherName).Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	msg := &nats.Msg{
		Subject: subject,
		Header:  make(nats.Header),
		Data:    data,
	}

	otel.GetTextMapPropagator().Inject(ctx, HeadersCarrier(msg.Header))

	return p.connection.PublishMsg(msg)
}

func (p *Publisher) Close() {
	if p.connection != nil {
		p.connection.Close()
	}
}
