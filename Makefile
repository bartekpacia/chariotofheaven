all: client server

CLIENT := ./cmd/client
SERVER := ./cmd/server

client: $(CLIENT)/client.go
	go build $(CLIENT)/client.go

server: $(SERVER)/server.go
	go build $(SERVER)/server.go

clean:
	rm -f ./client ./server
