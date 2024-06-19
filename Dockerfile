FROM golang:1.19 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go-rpc-server ./cmd/server/main.go

FROM scratch

COPY --from=builder /go-rpc-server /go-rpc-server
COPY --from=builder /app/core /core

EXPOSE 50051

CMD ["/go-rpc-server", "-network", "tcp", "-port", "50051"]