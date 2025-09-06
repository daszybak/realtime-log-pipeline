package binance_rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQClient interface {
	CreateChannel() (*amqp091.Channel, error)
}

type Client[T any] struct {
	ch *amqp091.Channel
}

// TODO Accept as parameters when creating the `client`.
const (
	exchangeName = "binance"
	queueName    = "book_ticker"
	routingKey   = "book_ticker"
)

func New[T any](client RabbitMQClient) (*Client[T], error) {
	ch, err := client.CreateChannel()
	if err != nil {
		return nil, fmt.Errorf("couldn't create RabbitMQ channel: %w", err)
	}

	// TODO Switch to a "topic"̨ exchange.
	err = ch.ExchangeDeclare(
		exchangeName,
		"direct",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		return nil, fmt.Errorf(
			"couldn't declare RabbitMQ %s exchange: %w",
			exchangeName,
			err,
		)
	}

	queue, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		return nil, fmt.Errorf(
			"couldn't declare RabbitMQ %s queue: %w",
			queueName,
			err,
		)
	}

	err = ch.QueueBind(queue.Name, routingKey, exchangeName, false, nil)
	if err != nil {
		ch.Close()
		return nil, fmt.Errorf(
			"couldn't bind RabbitMQ %s queue to %s exchange with route key %s: %w",
			queueName,
			exchangeName,
			routingKey,
			err,
		)
	}

	return &Client[T]{
		ch: ch,
	}, nil
}

// TODO Refactor.
type Message[T any] struct {
	Data    T              `json:"data"`
	TraceID string         `json:"trace_id"`
	Headers map[string]any `json:"headers"`
}

// TODO Create `pkg/rabbitmq/publisher.go` to create
// reusable publisher.
func (client *Client[T]) Publish(msg *Message[T]) error {
	body, err := json.Marshal(msg.Data)
	if err != nil {
		return fmt.Errorf("couldn't marshal data: %w", err)
	}
	err = client.ch.Publish(
		exchangeName,
		routingKey,
		true,
		false,
		amqp091.Publishing{
			Headers:       msg.Headers,
			CorrelationId: msg.TraceID,
			Body:          body,
		},
	)
	if err != nil {
		return fmt.Errorf("couldn't publish message to RabbitMQ: %w", err)
	}

	return nil
}

// TODO Create `pkg/rabbitmq/consumer.go` to create
// reusable consumer.
func (client *Client[T]) Consume(ctx context.Context) (<-chan *Message[T], error) {
	deliveries, err := client.ch.Consume(queueName, "", false, false, false, true, nil)
	if err != nil {
		return nil, fmt.Errorf("couldn't create consumer %s: %w", queueName, err)
	}
	out := make(chan *Message[T], 1000)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case d, ok := <-deliveries:
				if !ok {
					return
				}
				data := new(T)
				err := json.Unmarshal(d.Body, data)
				if err != nil {
					// TODO Log.
					d.Nack(false, false)
				} else {
					msg := &Message[T]{
						Data:    *data,
						TraceID: d.CorrelationId,
						Headers: d.Headers,
					}
					out <- msg
					d.Ack(false)
				}
			}
		}
	}()
	return out, nil
}
