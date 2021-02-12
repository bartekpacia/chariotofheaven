package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleWebsockets(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln("failed to upgrade connection:", err)
	}

	fmt.Println("connection established")
	err = ws.WriteMessage(1, []byte("hello to client from server!"))
	if err != nil {
		log.Println("failed to write message:", err)
	}

	inputAndSend(ws)
}

func inputAndSend(conn *websocket.Conn) {
	for {
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln("failed to read from stdin:", err)
		}

		err = conn.WriteMessage(websocket.TextMessage, []byte(input))
		if err != nil {
			log.Fatalln("failed to write msg:", err)
			return
		}
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "you've connected to the chariot server")
	})

	http.HandleFunc("/ws", handleWebsockets)

	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
