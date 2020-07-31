package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
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
			if _, err := ws.Write([]byte(fmt.Sprintf("[FROM CLIENT] %d time!", i))); err != nil {
				log.Fatal(err)
			}

			time.Sleep(time.Second * 5)
			i++
		}
	}()

	// Read message until user quits
	for {
		msg := make([]byte, 512)
		if n, err := ws.Read(msg); err != nil {
			log.Fatal(err)
		} else {
			log.Printf("Received: %s.\n", msg[:n])
		}
	}
}
