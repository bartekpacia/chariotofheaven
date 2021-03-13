all: client server pilot

CLIENT := ./cmd/client
SERVER := ./cmd/server
PILOT := ./cmd/pilot

client: $(CLIENT)/client.go
	go build $(CLIENT)/client.go $(CLIENT)/commands.go $(CLIENT)/state.go

server: $(SERVER)/server.go
	go build $(SERVER)/server.go

pilot: $(PILOT)/pilot.go
	go build $(PILOT)/pilot.go

clean:
	rm -f ./client ./server ./pilot
