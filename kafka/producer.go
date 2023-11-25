package kafka

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"
	"github.com/stepanbukhtii/easy-tools/econtext"
	"github.com/stepanbukhtii/easy-tools/elog"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/plugin/kotel"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
)

type Producer struct {
	client *kgo.Client
}

func NewProducer(brokers ...string) (*Producer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.AllowAutoTopicCreation(),
		kgo.WithHooks(kotel.NewTracer()),
	)
	if err != nil {
		return nil, err
	}

	return &Producer{client: client}, nil
}

func (p *Producer) Close() {
	if p.client != nil {
		p.client.Close()
	}
}

func (p *Producer) Produce(ctx context.Context, topic string, value any) error {
	body, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return p.ProduceRaw(ctx, topic, body)
}

func (p *Producer) ProduceRaw(ctx context.Context, topic string, body []byte) error {
	record := &kgo.Record{
		Key:   []byte(uuid.NewString()),
		Topic: topic,
		Value: body,
	}

	logger := econtext.Logger(ctx).With(
		slog.String(string(semconv.MessagingDestinationNameKey), topic),
		slog.String(elog.MessagingMessageBodyContent, string(body)),
	)

	if err := p.client.ProduceSync(ctx, record).FirstErr(); err != nil {
		logger.With(elog.Err(err)).ErrorContext(ctx, "kafka message producer failed")
		return err
	}

	logger.InfoContext(ctx, "kafka message produced")

	return nil
}

func (p *Producer) ProduceEvent(ctx context.Context, event Event) error {
	return p.Produce(ctx, event.Topic(), event)
}
