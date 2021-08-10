package main

import (
	"bufio"
	"flag"
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
	log.SetPrefix("pilot: ")

	flag.StringVar(&host, "host", "localhost", "host to connect to")
	flag.StringVar(&port, "port", "8080", "host' port to connect to")
	flag.Parse()
}

func main() {
	u := url.URL{
		Scheme: "ws",
		Host:   host + ":" + port,
		Path:   "/in",
	}

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalln("failed to connect to server:", err)
	}
	log.Printf("connected to server at %s on port %s\n", host, port)

	inputAndSend(ws)
}

func inputAndSend(ws *websocket.Conn) {
	for {
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadByte()
		if err != nil {
			log.Fatalln("failed to read from stdin:", err)
		}

		err = ws.WriteMessage(websocket.TextMessage, []byte{input})
		if err != nil {
			log.Fatalln("failed to write message to websocket connection:", err)
			return
		}
	}
}
