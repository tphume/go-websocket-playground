package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"math/rand"
	"time"
)

const (
	URL    = "ws://localhost:7777/echo-ping"
	ORIGIN = "http://localhost/"
)

func main() {
	// Connect to the server
	ws, err := websocket.Dial(URL, "", ORIGIN)
	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()

	// Send message periodically
	go func() {
		i := 1
		for {
			log.Println("sending message:", i)
			if _, err := ws.Write([]byte(fmt.Sprintf("[FROM CLIENT] %d time!", i))); err != nil {
				log.Fatal(err)
			}

			time.Sleep(time.Duration(rand.Intn(10) + 5))
			i++
		}
	}()

	// Read message until user quits
	for {
		reader, err := ws.NewFrameReader()
		if err != nil {
			log.Fatal("on read:", err)
		}

		switch reader.PayloadType() {
		case websocket.PingFrame:
			log.Println("ping received from server")

			pong, err := ws.NewFrameWriter(websocket.PongFrame)
			if err != nil {
				log.Println("on pong writer:", err)
			}

			if _, err := pong.Write([]byte("pong message from client")); err != nil {
				log.Println("on pong:", err)
			}
		default:
			msg := make([]byte, 512)
			if _, err := reader.Read(msg); err != nil {
				log.Println("on receive:", err)
				break
			}

			log.Println("received:", string(msg))
		}
	}
}
