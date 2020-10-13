/**
 * Copyright 2020 David Hyland-Wood
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *   http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	serverport = ":8082"
)

var (
	clients   = make(map[*websocket.Conn]bool) // connected clients
	broadcast = make(chan Message)             // broadcast channel
	upgrader  = websocket.Upgrader{}           // Upgrades HTTP requests to WebSockets
)

// Message object
type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

// String makes Message satisfy the Stringer interface.
func (a Message) String() string {
	return fmt.Sprintf("From %v:  %v", a.Username, a.Message)
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	clients[ws] = true

	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)

		// Echo to server's STDOUT
		fmt.Println("Received by Server: ", msg)

		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast

		// Echo to server's STDOUT
		fmt.Println("About to send to clients: ", msg)

		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	// Create a simple file server
	fs := http.FileServer(http.Dir("/Users/davidhyland-wood/Documents/GitHub/go-chat/public"))
	http.Handle("/", fs)

	// Configure websocket route
	http.HandleFunc("/ws", handleConnections)

	// Listen for incoming chat messages
	go handleMessages()

	// Start the server on localhost the designated port and log any errors
	log.Println("http server started on " + serverport)
	err := http.ListenAndServe(serverport, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	// Create a simple Web server
	// type countHandler struct {
	// 	mu sync.Mutex // guards n
	// 	n  int
	// }

	// func (h *countHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 	h.mu.Lock()
	// 	defer h.mu.Unlock()
	// 	h.n++
	// 	fmt.Fprintf(w, "count is %d\n", h.n)
	// }

	// func main() {
	// 	fmt.Println("Chat server started on port", serverport)

	// 	// Start server
	// 	http.Handle("/count", new(countHandler))
	// 	log.Fatal(http.ListenAndServe(serverport, nil))

	//Get test response
	// resp, err := http.Get("http://example.com/")
	// if err != nil {
	// 	// handle error
	// }
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// readableBody := string(body[:])
	// fmt.Println(readableBody)
}
