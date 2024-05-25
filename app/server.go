package main

import (
	"fmt"
	"io"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

type RequestLine struct {
	HttpMethod    string
	RequestTarget string
	HttpVersion   string
}

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
		requestLine := constructRequestLine(strings.Split(splittedBuffer[0], " "))

		if requestLine.RequestTarget == "/" {
			response := "HTTP/1.1 200 OK\r\n\r\n"
			conn.Write([]byte(response))
			continue
		}

		if strings.Contains(requestLine.RequestTarget, "/echo/") {
			targetResources := strings.Split(requestLine.RequestTarget, "/")
			finalResourse := targetResources[len(targetResources)-1]

			response := fmt.Sprintf(`HTTP/1.1 200 OK\r\n
				Content-Type: text/plain\r\n
				Content-Length: %d\r\n\r\n
				%s`,
				len(finalResourse),
				finalResourse)

			conn.Write([]byte(response))
			continue
		}

		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}

func constructRequestLine(requestLine []string) RequestLine {
	requestLineStruct := &RequestLine{
		HttpMethod:    requestLine[0],
		RequestTarget: requestLine[1],
		HttpVersion:   requestLine[2],
	}

	return *requestLineStruct
}
