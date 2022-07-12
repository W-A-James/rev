all: listener_mod/nc_listen.go dialer_mod/nc_dial.go
	go build -o build/listener listener
	go build -o build/dialer dialer

clean:
	rm -rf build/*
