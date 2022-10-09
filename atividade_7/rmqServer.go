package main

import (
	"encoding/binary"
	"encoding/json"
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
	ch.Qos(10000, 0, false)

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
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
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

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	var forever chan struct{}
	i := 0
	buf := [10000]int64{}

	for d := range msgs {
		unixTimestamp := int64(binary.LittleEndian.Uint64(d.Body[0:]))
		d.Ack(false)
		timeElapsed := time.Now().UTC().UnixMilli() - unixTimestamp
		buf[i] = int64(timeElapsed)
		i++
		if i == 9999 {
			var x = int64(0)
			for z := 0; z < 10000; z++ {
				x += int64(buf[z])
			}

			x = int64(x / 10000)
			fmt.Println(x)
			i = 0

			file, _ := json.MarshalIndent(x, "", " ")
			f, err := os.OpenFile("test_json_10000_PC.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println(err)
			}
			f.Write(file)
			f.WriteString("\n")
		}
	}

	<-forever
}
