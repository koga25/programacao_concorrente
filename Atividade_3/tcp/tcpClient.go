package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
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

	var tries = 10000
	for i := 0; i < tries; i++ {
		fmt.Fprintf(c, "sendTimeBetween\n")

		message, _ := bufio.NewReader(c).ReadString('\n')
		
		if message == "END" {
			fmt.Println("Server asked to disconnect. TCP client exiting...")
			return
		}

		fmt.Println(strconv.Itoa(i) + " ->: " + message)
	}

	c.Write([]byte("END"))

}