package kafka

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/stepanbukhtii/easy-tools/elog"
	"github.com/twmb/franz-go/pkg/kgo"
)

const ()

type HandlerFunc func(record *kgo.Record) error

type ConsumerGroup struct {
	client    *kgo.Client
	brokers   []string
	groupName string
	topics    []string
	handlers  map[string]HandlerFunc
}

func NewConsumerGroup(groupName string, brokers ...string) *ConsumerGroup {
	return &ConsumerGroup{
		brokers:   brokers,
		groupName: groupName,
		handlers:  make(map[string]HandlerFunc),
	}
}

func (c *ConsumerGroup) Shutdown() {
	if c.client != nil {
		c.client.Close()
	}
}

func (c *ConsumerGroup) Add(topic string, handler HandlerFunc) {
	c.topics = append(c.topics, topic)
	c.handlers[topic] = handler
}

func (c *ConsumerGroup) Consume(ctx context.Context) error {
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(c.brokers...),
		kgo.ConsumerGroup(c.groupName),
		kgo.ConsumeTopics(c.topics...),
		kgo.AllowAutoTopicCreation(),
	)
	if err != nil {
		return err
	}

	c.client = cl

	for {
		fetches := c.client.PollFetches(ctx)
		if err := fetches.Err(); err != nil {
			if errors.Is(err, kgo.ErrClientClosed) {
				return nil
			}

			slog.With(elog.Err(err)).Error("Poll fetches error")
			time.Sleep(time.Second)
			continue
		}

		iter := fetches.RecordIter()
		for !iter.Done() {
			c.handle(iter.Next())
		}
	}
}

func (c *ConsumerGroup) handle(record *kgo.Record) {
	handler, ok := c.handlers[record.Topic]
	if !ok {
		return
	}

	logger := slog.With(
		slog.String(elog.KafkaTopicName, record.Topic),
		slog.Int64(elog.KafkaTopicOffset, record.Offset),
		slog.Int64(elog.KafkaTopicPartition, int64(record.Partition)),
	)

	if len(record.Key) > 0 {
		logger = logger.With(slog.String(elog.KafkaMessageKey, string(record.Key)))
	}

	start := time.Now()
	err := handler(record)
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
