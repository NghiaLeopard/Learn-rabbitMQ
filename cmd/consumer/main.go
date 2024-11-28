package main

import (
	"Learn-rabbitMQ/internal"
	"context"
	"log"
	"time"

	"golang.org/x/sync/errgroup"
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

	messageBus, err := client.Consume("customers_created", "email-service", false)

	if err != nil {
		panic(err)
	}

	var blocking chan struct{}

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)

	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	g.SetLimit(10)

	// go func() {
	// 	for message := range messageBus {
	// 		log.Println("new message: %v", message)

	// 		// lần đầu chưa được gửi nên message.Redelivered = false mà không thay đổi logic hoặc trạng thái gì để phân biệt khi nhận lại
	// 		// lần sau thành true nên được bỏ qua và sẽ xóa message đó
	// 		if !message.Redelivered {
	// 			// Nó sẽ được gửi lại và đưa vào hàng chờ
	// 			message.Nack(false, true)
	// 			continue
	// 		}

	// 		if err := message.Ack(false); err != nil {
	// 			log.Println("message failure")
	// 			continue
	// 		}

	// 		log.Println("Acknowledge message: %s /n", message.MessageId)

	// 	}
	// }()

	go func() {
		for message := range messageBus {
			msg := message
			g.Go(func() error {
				log.Printf("new message: %v", msg)
				time.Sleep(10 * time.Second)
				if err := msg.Ack(false); err != nil {
					log.Println("Ack message failed")
					return err
				}
				return nil
			})
		}
	}()

	log.Println("Consuming, to close the program press CTRL + C")

	<-blocking
}
