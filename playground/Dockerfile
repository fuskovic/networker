FROM golang:latest

RUN mkdir -p /app/tls

WORKDIR /app

COPY server.go ./server.go

CMD go run server.go