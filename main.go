package main

import (
	"fmt"
	"os"
)

var mainDebug = false

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "debug" {
			fmt.Println("debug enabled")
			mainDebug = true
		}
	}

	fmt.Println("starting service...")
	serverPort, readerAddress, readerPort := configGet()
	if mainDebug {
		fmt.Println("serverPort:", serverPort, " readerAddress:", readerAddress, " readerPort:", readerPort)
	}

	readerConnect(readerAddress, readerPort)
	readerListen(serverSend)

	serverConnect(serverPort)
	serverListen(func(data string) {
		if mainDebug {
			fmt.Println("server received:", data)
		}
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
