package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleHomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func handleWebsockets(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln("failed to upgrade connection:", err)
	}

	fmt.Println("server: connection established")
	err = ws.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println("failed to write message:", err)
	}

	listen(ws)
}

func listen(conn *websocket.Conn) {
	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Fatalln(err)
			return
		}

		fmt.Println("message:", string(msg))

		if err := conn.WriteMessage(messageType, msg); err != nil {
			log.Println(err)
			return
		}
	}
}

func main() {
	http.HandleFunc("/", handleHomePage)
	http.HandleFunc("/ws", handleWebsockets)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
