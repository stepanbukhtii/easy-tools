package kafka

import (
	"context"

	"github.com/google/uuid"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Producer struct {
	client *kgo.Client
}

func NewProducer(brokers ...string) (*Producer, error) {
	client, err := kgo.NewClient(kgo.SeedBrokers(brokers...))
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

func (p *Producer) Produce(ctx context.Context, topic string, body []byte) error {
	record := &kgo.Record{
		Key:   []byte(uuid.NewString()),
		Topic: topic,
		Value: body,
	}

	if err := p.client.ProduceSync(ctx, record).FirstErr(); err != nil {
		return err
	}

	return nil
}
