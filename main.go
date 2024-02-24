package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)



var clients = make(map[string]*websocket.Conn)
var upgrader = websocket.Upgrader{}

func main() {
	http.HandleFunc("/ws", wsHandler)
	http.Handle("/", http.FileServer(http.Dir("./static/")))
	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	defer conn.Close()

	// Get the username from the client
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Failed to read username:", err)
		return
	}
	username := strings.TrimSpace(string(msg))

	
	clients[username] = conn

	log.Printf("User '%s' connected", username)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if err.Error() == "websocket: close 1001 (going away)" {
				log.Printf("User '%s' disconnected", username)
				delete(clients, username)
				break
			} else {
				log.Println("Error reading message:", err)
				break
			}
		}
		broadcastMessage(username, string(msg))
	}
}

func broadcastMessage(sender string, message string) {
	for username, conn := range clients {
		if  username != sender{
			conn.WriteMessage(websocket.TextMessage, []byte(sender+": "+message))
			// conn.WriteMessage(websocket.TextMessage, []byte(sender.username+": "+message))
		}
	}
}
