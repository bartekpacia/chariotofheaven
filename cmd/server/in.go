package main

import (
	"log"
	"net/http"
)

func handleInWebsockets(w http.ResponseWriter, req *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatalln("failed to upgrade /in to websocket connection:", err)
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
