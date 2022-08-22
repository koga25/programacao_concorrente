package main

import (
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
		fmt.Println("Please provide a host:port string")
		return
	}
	CONNECT := arguments[1]

	s, err := net.ResolveUDPAddr("udp4", CONNECT)
	c, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("The UDP server is %s\n", c.RemoteAddr().String())
	defer c.Close()

	test_flag := arguments[2]

	tries := 10000
	var rtts = make([]int, tries)

	for i := 0; i < tries; i++ {
		var initialTime = time.Now()
		data := []byte("sendTimeBetween\n")
		_, err = c.Write(data)

		if err != nil {
			fmt.Println(err)
			return
		}

		buffer := make([]byte, 1024)
		n, _, err := c.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}
		var finalTime = time.Now()

		var rtt = int(finalTime.Sub(initialTime).Nanoseconds())
		rtts[i] = rtt

		fmt.Printf((strconv.Itoa(i) + "º Reply: %s\n"), string(buffer[0:n]))
	}

	if test_flag == "true" {
		file_name := arguments[3]
		file, _ := json.MarshalIndent(rtts, "", " ")
		_ = ioutil.WriteFile(file_name, file, 0644)
	}
}
