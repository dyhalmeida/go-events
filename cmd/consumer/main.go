package main

import (
	"fmt"

	"github.com/dyhalmeida/go-events/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	messages := make(chan amqp.Delivery)

	go rabbitmq.Consumer(ch, messages)

	for message := range messages {
		fmt.Println(string(message.Body))
		message.Ack(false)
	}
}
