package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide host:port.")
		return
	}

	fmt.Println("Starting TCP connection at " + arguments[1])

	CONNECT := arguments[1]
	c, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}

	var receivingExpression = ""

	for {
		fmt.Fprintf(c, "sendNextLetter\n")

		message, _ := bufio.NewReader(c).ReadString('\n')
		
		if message == "END" {
			fmt.Println("Message received fully. TCP client exiting...")
			return
		}

		var nextLetter = string(message[0])

		receivingExpression = receivingExpression + nextLetter
		fmt.Println("->: " + receivingExpression)
	}
}