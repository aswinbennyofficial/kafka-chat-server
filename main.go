package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)


// Create a map where key is username and value is a pointer to a WebSocket connection
var clients = make(map[string]*websocket.Conn)
// Specify the parameters to be used for upgrading HTTP connections to WebSocket connections
var upgrader = websocket.Upgrader{}

/*
main is the entry point of the program. It creates a WebSocket endpoint and a file server to serve static files. It starts the server on port 8080.
*/

func main() {
	// Create a WebSocket endpoint
	http.HandleFunc("/ws", wsHandler)

	/* Create a file server to serve static files
	 Whenever a request is made to the root URL, the server will serve the index.html file
	 from the static directory. 
	 */
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})
	
	// Start the server on port 8080
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

/*
wsHandler is a function that handles WebSocket connections. It reads the username from the client, allows the client to send messages, and broadcasts the messages to all other clients.
*/
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

	// Continuously read messages from the client and broadcast them to all other clients
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
		// Broadcast the message to all other clients
		broadcastMessage(username, string(msg))
	}
}

/*
broadcastMessage sends a message to all clients except the sender.
*/
func broadcastMessage(sender string, message string) {
	// Iterate over all clients
	for username, conn := range clients {
		// If the client is not the sender, send the message to the client
		if  username != sender{
			conn.WriteMessage(websocket.TextMessage, []byte(sender+": "+message))
		}
	}
}
