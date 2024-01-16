package main

import (
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/sigurn/crc16"
)

var readerConnected = false
var readerConnection net.Conn
var readerListeners []func(string) = []func(string){}

func readerConnect(address, port string) {
	if readerConnected {
		return
	}

	var err error
	readerConnection, err = net.Dial("tcp", address+":"+port)
	if err != nil {
		fmt.Printf("fail to connect to reader: %v", err)
		os.Exit(1)
	}
	readerConnected = true
	time.AfterFunc(100*time.Millisecond, func() { readerSend("c000") })

	go readerProcess()
}

func readerDisconnect() {
	if !readerConnected {
		return
	}

	readerListeners = []func(string){}
	readerConnection.Close()
	readerConnected = false
}

func readerProcessRead(buffer []byte, mLen int, err error) {
	if err != nil {
		fmt.Println("error reading: ", err.Error())
		return
	}

	if mLen < 8 {
		fmt.Println("data err: too short: ", hex.EncodeToString(buffer[:mLen]))
		return
	}

	if buffer[0] != 0xAA || buffer[1] != 0xAA || buffer[2] != 0xFF {
		fmt.Println("data err: did not start with 0xAAAAFF: ", hex.EncodeToString(buffer[:mLen]))
		return
	}

	if (mLen - 3) > int(buffer[3]) {
		readerProcessRead(buffer[0:buffer[3]+3], int(buffer[3]+3), nil)
		readerProcessRead(buffer[buffer[3]+3:], mLen-int(buffer[3]+3), nil)
		return
	}

	if mLen-3 != int(buffer[3]) {
		fmt.Println("data err: length did not match: ", hex.EncodeToString(buffer))
		return
	}

	table := crc16.MakeTable(crc16.CRC16_CCITT_FALSE)
	crc := crc16.Checksum(buffer[0:mLen-2], table)
	if buffer[mLen-2] != byte(crc>>8) || buffer[mLen-1] != byte(crc%256) {
		fmt.Println("data err: crc did not match: ", hex.EncodeToString(buffer))
		return
	}

	message := "unknown message from reader: " + hex.EncodeToString(buffer)
	if buffer[4] == 0xC0 && buffer[5] == 0x00 && buffer[6] == 0x00 {
		message = "stopped"
	} else if buffer[4] == 0xC1 && buffer[5] == 0x02 && buffer[6] == 0x00 {
		message = "tag " + hex.EncodeToString(buffer[10:mLen-5])
	}

	for _, v := range readerListeners {
		v(message)
	}
}

func readerProcess() {
	for {
		buffer := make([]byte, 1024)
		mLen, err := readerConnection.Read(buffer)
		readerProcessRead(buffer, mLen, err)
	}
}

func readerListen(callback func(data string)) {
	readerListeners = append(readerListeners, callback)
}

func readerSend(data string) {
	hexData, err := hex.DecodeString(data)
	if err != nil {
		fmt.Println("error encoding: ", err.Error())
		return
	}

	buffer := []byte{0xAA, 0xAA, 0xFF, byte(len(hexData) + 3)}
	buffer = append(buffer, hexData...)
	table := crc16.MakeTable(crc16.CRC16_CCITT_FALSE)
	crc := crc16.Checksum(buffer, table)
	buffer = append(buffer, byte(crc>>8), byte(crc%256))
	readerConnection.Write(buffer)
}
