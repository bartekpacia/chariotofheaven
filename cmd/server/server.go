package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var port string

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	commands = make(chan string)
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("server: ")

	flag.StringVar(&port, "port", "8080", "port to listen on")
	flag.Parse()
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "you've connected to the chariot server")
	})

	http.HandleFunc("/in", handleInWebsockets)
	http.HandleFunc("/out", handleOutWebsockets)

	log.Printf("started on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleInWebsockets(w http.ResponseWriter, req *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatalln("failed to upgrade INPUT to websocket connection:", err)
	}

	log.Println("pilot connected")

	// Receive commands from pilot and send them to client in a blocking manner
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("pilot disconnected")
			break
		}

		log.Printf("received command: %s\n", msg)
		commands <- string(msg)
		log.Printf("sent command: %s\n", msg)
	}
}

func handleOutWebsockets(w http.ResponseWriter, req *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := upgrader.Upgrade(w, req, nil)
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
