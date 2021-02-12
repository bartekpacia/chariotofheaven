package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/pkg/errors"

	"github.com/gorilla/websocket"
)

var (
	host      string
	port      string
	serverURL url.URL
)

func init() {
	log.SetFlags(0)
	flag.StringVar(&host, "host", "localhost", "host to connect to")
	flag.StringVar(&port, "port", "8080", "host's port to connect to")
}

func main() {
	serverURL = url.URL{
		Scheme: "ws",
		Host:   host + ":" + port,
		Path:   "/out",
	}

	ws, err := connect(serverURL)
	if err != nil {
		log.Fatalln("client: failed to connect to server:", err)
	}
	fmt.Println("client: connected to server")

	listenWebsockets(ws)
}

func connect(u url.URL) (ws *websocket.Conn, err error) {
	ws, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return ws, nil
}

func listenWebsockets(ws *websocket.Conn) {
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("client: failed to read message from websocket connection")

			for {
				time.Sleep(3 * time.Second)
				fmt.Println("client: attempting to reconnect to server")
				ws, err = connect(serverURL)
				if err != nil {
					continue
				}
				break
			}
			continue
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
