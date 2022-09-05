package main

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"
)

var eventTime = time.Date(2025, 07, 29, 14, 30, 45, 100, time.Local)

type RemainingTime struct {
	// Seconds until the event starts
	Seconds int
}

// Handler is the struct which exposes the RemainingTime Server methods
type Handler struct {
}

// New returns the object for the RPC handler
func New() *Handler {
	h := &Handler{}
	err := rpc.Register(h)
	if err != nil {
		fmt.Println(err)
	}
	return h
}

// getRemainingTime function seconds until the event starts
func (rh *Handler) GetRemainingTime(payload string, reply *RemainingTime) error {
	var rt RemainingTime
	var today = time.Now()
	rt.Seconds = int(eventTime.Sub(today).Seconds())

	*reply = rt
	return nil
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide port number")
		return
	}

	PORT := ":" + arguments[1]

	var handler = new(Handler)
	err := rpc.Register(handler)
	if err != nil {
		fmt.Println(err)
		return
	}

	rpc.HandleHTTP()

	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = http.Serve(l, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
