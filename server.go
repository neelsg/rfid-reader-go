package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var serverConnected = false
var serverConnections []*websocket.Conn = [](*websocket.Conn){}
var serverListeners []func(string) = []func(string){}

func serverConnect(port string) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		websocket, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("websocket connected!")
		serverConnections = append(serverConnections, websocket)
		serverProcess(websocket)
	})
	http.ListenAndServe(":"+port, nil)
}

func serverProcess(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		for _, v := range serverListeners {
			v(string(message))
		}
	}
}

func serverSend(data string) {
	for _, v := range serverConnections {
		if err := v.WriteMessage(websocket.TextMessage, []byte(data)); err != nil {
			log.Println(err)
		}
	}
}

func serverListen(callback func(data string)) {
	serverListeners = append(serverListeners, callback)
}
