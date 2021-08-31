package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func handleOutWebsockets(w http.ResponseWriter, req *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatalln("failed to upgrade /out to websocket connection:", err)
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
