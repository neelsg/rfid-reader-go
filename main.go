package main

import (
	"fmt"
)

func main() {
	fmt.Println("starting service...")
	serverPort, readerAddress, readerPort := configGet()

	readerConnect(readerAddress, readerPort)
	readerListen(serverSend)

	serverConnect(serverPort)
	serverListen(func(data string) {
		if data == "start" {
			readerSend("c102050000")
			serverSend("starting")
		} else if data == "stop" {
			readerSend("c000")
		} else if data == "reconnect" {
			readerDisconnect()
			readerConnect(readerAddress, readerPort)
			serverSend("reconnecting")
		} else {
			serverSend("unknown command")
		}
	})
}
