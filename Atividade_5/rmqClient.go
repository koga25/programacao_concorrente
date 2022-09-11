package main

import (
	// "bufio"
	// "encoding/json"
	"fmt"
	// "io/ioutil"
	// "net"
	"os"
	// "strconv"
	// "time"
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

	conn, err := amqp.Dial("amqp://guest:guest@" + CONNECT + "/")
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

	queue, err := ch.QueueDeclare(
		"timeBetween", // name
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

	msgs, err := ch.Consume(
		queue.Name, // queue
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
	
	go func() {
	for d := range msgs {
		fmt.Println("Received a message: %s", string(d.Body))
	}
	}()
	
	fmt.Println(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

	// test_flag := arguments[2]

	// tries := 10000
	// var rtts = make([]int, tries)
	// for i := 0; i < tries; i++ {
	// 	var initialTime = time.Now()
	// 	fmt.Fprintf(c, "sendTimeBetween\n")

	// 	message, _ := bufio.NewReader(c).ReadString('\n')
	// 	var finalTime = time.Now()

	// 	var rtt = int(finalTime.Sub(initialTime).Nanoseconds())
	// 	rtts[i] = rtt

	// 	if message == "END" {
	// 		fmt.Println("Server asked to disconnect. TCP client exiting...")
	// 		os.Exit(0)
	// 	}

	// 	fmt.Println(strconv.Itoa(i) + " ->: " + message)
	// }

	// c.Write([]byte("END"))

	// if test_flag == "true" {
	// 	file_name := arguments[3]
	// 	file, _ := json.MarshalIndent(rtts, "", " ")
	// 	_ = ioutil.WriteFile(file_name, file, 0644)
	// }
}
