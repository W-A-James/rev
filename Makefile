all: server/server.go client/client.go
	go build -o build/rev_server server
	go build -o build/rev_client client

clean:
	rm -rf build/*
