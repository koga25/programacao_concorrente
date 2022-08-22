package main

import (
        "fmt"
        "net"
        "os"
		"strconv"
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
		
		var tries = 10000
        for i := 0; i < tries; i++ {
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
                fmt.Printf((strconv.Itoa(i) + "º Reply: %s\n"), string(buffer[0:n]))
        }
}
    