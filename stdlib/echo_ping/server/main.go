package main

import (
	"golang.org/x/net/websocket"
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
	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe("localhost:7777", nil))
}

func Handler(ws *websocket.Conn) {
	defer ws.Close()

	// Sends ping message periodically
	go func() {
		for {
			time.Sleep(time.Duration(rand.Intn(10)+5) * time.Second)

			// FrameWriter to send ping message with correct opcode
			ping, err := ws.NewFrameWriter(websocket.PingFrame)
			if err != nil {
				log.Println("on ping writer:", err)
				return
			}

			if _, err = ping.Write([]byte("ping message from server")); err != nil {
				log.Println("on ping:", err)
				break
			}
			log.Println("ping message sent")

			_ = ping.Close()
		}
	}()

	//Read message infinitely
	for {
		reader, err := ws.NewFrameReader()
		if err != nil {
			log.Fatal("on read:", err)
		}

		// Read
		msg := make([]byte, 512)
		if _, err := reader.Read(msg); err != nil {
			log.Println("on receive:", err)
			break
		}

		switch reader.PayloadType() {
		case websocket.PongFrame:
			log.Println("pong received from client")
		default:
			if _, err := ws.Write(msg); err != nil {
				log.Println("on echo:", err)
			}

			log.Println("echo message back")
		}
	}
}
