package main

import (
	"fmt"
	"io"
	"log/slog"
	"simple-http-server/server"
)

func main() {
	s := server.NewServer()
	s.HandleFunc("POST /", middlewareHandler(helloHandler))
	s.Run()
}

func helloHandler(w io.WriteCloser, r *server.Request) {
	body := "Hello, world!"
	contentType := "text/plain"
	contentLength := len(body)
	response := fmt.Sprintf("HTTP/1.1 %d %s\r\nContent-Type: %s\r\nConent-Lenth: %d\r\n\r\n%s", 200, "OK", contentType, contentLength, body)
	w.Write([]byte(response))
}

func middlewareHandler(next server.HandlerFunc) server.HandlerFunc {
	return func(w io.WriteCloser, r *server.Request) {
		slog.Info("start middlewareHandler")
		next(w, r)
		slog.Info("end middlewareHandler")
	}
}
