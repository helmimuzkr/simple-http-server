package server

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strings"
)

const (
	HTTP_VERSION = "HTTP/1.1"
)

type Server struct {
	Host    string
	Port    string
	handler map[string]Handler
}

func NewServer() *Server {
	return &Server{
		Host:    "localhost",
		Port:    ":6969",
		handler: make(map[string]Handler),
	}
}

func (s *Server) Run() error {
	slog.Info("Running the server", "port", s.Port)
	l, err := net.Listen("tcp", s.Host+s.Port)
	if err != nil {
		return err
	}
	defer l.Close()

	// accept connections
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go s.processConnection(conn)
	}
}

func (s *Server) processConnection(conn net.Conn) {
	defer conn.Close()

	// read incoming request
	reader := bufio.NewReader(conn)
	r, err := ReadRequest(reader)
	if err != nil {
		slog.Error("Error during reading request", "err", err)
		return
	}

	h, ok := s.handler[r.Path]

	// validate url
	if !ok {
		slog.Error("Not found")
		response := fmt.Sprintf("%s %d %s\r\n\r\n", HTTP_VERSION, 404, "NOT FOUND")
		conn.Write([]byte(response))
		return
	}

	// validate method
	if h.method != r.Method {
		slog.Error("Method not allowed")
		response := fmt.Sprintf("%s %d %s\r\n\r\n", HTTP_VERSION, 405, "METHOD NOT ALLOWED")
		conn.Write([]byte(response))
		return
	}

	// handle the request
	slog.Info("Route the connection", "method", r.Method, "path", r.Path)
	// TODO: create wrapper for response writer
	h.handlerFunc(conn, r)
}

func (s *Server) HandleFunc(p string, h HandlerFunc) {
	keys := strings.Split(p, " ")
	if len(keys) < 2 {
		return
	}
	method := keys[0]
	path := keys[1]
	s.handler[path] = Handler{
		method:      method,
		handlerFunc: h,
	}
}

type Handler struct {
	method      string
	handlerFunc HandlerFunc
}

type HandlerFunc func(io.WriteCloser, *Request)
