package main

import (
	"golang.org/x/net/websocket"
	"io"
	"log"
	"math/rand"
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

	// Sends ping message periodically
	go func() {
		// FrameWriter to send ping message with correct opcode
		ping, err := ws.NewFrameWriter(websocket.PingFrame)
		if err != nil {
			return
		}

		for {
			if _, err = ping.Write([]byte("ping message from server")); err != nil {
				log.Println("on ping:", err)
				break
			}

			_ = ping.Close()
			time.Sleep(time.Duration(rand.Intn(10) + 5))
		}
	}()

	//Read message infinitely
	if _, err := io.Copy(ws, ws); err != nil {
		return
	}
}
