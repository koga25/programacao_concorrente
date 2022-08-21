package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"strconv"
)

var concurrentClientsCounter = 0

var sendingExpression = "socorram me subi no onibus em marrocos"
var sendingExpressionLength = len(sendingExpression)
var progressTracker map[int]int
var clientIdCounter = 0

func handleConnection(c net.Conn) {
	clientIdCounter++
	var id = clientIdCounter
	progressTracker[id] = 0
	fmt.Println("New client connection established. id: " + strconv.Itoa(id))

	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		// Holds the message received from client, but its not really useful.
		temp := strings.TrimSpace(string(netData))
		fmt.Println(temp + " from id: " + strconv.Itoa(id))

		var currentIndex = progressTracker[id]

		if currentIndex == sendingExpressionLength {
			c.Write([]byte("END"))
			fmt.Println("Finished sending expression to id: " + strconv.Itoa(id))
			break
		}

		var nextLetter = string(sendingExpression[currentIndex]) + "\n"
		c.Write([]byte(nextLetter))

		progressTracker[id] = currentIndex + 1
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

	progressTracker = make(map[int]int)

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