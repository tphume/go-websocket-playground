package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	// Create a server that does not check origin
	server := websocket.Server{Handshake: nil, Handler: Handler}

	// Register it with the default mux
	http.Handle("/echo-ping", server)

	// Start the server
	log.Fatal(http.ListenAndServe("localhost:7777", nil))
}

func Handler(ws *websocket.Conn) {
	defer ws.Close()

	// Sends message periodically
	go func() {
		i := 1
		for {
			if _, err := ws.Write([]byte(fmt.Sprintf("[FROM SERVER] %d time!", i))); err != nil {
				log.Println("Connection Closed by client")
			}

			time.Sleep(time.Second * 10)
			i++
		}
	}()

	//Read message infinitely
	if _, err := io.Copy(ws, ws); err != nil {
		return
	}
}
