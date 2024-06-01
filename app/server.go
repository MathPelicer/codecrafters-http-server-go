package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

const (
	FILE_PATH = "/tmp/data/codecrafters.io/http-server-tester/"
	STATUS404 = "HTTP/1.1 404 Not Found\r\n\r\n"
	STATUS201 = "HTTP/1.1 201 Created\r\n\r\n"
	STATUS200 = "HTTP/1.1 200 OK\r\n\r\n"
)

type Request struct {
	HttpMethod    string
	RequestTarget string
	HttpVersion   string
	Headers       map[string]string
	Body          string
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
			headers[strings.ToLower(headerSplit[0])] = strings.TrimSpace(headerSplit[1])
			continue
		}
	}

	return &Request{
		HttpMethod:    firstLine[0],
		RequestTarget: firstLine[1],
		HttpVersion:   firstLine[2],
		Headers:       headers,
		Body:          strings.Trim(request[len(request)-1], "\x00"),
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

	switch request.HttpMethod {
	case "GET":
		handleGetMethod(request, conn)
	case "POST":
		handlePostMethod(request, conn)
	}

	conn.Write([]byte(STATUS404))
}

func handlePostMethod(request *Request, conn net.Conn) {
	if strings.HasPrefix(request.RequestTarget, "/files/") {
		pathParam := getFinalPathParam(request.RequestTarget)
		writeFileContent(pathParam, request.Body)

		conn.Write([]byte(STATUS201))
	}
}

func handleGetMethod(request *Request, conn net.Conn) {
	var response string

	if request.RequestTarget == "/" {
		conn.Write([]byte(STATUS200))
	}

	if strings.HasPrefix(request.RequestTarget, "/echo/") {
		pathParam := getFinalPathParam(request.RequestTarget)
		encodingType, isValid := isEncodingTypeValid(request)
		if isValid {
			response = constructEncodedTextResponse(pathParam, encodingType)
			conn.Write([]byte(response))
			return
		}

		response = constructTextResponse(pathParam)
		conn.Write([]byte(response))
		return
	}

	if strings.HasPrefix(request.RequestTarget, "/files/") {
		pathParam := getFinalPathParam(request.RequestTarget)
		finalResourse, err := readFileContent(pathParam)

		if err != nil {
			conn.Write([]byte(STATUS404))
			return
		}

		response := constructOctetStreamResponse(finalResourse)
		conn.Write([]byte(response))
		return
	}

	if request.RequestTarget == "/user-agent" {
		response := constructTextResponse(request.Headers["user-agent"])

		conn.Write([]byte(response))
		return
	}
}

func isEncodingTypeValid(request *Request) (encodingType string, isValidEncoding bool) {
	encodingType, isEncoded := request.Headers["accept-encoding"]
	if isEncoded &&
		encodingType == "gzip" {
		return encodingType, true
	}

	return "", false
}

func getFinalPathParam(requestTarget string) string {
	targetResources := strings.Split(requestTarget, "/")
	return targetResources[len(targetResources)-1]
}

func constructTextResponse(content string) string {
	response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",
		len(content),
		content)

	return response
}

func constructEncodedTextResponse(content string, encoding string) string {
	response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Encoding: %s\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",
		encoding,
		len(content),
		content)

	return response
}

func constructOctetStreamResponse(content string) string {
	response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s",
		len(content),
		content)

	return response
}

func readFileContent(filename string) (fileContent string, err error) {
	filePath := FILE_PATH + filename
	readFileContent, err := os.ReadFile(filePath)

	if err != nil {
		return
	}

	fileContent = string(readFileContent)
	return
}

func writeFileContent(filename string, fileContent string) {
	filePath := FILE_PATH + filename
	file, err := os.Create(filePath)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	_, err = file.WriteString(fileContent)

	if err != nil {
		panic(err)
	}
}
