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
	host string
	port string
)

var (
	red    *gpiod.Line
	green  *gpiod.Line
	yellow *gpiod.Line
	servo  *gpiod.Line
	dir    *gpiod.Line
	step   *gpiod.Line
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

	red, err = chip.RequestLine(rpi.GPIO17, gpiod.AsOutput())
	if err != nil {
		log.Fatalln("client: failed to request GPIO14 (red):", err)
	}

	green, err = chip.RequestLine(rpi.GPIO22, gpiod.AsOutput())
	if err != nil {
		log.Fatalln("client: failed to request GPIO15 (green):", err)
	}

	yellow, err = chip.RequestLine(rpi.GPIO27, gpiod.AsOutput())
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

	serverURL := url.URL{
		Scheme: "ws",
		Host:   host + ":" + port,
		Path:   "/out",
	}

	ws, err := connect(serverURL)
	if err != nil {
		log.Fatalln("client: failed to connect to server:", err)
	}
	fmt.Println("client: connected to server")

	listenWebsockets(serverURL, ws)
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
		_, msg, err := ws.ReadMessage()
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
		fmt.Printf("client: received command: %#v\n", string(msg))

		matchCommand(string(msg))
	}
}

func matchCommand(command string) {
	switch command {
	case "w":
		resetMovePins()
		green.SetValue(1)
	case "b":
		resetMovePins()
		red.SetValue(1)
	case "a":
		dir.SetValue(0)
		go startStepping()
		// start turning left
	case "d":
		dir.SetValue(1)
		go startStepping()
		// start turning right
	case "z":

	case "x":
		resetMovePins()
		yellow.SetValue(1)
	default:
		fmt.Printf("command %s not matched\n", command)
	}
}

func startStepping() {
	for {
		step.SetValue(1)
		time.After(1 * time.Second)
		step.SetValue(1)
	}

}

func resetMovePins() {
	red.SetValue(0)
	green.SetValue(0)
	yellow.SetValue(0)
}
