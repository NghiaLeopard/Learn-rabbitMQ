package main

import (
	"Learn-rabbitMQ/internal"
	"context"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := internal.ConnectRabbitMQ("guest", "guest", "localhost:5672", "customer")

	if err != nil {
		panic(err)
	}

	defer conn.Close()

	client, err := internal.NewRabbitClient(conn)

	if err != nil {
		panic(err)
	}

	defer client.Close()

	if err := client.CreateQueue("customers_created", true, false); err != nil {
		panic(err)
	}

	if err := client.CreateQueue("customers_test", false, true); err != nil {
		panic(err)
	}

	if err := client.CreateBinding("customers_created", "customers.created.*", "customer_event"); err != nil {
		panic(err)
	}

	if err := client.CreateBinding("customers_test", "customers.*", "customer_event"); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	for i := 0; i < 10; i++ {
		if err := client.Send(ctx, "customer_event", "customers.created.us", amqp091.Publishing{
			ContentType:  "text/plain",
			DeliveryMode: amqp091.Persistent,
			Body:         []byte(`An cool message between service`),
		}); err != nil {
			panic(err)
		}
	}

	time.Sleep(10 * time.Second)

	log.Print(client)
}
