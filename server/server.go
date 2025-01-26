package server

import (
	"bufio"
	"io"
	"log/slog"
	"net"
	"strings"
)

type Server struct {
	Port    string
	handler map[string]Handler
}

func NewServer() *Server {
	return &Server{Port: "6969", handler: make(map[string]Handler)}
}

func (s *Server) AddHandler(p string, h HandlerFunc) {
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

func (s *Server) Run() error {
	slog.Info("Running the server", "port", s.Port)
	l, err := net.Listen("tcp", ":"+s.Port)
	if err != nil {
		return err
	}
	defer l.Close()

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

	// read request
	reader := bufio.NewReader(conn)
	r, err := ReadRequest(reader)
	if err != nil {
		slog.Error("Error during reading request", "err", err)
		return
	}

	// route the handler
	h, ok := s.handler[r.Path]
	if !ok {
		slog.Error("URL not found")
		// TODO: handle error to client
		return
	}
	if h.method != r.Method {
		slog.Error("Method not allowed")
		// TODO: handle error to client
		return
	}
	slog.Info("Route connection", "path", r.Path)
	h.handlerFunc.Serve(conn, r)
}

type Handler struct {
	method      string
	handlerFunc HandlerFunc
}

type HandlerFunc func(io.WriteCloser, *Request)

func (fn HandlerFunc) Serve(w io.WriteCloser, r *Request) {
	fn(w, r)
}
