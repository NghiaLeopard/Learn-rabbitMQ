package main

import (
	"Learn-rabbitMQ/internal"
	"log"
	"time"
)

func main() {
	conn, err := internal.ConntectRabbitMQ("guest", "guest", "localhost:5672", "customers")

	if err != nil {
		panic(err)
	}

	defer conn.Close()

	client, err := internal.NewRabbitClient(conn)

	if err != nil {
		panic(err)
	}


	defer client.Close()

	time.Sleep(10 * time.Second)

	log.Print(client)
}
