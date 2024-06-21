FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go run createFiles.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go-rpc-server ./cmd/server/main.go

FROM scratch

COPY --from=builder /go-rpc-server /go-rpc-server
COPY --from=builder /app/core /core

EXPOSE 50051
EXPOSE 8080

CMD ["/go-rpc-server", "-network", "tcp", "-grpc-port", "50051", "-host", "127.0.0.1", "-gateway-port", "8080"]
