package main

import (
        "fmt"
        "net"
        "os"
        "strconv"
        "time"
)

var eventTime = time.Date(2025, 07, 29, 14, 30, 45, 100, time.Local)

func main() {
        arguments := os.Args
        if len(arguments) == 1 {
                fmt.Println("Please provide a port number!")
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
        buffer := make([]byte, 1024)

        for {
                n, addr, err := connection.ReadFromUDP(buffer)
                fmt.Print("-> ", string(buffer[0:n-1]))

				var today = time.Now()
				var timeBetween = int(eventTime.Sub(today).Seconds())
				

                data := []byte(strconv.Itoa(timeBetween) + "\n")
                fmt.Printf("data: %s\n", string(data))
                _, err = connection.WriteToUDP(data, addr)
                if err != nil {
                        fmt.Println(err)
                        return
                }
        }
}