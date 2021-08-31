package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"

	"github.com/stianeikeland/go-rpio/v4"
)

var (
	host     string
	port     string
	interval int
)

var signalChan = make(chan os.Signal)

// Pins used to communicate with physical parts.
var (
	red    rpio.Pin
	green  rpio.Pin
	yellow rpio.Pin
	servo  rpio.Pin
	dir    rpio.Pin
	step   rpio.Pin
)

var chariot ChariotState

func init() {
	log.SetFlags(0)
	log.SetPrefix("client: ")

	flag.StringVar(&host, "host", "localhost", "host to connect to")
	flag.StringVar(&port, "port", "8080", "host's port to connect to")
	flag.IntVar(&interval, "interval", 1, "stepper motor interval (ms)")
	flag.Parse()

	err := rpio.Open()
	if err != nil {
		log.Fatalln("failed to initialize GPIO:", err)
	}

	red = rpio.Pin(4)
	red.Output()

	yellow = rpio.Pin(17)
	yellow.Output()

	green = rpio.Pin(22)
	green.Output()

	servo = rpio.Pin(10)
	servo.Output()

	dir = rpio.Pin(21)
	dir.Output()

	step = rpio.Pin(20)
	step.Output()

	time.Sleep(time.Millisecond * 500)
	setPinsHigh(&red, &green, &yellow)
	time.Sleep(time.Millisecond * 500)
	setPinsLow(&red, &green, &yellow)
	time.Sleep(time.Millisecond * 500)
	setPinsHigh(&red, &green, &yellow)
	time.Sleep(time.Millisecond * 500)
	setPinsLow(&red, &green, &yellow)
}

func main() {
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		setPinsLow(&red, &green, &yellow, &servo, &dir, &step)
		log.Println("clean up and shutdown")
		os.Exit(0)
	}()

	u := url.URL{
		Scheme: "ws",
		Host:   host + ":" + port,
		Path:   "/out",
	}

	ws, err := connect(u)
	if err != nil {
		log.Fatalln("failed to connect to server:", err)
	}
	log.Printf("connected to server at %s on port %s\n", host, port)

	chariot = ChariotState{
		MovingState:  NotMoving,
		TurningState: NotTurning,
	}

	go startTurner()

	listenWebsockets(u, ws)
}

func connect(u url.URL) (ws *websocket.Conn, err error) {
	ws, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
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

		cmd := string(msg)
		chariot.MapCommandToState(cmd)
		startMover(&chariot)
	}
}

func startMover(c *ChariotState) {
	switch c.MovingState {
	case MovingForward:
		setPinsLow(&red, &green, &yellow)
		setPinsHigh(&green)

	case MovingBackward:
		setPinsLow(&red, &green, &yellow)
		setPinsHigh(&red)

	case NotMoving:
		setPinsLow(&red, &green, &yellow)
		setPinsHigh(&yellow)
	}
}

func startTurner() {
	for {
		switch chariot.TurningState {
		case TurningLeft:
			dir.Low()
			step.High()
			time.Sleep(time.Millisecond * time.Duration(interval))
			step.Low()
			time.Sleep(time.Millisecond * time.Duration(interval))
		case TurningRight:
			dir.High()
			step.High()
			time.Sleep(time.Millisecond * time.Duration(interval))
			step.Low()
			time.Sleep(time.Millisecond * time.Duration(interval))
		case NotTurning:
			step.Low()
		}
	}
}

// setPinsHigh sets pins to high.
func setPinsHigh(pins ...*rpio.Pin) {
	for _, pin := range pins {
		if pin != nil {
			pin.High()
		}
	}
}

// setPinsLow pins to low.
func setPinsLow(pins ...*rpio.Pin) {
	for _, pin := range pins {
		if pin != nil {
			pin.Low()
		}
	}
}
