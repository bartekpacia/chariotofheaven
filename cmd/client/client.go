package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/pkg/errors"

	"github.com/gorilla/websocket"
	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
)

var (
	// Host with server running to connect to.
	host string
	// Port on host to connect to, on which the server is listening.
	port string
)

var (
	commands = make(chan string)
	turnChan = make(chan struct{})
)

// Pins used to communicate with physical parts.
var (
	red    *gpiod.Line
	green  *gpiod.Line
	yellow *gpiod.Line
	servo  *gpiod.Line
	dir    *gpiod.Line
	step   *gpiod.Line
)

const (
	MoveForward  = "w"
	MoveBackward = "s"
	MoveStop     = "q"

	TurnLeft  = "a"
	TurnRight = "d"
	TurnStop  = "z"

	StopAll = "x"
)

func init() {
	log.SetFlags(0)
	flag.StringVar(&host, "host", "localhost", "host to connect to")
	flag.StringVar(&port, "port", "8080", "host's port to connect to")
}

func initGPIO() {
	chip, err := gpiod.NewChip("gpiochip0", gpiod.WithConsumer("softwire"))
	if err != nil {
		log.Fatalln("client: failed to get chip:", err)
	}

	red, err = chip.RequestLine(rpi.GPIO17, gpiod.AsOutput(1))
	if err != nil {
		log.Fatalln("client: failed to request GPIO14 (red):", err)
	}

	green, err = chip.RequestLine(rpi.GPIO22, gpiod.AsOutput(1))
	if err != nil {
		log.Fatalln("client: failed to request GPIO15 (green):", err)
	}

	yellow, err = chip.RequestLine(rpi.GPIO27, gpiod.AsOutput(1))
	if err != nil {
		log.Fatalln("client: failed to request GPIO15 (yellow):", err)
	}

	servo, err = chip.RequestLine(rpi.GPIO10, gpiod.AsOutput())
	if err != nil {
		log.Fatalln("client: failed to request GPIO10 (servo):", err)
	}

	dir, err = chip.RequestLine(rpi.GPIO21, gpiod.AsOutput())
	if err != nil {
		log.Fatalln("client: failed to request GPIO21 (dir):", err)
	}

	step, err = chip.RequestLine(rpi.GPIO20, gpiod.AsOutput())
	if err != nil {
		log.Fatalln("client: failed to request GPIO20 (step):", err)
	}
}

func initBlink() {
	time.After(time.Second * 1)

	red.SetValue(1)
	green.SetValue(1)
	yellow.SetValue(1)

	time.After(time.Second * 1)

	red.SetValue(0)
	green.SetValue(0)
	yellow.SetValue(0)
}

func main() {
	flag.Parse()
	initGPIO()
	initBlink()

	u := url.URL{
		Scheme: "ws",
		Host:   host + ":" + port,
		Path:   "/out",
	}

	ws, err := connect(u)
	if err != nil {
		log.Fatalln("client: failed to connect to server:", err)
	}
	fmt.Println("client: connected to server")

	go processCommands()
	listenWebsockets(u, ws)
}

func connect(u url.URL) (ws *websocket.Conn, err error) {
	ws, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return ws, nil
}

func listenWebsockets(u url.URL, ws *websocket.Conn) {
	for {
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("client: failed to read message from websocket connection")

			for {
				time.Sleep(3 * time.Second)
				fmt.Println("client: attempting to reconnect to server")
				ws, err = connect(u)
				if err != nil {
					continue
				}
				break
			}
			continue
		}

		if msgType != websocket.TextMessage {
			log.Fatalln("client: received message of type other than TextMessage")
		}

		fmt.Printf("client: received command: %#v\n", string(msg))
		commands <- string(msg)
	}
}

func processCommands() {
	for {
		command := <-commands

		switch command {
		case MoveForward:
			resetMovePins()
			green.SetValue(1)

		case MoveBackward:
			resetMovePins()
			red.SetValue(1)

		case MoveStop:
			resetMovePins()
			yellow.SetValue(1)

		case TurnLeft:
			dir.SetValue(0)
			//turnChan <- struct{}{}
			go startStepping()

		case TurnRight:
			dir.SetValue(1)
			//turnChan <- struct{}{}
			go startStepping()

		case TurnStop:
			//turnChan <- struct{}{}

		case StopAll:
			resetMovePins()
			//turnChan <- struct{}{}
			yellow.SetValue(1)

		default:
			fmt.Printf("command %s not matched\n", command)
		}
	}
}

func startStepping() {
	for {
		select {
		case <-turnChan:

		case <-time.After(time.Millisecond * 500):
			step.SetValue(1)
			time.After(time.Millisecond * 500)
			step.SetValue(1)
		}
	}

}

func resetMovePins() {
	red.SetValue(0)
	green.SetValue(0)
	yellow.SetValue(0)
}
