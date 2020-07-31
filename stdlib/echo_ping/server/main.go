package main

import (
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
)

func main() {
	// Create a server that does not check origin
	server := websocket.Server{Handshake: nil, Handler: Handler}

	// Register it with the default mux
	http.Handle("/echo-ping", server)

	// Start the server
	log.Fatal(http.ListenAndServe("0.0.0.0:7777", nil))
}

// Our websocket handler will echo any client message  and randomly pings the client by itself
func Handler(ws *websocket.Conn) {
	if _, err := io.Copy(ws, ws); err != nil {
		ws.Close()
	}
}
