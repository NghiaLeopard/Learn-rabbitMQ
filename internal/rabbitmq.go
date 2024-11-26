package internal

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	// Connect rabbitmq
	conn *amqp.Connection

	// COnnect channel , used to send message
	ch *amqp.Channel
}

func ConntectRabbitMQ(username, password, host, vhost string) (*amqp.Connection, error) {
	return amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, host, vhost))
}

func NewRabbitClient(conn *amqp.Connection) (RabbitClient, error) {
	ch, err := conn.Channel()

	if err != nil {
		return RabbitClient{}, err
	}

	return RabbitClient{
		conn: conn,
		ch:   ch,
	}, nil
}

func (rc *RabbitClient) Close() error {
	err := rc.ch.Close()

	if err != nil {
		return err
	}

	return nil
}
