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
		Path:   "/ws",
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalln("failed to dial:", err)
	}

	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Fatalln("failed to read msg from conn:", err)
	}

	fmt.Println("got message:", string(msg))
}
