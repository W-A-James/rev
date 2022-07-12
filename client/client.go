package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	host := flag.String("host", "", "Host IP")
	port := flag.String("port", "", "Host Port")
	flag.Parse()

	if *host == "" {
		fmt.Println("Must specify target host")
		flag.Usage()
		os.Exit(1)
	}

	if *port == "" {
		fmt.Println("Must specify target port")
		flag.Usage()
		os.Exit(1)
	}

	addr := *host + ":" + *port
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	go io.Copy(os.Stdout, conn)
	io.Copy(conn, os.Stdin)
	conn.Close()
}
