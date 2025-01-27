package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"testing"
	"time"
)

const (
	host = "localhost"
	port = ":6969"
)

func helloHandler(w io.WriteCloser, r *Request) {
	body := "Hello, world!"
	contentType := "text/plain"
	contentLength := len(body)
	response := fmt.Sprintf("HTTP/1.1 %d %s\r\nContent-Type: %s\r\nConent-Lenth: %d\r\n\r\n%s", 200, "OK", contentType, contentLength, body)
	w.Write([]byte(response))
}

func TestServer(t *testing.T) {
	// run the server
	s := NewServer()
	s.Host = host
	s.Port = port
	s.HandleFunc("GET /", helloHandler)
	go func() {
		if err := s.Run(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	time.Sleep(100 * time.Millisecond) // waiting for the server start

	// create test case
	testCase := createTestCase()

	// run test
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			// connect to the server
			conn, err := net.Dial("tcp", host+port)
			if err != nil {
				t.Fatalf("Failed to connect to server: %v", err)
			}
			defer conn.Close()

			// send request
			conn.Write([]byte(tc.request))

			// read the response
			reader := bufio.NewReader(conn)
			requestLine, err := reader.ReadString('\n')
			if err != nil {
				t.Fatalf("failed to read request line: %v", err)
			}

			// validate status
			if !strings.Contains(requestLine, tc.expectedStatus) {
				t.Errorf("expected status %q, got %q", tc.expectedStatus, requestLine)
			}
		})
	}
}

type testCase struct {
	name           string
	request        string
	expectedStatus string
}

func createTestCase() []testCase {
	return []testCase{
		{
			name:           "Valid GET Request",
			request:        "GET / HTTP:/1.1\r\n\r\n",
			expectedStatus: "200 OK",
		},
		{
			name:           "Not found GET Request",
			request:        "GET /users HTTP:/1.1\r\n\r\n",
			expectedStatus: "404 NOT FOUND",
		},
	}
}
