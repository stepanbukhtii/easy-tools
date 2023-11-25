package rabbitmq

import (
	"context"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	deliveryMode uint8
	connection   *Connection
	channel      *amqp.Channel
}

func NewPublisher(connection *Connection) (*Publisher, error) {
	channel, err := connection.Channel(context.Background())
	if err != nil {
		return nil, err
	}

	return &Publisher{
		deliveryMode: amqp.Persistent,
		connection:   connection,
		channel:      channel,
	}, nil
}

func (p *Publisher) Publish(exchange, routingKey string, body []byte) error {
	if p.channel.IsClosed() {
		if err := p.createChannel(); err != nil {
			return err
		}
	}

	msg := amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: p.deliveryMode,
		MessageId:    uuid.NewString(),
		Timestamp:    time.Now(),
		Body:         body,
	}
	return p.channel.Publish(exchange, routingKey, false, false, msg)
}

func (p *Publisher) PublishQueue(queue string, body []byte) error {
	if p.channel.IsClosed() {
		if err := p.createChannel(); err != nil {
			return err
		}
	}

	msg := amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: p.deliveryMode,
		MessageId:    uuid.NewString(),
		Timestamp:    time.Now(),
		Body:         body,
	}
	return p.channel.Publish("", queue, false, false, msg)
}

func (p *Publisher) createChannel() error {
	channel, err := p.connection.Channel(context.Background())
	if err != nil {
		return err
	}

	p.channel = channel

	return nil
}
