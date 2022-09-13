package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide host:port.")
		os.Exit(0)
	}

	fmt.Println("Starting connection at RabbitMQ in address " + arguments[1])

	CONNECT := arguments[1]

	conn, err := amqp.Dial("amqp://guest:guest@localhost:" + CONNECT + "/")
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

	messages, err := ch.Consume(
		responseQueue.Name, // queue
		"",                 // consumer
		true,               // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)
	if err != nil {
		fmt.Println(err)
		return
	}

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	test_flag := arguments[2]

	tries := 10000
	var rtts = make([]int, tries+1)
	var wg sync.WaitGroup
	for i := 0; i < tries; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Printf("%i\n", i)
			var initialTime = time.Now()

			fmt.Println("Sending new request")

			body := "sendTimeBetween"
			var err = ch.PublishWithContext(
				ctx,               // context
				"",                // exchange
				requestQueue.Name, // routing key
				false,             // mandatory
				false,             // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(body),
				},
			)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println("Entrando no loop de response")

			mes := <-messages
			response, _ := strconv.Atoi(string(mes.Body))
			fmt.Printf("Received a message: %d\n", response)

			fmt.Println("received new response")

			var finalTime = time.Now()

			var rtt = int(finalTime.Sub(initialTime).Nanoseconds())
			rtts[i] = rtt

		}()
	}

	if test_flag == "true" {
		file_name := arguments[3]
		file, _ := json.MarshalIndent(rtts, "", " ")
		_ = ioutil.WriteFile(file_name, file, 0644)
	}

	wg.Wait()
}
