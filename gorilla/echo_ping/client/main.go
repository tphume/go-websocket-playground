package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"time"
)

const URL = "ws://localhost:7777/echo-ping"

func main() {
	log.Println("Connecting to server...")
	conn, _, err := websocket.DefaultDialer.Dial(URL, nil)
	if err != nil {
		log.Fatal("on dial:", err)
	}
	defer conn.Close()

	client := &Client{Conn: conn}

	// Register a pong handler
	conn.SetPingHandler(client.HandlePong)

	// Send a message to server every 5 seconds
	go func() {
		for {
			i := 1
			for {
				err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Message numero: %d", i)))
				if err != nil {
					conn.Close()
					log.Fatal("on send:", err)
				}

				log.Printf("message numero %d sent to server\n", i)
				time.Sleep(time.Second * 5)
				i++
			}
		}
	}()

	// Read echo message
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Fatal("on receive:", err)
		}

		log.Println("received:", string(message))
	}
}

type Client struct {
	Conn *websocket.Conn
}

func (c *Client) HandlePong(appData string) error {
	log.Println("received ping from server")
	err := c.Conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(time.Second))
	if err == websocket.ErrCloseSent {
		return nil
	} else if e, ok := err.(net.Error); ok && e.Temporary() {
		return nil
	}
	return err
}
