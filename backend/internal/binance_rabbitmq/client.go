package binance_rabbitmq

import (
	"encoding/json"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQClient interface {
	CreateChannel() (*amqp091.Channel, error)
}

type Client struct {
	ch *amqp091.Channel
}

// TODO Accept as parameters when creating the `client`.
const (
	exchangeName = "binance"
	queueName    = "book_ticker"
	routingKey   = "book_ticker"
)

func New(client RabbitMQClient) (*Client, error) {
	ch, err := client.CreateChannel()
	if err != nil {
		return nil, fmt.Errorf("couldn't create RabbitMQ channel: %w", err)
	}

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

	return &Client{
		ch: ch,
	}, nil
}

// TODO Refactor.
type Message struct {
	Data    any
	TraceID string
	Headers map[string]interface{}
}

// TODO Create `pkg/rabbitmq/publisher.go` to create
// reusable publisher.
func (client *Client) Publish(msg *Message) error {
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

// TODO Add consumer.
