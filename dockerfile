FROM golang:1.23.0-alpine

RUN mkdir /app/sh

WORKDIR /app/sh


ADD . /app/sh

RUN go test -v ./...

RUN go build -o main ./cmd/*.go

CMD ["/app/sh/main"]
