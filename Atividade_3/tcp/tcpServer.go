package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"strconv"
	"time"
)

var concurrentClientsCounter = 0
var clientIdCounter = 0 

var eventTime = time.Date(2025, 07, 29, 14, 30, 45, 100, time.Local)

func handleConnection(c net.Conn) {
	clientIdCounter++
	var id = clientIdCounter
	fmt.Println("New client connection established. id: " + strconv.Itoa(id))

	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		// Holds the message received from client, but its not really useful.
		message := strings.TrimSpace(string(netData))

		if message == "END" {
			c.Write([]byte("END"))
			fmt.Println("Client id: " + strconv.Itoa(id) + " asked to disconnect.")
			break
		}

		var today = time.Now()
		var timeBetween = int(eventTime.Sub(today).Seconds())
		
		c.Write([]byte(strconv.Itoa(timeBetween) + "\n"))
	}

	concurrentClientsCounter--
	fmt.Println("A client disconnected. id: " + strconv.Itoa(id))
	c.Close()
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide port number")
		return
	}

	PORT := ":" + arguments[1]
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
		concurrentClientsCounter++
	}
}