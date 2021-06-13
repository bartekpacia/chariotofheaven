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
	log.SetFlags(0)
	log.SetPrefix("server: ")

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
		log.Fatalln("failed to upgrade INPUT to websocket connection:", err)
	}

	log.Println("pilot connected")

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("pilot disconnected")
			break
		}

		commands <- string(msg)
		log.Printf("received command: %s\n", msg)
	}
}

func handleOutWebsockets(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln("failed to upgrade OUTPUT to websocket connection:", err)
	}

	log.Println("client connected")

	for {
		event := <-commands
		log.Printf("sent command: %#v\n", event)
		err = ws.WriteMessage(websocket.TextMessage, []byte(event))
		if err != nil {
			log.Fatalln("failed to write message to /out websocket connection:", err)
		}
	}
}
