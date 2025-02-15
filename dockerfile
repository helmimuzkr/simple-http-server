FROM golang:1.23.0-alpine

# when creating subdirectories hanging off from a non-existing parent directory(s) you must pass the -p flag to mkdir
RUN mkdir -p /app/shs

WORKDIR /app/shs

ADD . /app/shs

RUN go test -v ./...

RUN go build -o main ./cmd/*.go

CMD ["/app/shs/main"]
