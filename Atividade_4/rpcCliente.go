package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/rpc"
	"os"
	"strconv"
	"time"
)

type RemainingTime struct {
	// Seconds until the event starts
	Seconds int
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide host:port.")
		os.Exit(0)
	}

	fmt.Println("Starting TCP connection at " + arguments[1])

	CONNECT := arguments[1]
	c, err := rpc.DialHTTP("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	var reply RemainingTime

	test_flag := arguments[2]

	tries := 10000
	var rtts = make([]int, tries)
	for i := 0; i < tries; i++ {
		var initialTime = time.Now()

		c.Call("Handler.GetRemainingTime", "SendRemainingTime\n", &reply)

		var finalTime = time.Now()
		var rtt = int(finalTime.Sub(initialTime).Nanoseconds())
		rtts[i] = rtt

		fmt.Println(strconv.Itoa(i) + " ->: " + strconv.Itoa(reply.Seconds))
	}

	if test_flag == "true" {
		file_name := arguments[3]
		file, _ := json.MarshalIndent(rtts, "", " ")
		_ = ioutil.WriteFile(file_name, file, 0644)
	}
}
