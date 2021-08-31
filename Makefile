all: client server pilot

CLIENT := ./cmd/client
SERVER := ./cmd/server
PILOT := ./cmd/pilot

client: $(CLIENT)/client.go $(CLIENT)/commands.go $(CLIENT)/states.go
	go build $(CLIENT)/client.go $(CLIENT)/commands.go $(CLIENT)/states.go

server: $(SERVER)/server.go $(SERVER)/in.go $(SERVER)/out.go
	go build $(SERVER)/server.go $(SERVER)/in.go $(SERVER)/out.go

pilot: $(PILOT)/pilot.go
	go build $(PILOT)/pilot.go

clean:
	rm -f ./client ./server ./pilot
