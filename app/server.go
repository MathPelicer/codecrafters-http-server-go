package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

type Request struct {
	HttpMethod    string
	RequestTarget string
	HttpVersion   string
	Headers       map[string]string
}

func ParseRequest(request []string) *Request {
	firstLine := strings.Split(request[0], " ")

	headers := make(map[string]string)
	for _, element := range request[1:] {
		if strings.Contains(element, "\x00") {
			break
		}

		if element != "" {
			headerSplit := strings.Split(element, ":")
			headers[headerSplit[0]] = strings.TrimSpace(headerSplit[1])
			continue
		}
	}

	return &Request{
		HttpMethod:    firstLine[0],
		RequestTarget: firstLine[1],
		HttpVersion:   firstLine[2],
		Headers:       headers,
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

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

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		if err != io.EOF {
			os.Exit(1)
		}
	}

	splittedBuffer := strings.Split(string(buffer), "\r\n")
	request := ParseRequest(splittedBuffer)

	if request.RequestTarget == "/" {
		response := "HTTP/1.1 200 OK\r\n\r\n"
		conn.Write([]byte(response))
	}

	if strings.Contains(request.RequestTarget, "/echo/") {
		targetResources := strings.Split(request.RequestTarget, "/")
		finalResourse := targetResources[len(targetResources)-1]

		response := constructTextResponse(len(finalResourse), finalResourse)

		conn.Write([]byte(response))
	}

	if request.RequestTarget == "/user-agent" {
		response := constructTextResponse(len(request.Headers["User-Agent"]),
			request.Headers["User-Agent"])

		conn.Write([]byte(response))
	}

	conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
}

func constructTextResponse(contentLen int, content string) string {
	response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",
		contentLen,
		content)

	return response
}
