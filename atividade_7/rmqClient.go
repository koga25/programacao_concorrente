package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"reflect"
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	buf := [1024]byte{}
	t := reflect.TypeOf(buf)
	fmt.Println(t.Size())
	for b := 0; b < 10000; b++ {
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
	}

}
