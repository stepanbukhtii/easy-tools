package kafka

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/plugin/kotel"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"

	"github.com/stepanbukhtii/easy-tools/elog"
)

var DefaultRetryInterval = 10 * time.Second

type HandlerFunc func(ctx context.Context, record *kgo.Record) error

type ConsumerGroup struct {
	brokers   []string
	groupName string
	topics    []string
	tracer    *kotel.Tracer
	client    *kgo.Client
	wg        sync.WaitGroup
	mu        sync.Mutex
	isClosed  bool
	handlers  map[string]HandlerFunc
}

func NewConsumerGroup(groupName string, brokers ...string) *ConsumerGroup {
	return &ConsumerGroup{
		brokers:   brokers,
		groupName: groupName,
		tracer:    kotel.NewTracer(),
		handlers:  make(map[string]HandlerFunc),
	}
}

func (c *ConsumerGroup) Shutdown() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isClosed {
		return
	}
	c.isClosed = true

	c.wg.Wait()

	c.client.Close()
}

func (c *ConsumerGroup) Add(topic string, handler HandlerFunc) {
	c.topics = append(c.topics, topic)
	c.handlers[topic] = handler
}

func (c *ConsumerGroup) Consume(ctx context.Context) error {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(c.brokers...),
		kgo.ConsumerGroup(c.groupName),
		kgo.ConsumeTopics(c.topics...),
		kgo.DisableAutoCommit(),
		kgo.AllowAutoTopicCreation(),
		kgo.WithHooks(c.tracer),
	)
	if err != nil {
		return err
	}

	c.client = client

	for {
		fetches := c.client.PollFetches(ctx)
		if err := fetches.Err(); err != nil {
			if errors.Is(err, kgo.ErrClientClosed) {
				slog.InfoContext(ctx, "kafka consumer closed")
				return nil
			}

			slog.With(elog.Err(err)).ErrorContext(ctx, "poll fetches failed, retrying...")
			time.Sleep(DefaultRetryInterval)
			continue
		}

		fetches.EachPartition(func(partition kgo.FetchTopicPartition) {
			go func(partition kgo.FetchTopicPartition) {
				c.wg.Add(1)
				partition.EachRecord(func(record *kgo.Record) {
					c.handle(record)

					if err := c.client.CommitRecords(record.Context, record); err != nil {
						slog.With(elog.Err(err)).ErrorContext(ctx, "commit record failed")
					}
				})
				c.wg.Done()
			}(partition)
		})
	}
}

func (c *ConsumerGroup) handle(record *kgo.Record) {
	handler, ok := c.handlers[record.Topic]
	if !ok {
		return
	}

	ctx, span := c.tracer.WithProcessSpan(record)
	defer span.End()

	logger := slog.With(
		slog.String(string(semconv.MessagingSystemKey), "kafka"),
		slog.String(string(semconv.MessagingDestinationNameKey), record.Topic),
		slog.String(string(semconv.MessagingConsumerGroupNameKey), c.groupName),
		slog.Int64(string(semconv.MessagingKafkaOffsetKey), record.Offset),
		slog.Int64(string(semconv.MessagingDestinationPartitionIDKey), int64(record.Partition)),
	)

	if len(record.Key) > 0 {
		logger = logger.With(slog.String(string(semconv.MessagingKafkaMessageKeyKey), string(record.Key)))
	}

	if len(record.Value) > 0 {
		logger = logger.With(slog.String(elog.MessagingMessageBodyContent, string(record.Value)))
	}

	start := time.Now()
	err := handler(ctx, record)
	duration := time.Now().Sub(start)

	logger = logger.With(slog.Duration(elog.MessagingOperationDuration, duration))

	if err != nil {
		logger.With(elog.Err(err)).ErrorContext(ctx, "kafka message processing failed")
		return
	}

	logger.InfoContext(ctx, "kafka message processed")
}
