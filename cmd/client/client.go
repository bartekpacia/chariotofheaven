package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
)

var (
	// Host with server running to connect to.
	host string
	// Port on host to connect to, on which the server is listening.
	port string

	interval int
)

var (
	signalChan = make(chan os.Signal, 1)
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

var chariot Chariot

func init() {
	log.SetFlags(0)
	log.SetPrefix("client: ")

	flag.StringVar(&host, "host", "localhost", "host to connect to")
	flag.StringVar(&port, "port", "8080", "host's port to connect to")
	flag.IntVar(&interval, "interval", 50, "stepper motor interval")
	flag.Parse()
}

func initGPIO() {
	chip, err := gpiod.NewChip("gpiochip0", gpiod.WithConsumer("softwire"))
	if err != nil {
		log.Fatalln("failed to get chip:", err)
	}

	red, err = chip.RequestLine(rpi.GPIO17, gpiod.AsOutput(1))
	if err != nil {
		log.Fatalln("failed to request GPIO14 (red):", err)
	}

	green, err = chip.RequestLine(rpi.GPIO22, gpiod.AsOutput(1))
	if err != nil {
		log.Fatalln("failed to request GPIO15 (green):", err)
	}

	yellow, err = chip.RequestLine(rpi.GPIO27, gpiod.AsOutput(1))
	if err != nil {
		log.Fatalln("failed to request GPIO15 (yellow):", err)
	}

	servo, err = chip.RequestLine(rpi.GPIO10, gpiod.AsOutput())
	if err != nil {
		log.Fatalln("failed to request GPIO10 (servo):", err)
	}

	dir, err = chip.RequestLine(rpi.GPIO21, gpiod.AsOutput())
	if err != nil {
		log.Fatalln("failed to request GPIO21 (dir):", err)
	}

	step, err = chip.RequestLine(rpi.GPIO20, gpiod.AsOutput())
	if err != nil {
		log.Fatalln("failed to request GPIO20 (step):", err)
	}
}

func initBlink() {
	time.Sleep(time.Millisecond * 500)
	setActivePins(red, green, yellow)
	time.Sleep(time.Millisecond * 500)
	setActivePins()
	time.Sleep(time.Millisecond * 500)
	setActivePins(red, green, yellow)
	time.Sleep(time.Millisecond * 500)
	setActivePins()
}

func main() {
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		setActivePins()
		log.Println("shutting down...")
		os.Exit(0)
	}()

	initGPIO()
	initBlink()

	u := url.URL{
		Scheme: "ws",
		Host:   host + ":" + port,
		Path:   "/out",
	}

	ws, err := connect(u)
	if err != nil {
		log.Fatalln("failed to connect to server:", err)
	}
	log.Println("connected to server")

	chariot = Chariot{
		MovingState:      NotMoving,
		TurningDirection: Left,
		Turning:          false,
	}

	go startTurner()

	listenWebsockets(u, ws)
}

func connect(u url.URL) (ws *websocket.Conn, err error) {
	ws, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("dial: %v", err)
	}

	return ws, nil
}

func listenWebsockets(u url.URL, ws *websocket.Conn) {
	for {
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("failed to read message from websocket connection")

			for {
				time.Sleep(3 * time.Second)
				log.Println("attempting to reconnect to server")
				ws, err = connect(u)
				if err != nil {
					continue
				}
				break
			}
			continue
		}

		if msgType != websocket.TextMessage {
			log.Fatalln("received message of type other than TextMessage")
		}

		log.Printf("received command: %#v\n", string(msg))

		chariot.InterpretCommand(string(msg))
		execute(&chariot)
	}
}

func execute(c *Chariot) {
	switch c.MovingState {
	case MovingForward:
		setActivePins(green)

	case MovingBackward:
		setActivePins(red)

	case NotMoving:
		setActivePins(yellow)
	}

	switch c.TurningDirection {
	case Left:
		dir.SetValue(0)
	case Right:
		dir.SetValue(1)
	}
}

func startTurner() {
	for {
		if chariot.Turning {
			step.SetValue(1)
			time.Sleep(time.Millisecond * time.Duration(interval))
			step.SetValue(0)
			time.Sleep(time.Millisecond * time.Duration(interval))
		} else {
			step.SetValue(0)
		}
	}
}

func setActivePins(pins ...*gpiod.Line) {
	red.SetValue(0)
	green.SetValue(0)
	yellow.SetValue(0)

	for _, pin := range pins {
		if pin != nil {
			err := pin.SetValue(1)
			if err != nil {
				log.Fatalf("client: failed to set pin %v to 1\n", pin)
			}
		}
	}
}
