package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 5 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 10 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
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
		log.Fatalln("failed to dial:", err)
	}

	log.Println("connected to server")

	inputAndSend(ws)
}

func inputAndSend(ws *websocket.Conn) {
	// from here: https://github.com/gorilla/websocket/blob/c3dd95aea9779669bb3daafbd84ee0530c8ce1c1/examples/chat/client.go
	ticker := time.NewTicker(pingPeriod)
	commandChan := make(chan struct{})

	for {
		reader := bufio.NewReader(os.Stdin)
		b, err := reader.ReadByte()
		if err != nil {
			log.Fatalln("read from stdin:", err)
		}

		go send(ws, []byte{b}, commandChan)

		select {
		case <-commandChan:
			log.Println("command sent successfully")
		case <-ticker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Fatalln("failed to write ping message")
			}
		}
	}
}

func send(ws *websocket.Conn, data []byte, c chan struct{}) error {
	defer func() {
		c <- struct{}{}
	}()

	if err := ws.WriteMessage(websocket.TextMessage, data); err != nil {
		return fmt.Errorf("write message to websocket connection: %v", err)
	}

	return nil
}
