package main

import (
	"io"
	"log"
	"net"
)

func main() {
	// Listen for incoming connections
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	// Accept incoming connections
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		// Connect to the other server or device
		otherConn, err := net.Dial("tcp", "other-server:8080")
		if err != nil {
			log.Println(err)
			conn.Close()
			continue
		}

		// Copy data between the two connections
		go func() {
			_, err := io.Copy(conn, otherConn)
			if err != nil {
				log.Println(err)
			}
			conn.Close()
			otherConn.Close()
		}()
		go func() {
			_, err := io.Copy(otherConn, conn)
			if err != nil {
				log.Println(err)
			}
			conn.Close()
			otherConn.Close()
		}()
	}
}