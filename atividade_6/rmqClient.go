package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"reflect"
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

	URL := "amqp://guest:guest@localhost:" + arguments[1] + "/"

	conn, err := amqp.Dial(URL)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()
	var wg sync.WaitGroup

	for i := uint16(0); i < 1; i++ {
		wg.Add(1)
		go func(i uint16) {

			defer wg.Done()
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

			buf := [4096]byte{}
			t := reflect.TypeOf(buf)
			fmt.Println(t.Size())
			for b := 0; b < 10000; b++ {
				//encoded = binary.Write()
				binary.LittleEndian.PutUint64(buf[0:], uint64(time.Now().UTC().UnixMilli()))
				binary.LittleEndian.PutUint16(buf[8:], i)
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

			//log.Printf(" [x] Sent %s", body)

		}(i)
	}
	wg.Wait()
	/*
		test_flag := arguments[2]
			if test_flag == "true" {
				file_name := arguments[3]
				file, _ := json.MarshalIndent(rtts, "", " ")
				_ = ioutil.WriteFile(file_name, file, 0644)
			}

	*/
}
