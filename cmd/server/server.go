package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var port string

var upgrader = websocket.Upgrader{}

var commands = make(chan string)

func init() {
	log.SetFlags(0)
	log.SetPrefix("server: ")

	flag.StringVar(&port, "port", "8080", "port to listen on")
	flag.Parse()
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "you've connected to the chariot server")
	})

	http.HandleFunc("/in", handleInWebsockets)
	http.HandleFunc("/out", handleOutWebsockets)

	log.Printf("started on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
