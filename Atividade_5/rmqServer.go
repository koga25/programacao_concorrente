package main

import (
	"context"
	"fmt"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var concurrentClientsCounter = 0
var clientIdCounter = 0

var eventTime = time.Date(2025, 07, 29, 14, 30, 45, 100, time.Local)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide port number")
		return
	}

	URL := "amqp://guest:guest@localhost:" + arguments[1] + "/"

	conn, err := amqp.Dial(URL)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ch.Close()

	// Declarando uma fila
	queue, err := ch.QueueDeclare(
		"timeBetween", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create some logic here to only post content when someone asks
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Hello World!"
	err = ch.PublishWithContext(
		ctx,        // context
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		fmt.Println(err)
		return
	}
}
