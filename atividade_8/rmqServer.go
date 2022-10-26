package main

import (
	"fmt"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

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
	//prefetch settings is done here

	defer ch.Close()

	err = ch.ExchangeDeclare(
		"pubsub", // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	q, err := ch.QueueDeclare(
		"queue", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ch.QueueBind(
		q.Name,
		"",
		"pubsub",
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	var forever chan struct{}

	for {
		msgs, ok, err := ch.Get(
			q.Name, // queue
			false,
		)
		if err != nil {
			fmt.Println(err)
			return
		}
		if ok {

			msgs.Ack(false)
			fmt.Println(msgs.MessageCount)
			time.Sleep(2000000000)
		} else {
			time.Sleep(20000000)
		}
	}

	<-forever
}
