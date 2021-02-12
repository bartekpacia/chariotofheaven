package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

var (
	host string
	port string
)

func init() {
	log.SetFlags(0)
	flag.StringVar(&host, "host", "localhost", "host to connect to")
	flag.StringVar(&port, "port", "8080", "host's port to connect to")
}

func main() {
	u := url.URL{
		Scheme: "ws",
		Host:   host + ":" + port,
		Path:   "/out",
	}

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalln("client: failed to dial:", err)
	}

	fmt.Println("client: connected to server")

	listenWebsockets(ws)
}

func listenWebsockets(ws *websocket.Conn) {
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Fatalln("client: failed to read message from websocket connection:", err)
		}

		matchCommand(string(msg))

		fmt.Printf("client: received command: %#v\n", string(msg))
	}
}

func matchCommand(command string) {
	switch command {
	case "w":
		// start moving forward
	case "b":
		// start moving backward
	case "a":
		// start turning left
	case "d":
		// start turning right
	case "z":
		// stop turning
	case "x":
		// stop moving and turning
	}
}
