package rabbitmq

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Declare struct {
	connection *Connection
}

func (d Declare) Queue(queue string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	channel, err := d.connection.Channel(ctx)
	if err != nil {
		return err
	}

	_, err = channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	return channel.Close()
}

func (d Declare) QueueExchange(queue, exchange string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	channel, err := d.connection.Channel(ctx)
	if err != nil {
		return err
	}

	q, err := channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	if exchange != "" {
		err = channel.ExchangeDeclare(exchange, amqp.ExchangeDirect, true, false, false, false, nil)
		if err != nil {
			return err
		}

		if err := channel.QueueBind(q.Name, queue, exchange, false, nil); err != nil {
			return err
		}
	}

	return channel.Close()
}

func (d Declare) QueueExchangeRoutingKey(queue, exchange, routingKey string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	channel, err := d.connection.Channel(ctx)
	if err != nil {
		return err
	}

	q, err := channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	if exchange != "" {
		err = channel.ExchangeDeclare(exchange, amqp.ExchangeDirect, true, false, false, false, nil)
		if err != nil {
			return err
		}

		if err := channel.QueueBind(q.Name, routingKey, exchange, false, nil); err != nil {
			return err
		}
	}

	return channel.Close()
}

func (d Declare) QueueRetryDeclare(queue, exchange string, retryPeriod time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	channel, err := d.connection.Channel(ctx)
	if err != nil {
		return err
	}

	retryQueue := fmt.Sprintf("%s.retry", queue)

	queueArgs := amqp.Table{
		"x-dead-letter-exchange":    exchange,
		"x-dead-letter-routing-key": retryQueue,
	}

	q, err := channel.QueueDeclare(queue, true, false, false, false, queueArgs)
	if err != nil {
		return err
	}

	if exchange != "" {
		err = channel.ExchangeDeclare(exchange, amqp.ExchangeDirect, true, false, false, false, nil)
		if err != nil {
			return err
		}

		if err := channel.QueueBind(q.Name, queue, exchange, false, nil); err != nil {
			return err
		}
	}

	retryArgs := amqp.Table{
		"x-dead-letter-exchange":    exchange,
		"x-dead-letter-routing-key": q.Name,
		"x-message-ttl":             retryPeriod.Milliseconds(),
	}
	retryQ, err := channel.QueueDeclare(retryQueue, true, false, false, false, retryArgs)
	if err != nil {
		return err
	}

	if err := channel.QueueBind(retryQ.Name, retryQueue, exchange, false, nil); err != nil {
		return err
	}

	return channel.Close()
}
