package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

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
		Path:   "/in",
	}

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalln("pilot: failed to dial:", err)
	}

	fmt.Println("pilot: connected to server")

	inputAndSend(ws)
}

func inputAndSend(ws *websocket.Conn) {
	for {
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln("pilot: failed to read from stdin:", err)
		}

		err = ws.WriteMessage(websocket.TextMessage, []byte(input))
		if err != nil {
			log.Fatalln("pilot: failed to write message to websocket connection:", err)
			return
		}
	}
}
