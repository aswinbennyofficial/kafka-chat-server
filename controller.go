package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

var clients = make(map[string]*websocket.Conn)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	defer conn.Close()

	// The client sends its username as the first message
    _, msg, err := conn.ReadMessage()
    if err != nil {
        log.Println("Failed to read username:", err)
        return
    }
    // Trim any leading and trailing whitespace from the username
    username := strings.TrimSpace(string(msg))

    // Map the username to the WebSocket connection
    clients[username] = conn
    // Log that the user has connected
    log.Printf("User '%s' connected", username)

	
	for {
		// Read a message from the client
        _, msg, err := conn.ReadMessage()
        if err != nil {
            // If the client disconnects, remove the client from the map and break the loop
            if err.Error() == "websocket: close 1001 (going away)" {
                log.Printf("User '%s' disconnected", username)
                delete(clients, username)
                break
            } else {
                log.Println("Error reading message:", err)
                break
            }
        }
        
		PublishToRedis(string(msg))
	}

}
