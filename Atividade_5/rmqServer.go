package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
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

	// Declarando uma fila de envio de respostas
	responseQueue, err := ch.QueueDeclare(
		"responseTimeBetween", // name
		false,                 // durable
		false,                 // delete when unused
		false,                 // exclusive
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Declarando uma fila de requisiçoes de evento
	requestQueue, err := ch.QueueDeclare(
		"requestTimeBetween", // name
		false,                // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	requests, err := ch.Consume(
		requestQueue.Name, // queue
		"",                // consumer
		true,              // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	for r := range requests {
		var request = string(r.Body)
		fmt.Println(request)
		if request == "sendTimeBetween" {
			var today = time.Now()
			var timeBetween = int(eventTime.Sub(today).Seconds())
			err = ch.PublishWithContext(
				ctx,                // context
				"",                 // exchange
				responseQueue.Name, // routing key
				false,              // mandatory
				false,              // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(strconv.Itoa(timeBetween)),
				},
			)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
