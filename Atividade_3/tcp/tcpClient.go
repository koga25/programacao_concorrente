package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	test_flag := arguments[2]

	tries := 10000
	var rtts = make([]int, tries)
	for i := 0; i < tries; i++ {
		var initialTime = time.Now()
		fmt.Fprintf(c, "sendTimeBetween\n")

		message, _ := bufio.NewReader(c).ReadString('\n')
		var finalTime = time.Now()

		var rtt = int(finalTime.Sub(initialTime).Nanoseconds())
		rtts[i] = rtt

		if message == "END" {
			fmt.Println("Server asked to disconnect. TCP client exiting...")
			os.Exit(0)
		}

		fmt.Println(strconv.Itoa(i) + " ->: " + message)
	}

	c.Write([]byte("END"))

	if test_flag == "true" {
		file_name := arguments[3]
		file, _ := json.MarshalIndent(rtts, "", " ")
		_ = ioutil.WriteFile(file_name, file, 0644)
	}
}
