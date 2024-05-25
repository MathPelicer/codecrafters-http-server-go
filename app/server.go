package main

import (
	"fmt"
	"io"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		buffer := make([]byte, 1024)
		_, err = conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				os.Exit(1)
			}
		}

		splittedBuffer := strings.Split(string(buffer), "\r\n")
		requestLine := strings.Split(splittedBuffer[0], " ")

		if requestLine[1] == "/" {
			conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
			continue
		}

		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}
