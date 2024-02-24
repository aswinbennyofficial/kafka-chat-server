package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

type client struct {
	username string
	conn     *websocket.Conn
}

var clients = make(map[*websocket.Conn]*client)
var upgrader = websocket.Upgrader{}

func main() {
	http.HandleFunc("/ws", wsHandler)
	http.Handle("/", http.FileServer(http.Dir("./static")))
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

	client := &client{
		username: username,
		conn:     conn,
	}
	clients[conn] = client

	log.Printf("User '%s' connected", username)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		broadcastMessage(client, string(msg))
	}
}

func broadcastMessage(sender *client, message string) {
	for conn, client := range clients {
		if conn != sender.conn {
			client.conn.WriteMessage(websocket.TextMessage, []byte(sender.username+": "+message))
		}
	}
}
