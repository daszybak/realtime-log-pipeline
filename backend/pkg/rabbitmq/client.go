// Package rabbitmq wraps RabbitMQ into a client.
package rabbitmq

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type Client struct {
	connection     *amqp091.Connection
	connectionName string
}

func New(
	connectionName string,
	rabbitMQURL string,
) (*Client, error) {
	conn, err := amqp091.Dial(rabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("couldn't connect to RabbitMQ: %w", err)
	}

	conn.Properties.SetClientConnectionName(connectionName)
	// TODO Add a listener with `ChannelLogger.NotifyReturn` to
	// to handle any undeliverd messages
	return &Client{
		connection:     conn,
		connectionName: connectionName,
	}, nil
}

func (client *Client) GetConnection() *amqp091.Connection {
	return client.connection
}

func (client *Client) GetConnectionName() string {
	return client.connectionName
}

func (client *Client) IsConnected() bool {
	return client.connection != nil && !client.connection.IsClosed()
}

func (client *Client) Close() error {
	if client.connection != nil && !client.connection.IsClosed() {
		return client.connection.Close()
	}
	return nil
}

func (client *Client) CreateChannel() (*amqp091.Channel, error) {
	if !client.IsConnected() {
		return nil, fmt.Errorf("RabbitMQ connection is not available")
	}

	ch, err := client.connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("couldn't create RabbitMQ channel: %w", err)
	}

	return ch, nil
}
