package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"os"
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

	URL := "amqp://guest:guest@localhost:" + arguments[1] + "/"

	conn, err := amqp.Dial(URL)

	if err != nil {
		fmt.Println(err)
		return
	}
	var forever chan struct{}

	defer conn.Close()
	ch, err := conn.Channel()
	defer ch.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	buf := [8]byte{}

	for {
		//encoded = binary.Write()
		binary.LittleEndian.PutUint64(buf[0:], uint64(time.Now().UTC().UnixMilli()))
		err = ch.PublishWithContext(ctx,
			"pubsub", // exchange
			"",       // routing key
			false,    // mandatory
			false,    // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(buf[:]),
			})
		if err != nil {
			fmt.Println(err)
			return
		}

		count := get_message_count(ch)
		fmt.Println(count)
		if count == 5 {
			for {
				if get_message_count(ch) == 0 {
					break
				}
			}
		}
		time.Sleep(800000000)
	}

	<-forever
}

func get_message_count(ch *amqp.Channel) int {
	q, err := ch.QueueInspect(
		"queue", // name
	)

	if err != nil {
		fmt.Println(err)
	}
	return q.Messages
}
