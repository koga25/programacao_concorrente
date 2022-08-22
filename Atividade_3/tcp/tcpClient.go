package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide host:port.")
		os.Exit(0)
	}

	fmt.Println("Starting TCP connection at " + arguments[1])

	CONNECT := arguments[1]
	c, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	var rttTotal = 0
	var tries = 10000
	for i := 0; i < tries; i++ {
		var initialTime = time.Now()
		fmt.Fprintf(c, "sendTimeBetween\n")

		message, _ := bufio.NewReader(c).ReadString('\n')
		var finalTime = time.Now()

		var rtt = int(finalTime.Sub(initialTime).Milliseconds())
		rttTotal = rttTotal + rtt

		if message == "END" {
			fmt.Println("Server asked to disconnect. TCP client exiting...")
			os.Exit(0)
		}

		fmt.Println(strconv.Itoa(i) + " ->: " + message)
	}

	c.Write([]byte("END"))

	rttTotal = rttTotal / tries
	os.Exit(rttTotal)
}