package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

var traveller = []byte("I")[0]

func main() {
	isReceiver := os.Getenv("IS_RECEIVER") == "true"
	destination := os.Getenv("DESTINATION")
	port := os.Getenv("LISTEN_PORT")

	transmitter, err := connectRemote(destination)
	if err != nil {
		panic(err)
	}
	defer transmitter.Close()

	if isReceiver {
		finishLine := make(chan time.Time)
		go startReceiver(port, func(received byte) {
			if received == traveller {
				finishLine <- time.Now()
			} else {
				fmt.Println("Hitchiker!", received)
			}
		})()

		fmt.Println("Sending message")
		start := time.Now()
		_, err := transmitter.Write([]byte{traveller})
		if err != nil {
			panic(err)
		}

		fmt.Println("Awaiting arrival")
		end := <- finishLine
		elapsed := end.Sub(start)

		fmt.Println("Round trip took", elapsed.String())
	} else {
		startReceiver(port, func(received byte) {
			_, err := transmitter.Write([]byte{received})
			if err != nil {
				panic(err)
			}
			fmt.Println("Forwarded a traveller")
		})()
	}
}

func startReceiver(port string, callback func(byte)) func() {
	address, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp", address)
	return func() {
		buf := make([]byte, 1)

		for {
			fmt.Println("Listening for messages")
			_, _, err := conn.ReadFromUDP(buf)
			if err != nil {
				panic(err)
			}

			callback(buf[0])
		}
	}
}

func connectRemote(destination string) (*net.UDPConn, error) {
	destinationAddress, err := net.ResolveUDPAddr("udp", destination)
	if err != nil {
		return nil, err
	}

	return net.DialUDP("udp", nil, destinationAddress)
}