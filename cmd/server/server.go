package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	commands = make(chan string)
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "you've connected to the chariot server")
	})

	http.HandleFunc("/in", handleInWebsockets)
	http.HandleFunc("/out", handleOutWebsockets)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleInWebsockets(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln("server: failed to upgrade INPUT to websocket connection:", err)
	}

	fmt.Println("server: pilot connected")

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Fatalln("server: failed to read message from INPUT websocket connection:", err)
		}

		commands <- string(msg)
		fmt.Printf("server: received command: %#v\n", string(msg))
	}
}

func handleOutWebsockets(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln("server: failed to upgrade OUTPUT to websocket connection:", err)
	}

	fmt.Println("server: client connected")

	for {
		select {
		case event := <-commands:
			fmt.Printf("server: sent command: %#v\n", event)
			err = ws.WriteMessage(websocket.BinaryMessage, []byte(event))
			if err != nil {
				log.Println("server: failed to write message to /out websocket connection:", err)
			}
		}
	}
}
