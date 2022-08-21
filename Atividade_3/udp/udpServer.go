package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

var sendingExpression = "socorram me subi no onibus em marrocos"
var sendingExpressionLength = len(sendingExpression)
var progressTracker map[int]int
var clientIdCounter = 0

func handleConnection(connection *net.UDPConn, quit chan struct{}) {
	buffer := make([]byte, 1024)
	n, remoteAddr, err := 0, new(net.UDPAddr), error(nil)

	clientIdCounter++
	var id = clientIdCounter
	progressTracker[id] = 0
	fmt.Println("New client connection established. id: " + strconv.Itoa(id))

	for err == nil {
		n, remoteAddr, err = connection.ReadFromUDP(buffer)
		var message = string(buffer[0:n])

		fmt.Println(message + " from id: " + strconv.Itoa(id))

		var currentIndex = progressTracker[id]

		if currentIndex == sendingExpressionLength {
			connection.WriteToUDP([]byte("END"), remoteAddr)
			fmt.Println("Finished sending expression to id: " + strconv.Itoa(id))
			break
		}

		var nextLetter = string(sendingExpression[currentIndex]) + "\n"
		connection.WriteToUDP([]byte(nextLetter), remoteAddr)

		progressTracker[id] = currentIndex + 1
	}

	fmt.Println("A client disconnected. id: " + strconv.Itoa(id))
	quit <- struct{}{}
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide port number")
		return
	}

	PORT := ":" + arguments[1]
	s, err := net.ResolveUDPAddr("udp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}

	connection, err := net.ListenUDP("udp4", s)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer connection.Close()

	progressTracker = make(map[int]int)

	max_listeners := 5
	quit := make(chan struct{})

	for i := 0; i < max_listeners; i++ {
		go handleConnection(connection, quit)
	}
	<-quit // Waits for an error to happen.
}