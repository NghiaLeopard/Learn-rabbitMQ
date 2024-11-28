package internal

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	// Connect rabbitmq
	conn *amqp.Connection

	// Connect channel , used to send message
	ch *amqp.Channel
}

func ConnectRabbitMQ(username, password, host, vhost string) (*amqp.Connection, error) {
	return amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, host, vhost))
}

func NewRabbitClient(conn *amqp.Connection) (RabbitClient, error) {
	ch, err := conn.Channel()

	if err != nil {
		return RabbitClient{}, err
	}

	// để kích hoạt PublishWithDeferredConfirmWithContext
	if err := ch.Confirm(false); err != nil {
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

func (rc *RabbitClient) CreateQueue(userName string, durable, autoDelete bool) error {

	_, err := rc.ch.QueueDeclare(userName, durable, autoDelete, false, false, nil)

	if err != nil {
		return err
	}

	return nil
}

func (rc *RabbitClient) CreateBinding(name, binding, exchange string) error {
	err := rc.ch.QueueBind(name, binding, exchange, false, nil)

	return err
}

func (rc *RabbitClient) Send(ctx context.Context, exchange, routingKey string, options amqp.Publishing) error {

	// 	Nếu mandatory = true:
	// RabbitMQ kiểm tra xem có queue nào được ràng buộc (binding) với exchange và khớp với routing key của message.
	// Nếu không có queue nào khớp:
	// Message được trả về cho producer thông qua cơ chế Basic.Return.
	// Nếu có ít nhất một queue phù hợp, message được gửi đến queue đó.
	// Nếu mandatory = false:
	// Message bị hủy (dropped) nếu không có queue nào phù hợp, và producer không nhận được thông báo.

	// 	Nếu immediate = true:
	// RabbitMQ chỉ gửi message đến queue nếu queue đó có ít nhất một consumer đang active.
	// Nếu không có consumer nào, message được trả về cho producer.
	// Nếu immediate = false:
	// Message được đẩy vào queue bất kể trạng thái của consumer.
	// Consumer sẽ nhận được message sau khi nó trở thành active.
	confirmation, err := rc.ch.PublishWithDeferredConfirmWithContext(ctx, exchange, routingKey, true, false, options)

	if err != nil {
		return err
	}

	// Đợi đến khi nào server thông báo là nhận được tin nhắn
	log.Println(confirmation.Wait())
	return nil
}

func (rc *RabbitClient) Consume(queue, consumer string, autoAck bool) (<-chan amqp.Delivery, error) {
	return rc.ch.Consume(queue, consumer, autoAck, false, false, false, nil)
}
