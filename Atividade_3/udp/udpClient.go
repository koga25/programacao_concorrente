package main

import (
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
	s, err := net.ResolveUDPAddr("udp4", CONNECT)
	c, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()

	var receivingExpression = ""

	for {
		fmt.Fprintf(c, "sendNextLetter\n")

		buffer := make([]byte, 1024)
		n, _, err := c.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}
		var message = string(buffer[0:n])
		
		if message == "END" {
			fmt.Println("Message received fully. UDP client exiting...")
			return
		}

		var nextLetter = string(message[0])

		receivingExpression = receivingExpression + nextLetter
		fmt.Println("->: " + receivingExpression)
	}
}