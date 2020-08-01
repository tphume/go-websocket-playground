package main

import (
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{}

func main() {
	http.HandleFunc("/echo-ping", Handler)
	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe("localhost:7777", nil))
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// Upgrade connection to web socket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("on upgrade:", err)
		return
	}
	defer conn.Close()

	// Randomly send ping message to client
	go func() {
		defer conn.Close()

		for {
			time.Sleep(time.Duration(rand.Intn(10)+5) * time.Second)
			if err := conn.WriteMessage(websocket.PingMessage, []byte("ping from server na ja")); err != nil {
				log.Println("on ping:", err)
				break
			}

			log.Println("ping sent to client")
		}
	}()

	// Echo back messages sent from the client
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("on read:", err)
			break
		}

		log.Println("received:", string(message))
		if err := conn.WriteMessage(mt, message); err != nil {
			log.Println("write:", err)
			break
		}
	}
}
