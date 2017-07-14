package main

import (
	"io"
	"log"
	"net"
	"os"
)

func forward(conn net.Conn) {
	client, err := net.Dial("tcp", os.Args[2])
	if err != nil {
		log.Fatalf("Dial failed: %v", err)
	}
	log.Printf("Connected connection %v to remote %v\n", conn.LocalAddr(), conn.RemoteAddr())
	log.Printf("Connected remote address on %v to local %v\n", client.RemoteAddr(), client.LocalAddr())
	go func() {
		defer client.Close()
		defer conn.Close()
		io.Copy(client, conn)
	}()
	go func() {
		defer client.Close()
		defer conn.Close()
		io.Copy(conn, client)
	}()
}

func main() {
	log.SetPrefix("[main] ")
	log.Println("starting tcpforward")

	if len(os.Args) != 3 {
		log.Fatalf("Usage %s listen:port forward:port\nExample: tcpforward.exe 0.0.0.0:1234 12.113.133.192:80\n", os.Args[0])
		return
	}

	listener, err := net.Listen("tcp", os.Args[1])
	if err != nil {
		log.Fatalf("Failed to start listener: %v", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("ERROR: failed to accept listener: %v", err)
		}
		log.Printf("Accepted connection on %v\n", conn.LocalAddr())
		go forward(conn)
	}
}
