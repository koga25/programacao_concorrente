package main

import (
	"encoding/binary"
	"encoding/json"
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
		true,   // auto-ack
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
		publisherNumber := binary.LittleEndian.Uint16(d.Body[8:])
		if publisherNumber == 20 {
			timeElapsed := time.Now().UTC().UnixMicro() - unixTimestamp
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
				f, err := os.OpenFile("test_json_1024.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					fmt.Println(err)
				}
				f.Write(file)
				f.WriteString("\n")
			}
		}
		//i++
		//fmt.Printf("received %s.       %dnth message.\n", d.Body, i)
	}

	<-forever
}
