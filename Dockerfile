FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go run createFiles.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go-rpc-server ./cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /client ./cmd/client/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /doc-server ./cmd/doc/main.go

FROM scratch as server

COPY --from=builder /go-rpc-server /go-rpc-server
COPY --from=builder /app/core /core

EXPOSE 50051

CMD ["/go-rpc-server", "-network", "tcp", "-grpc-port", "50051"]

FROM scratch AS client
COPY --from=builder /client /client
CMD ["/client", "-c", "1", "-nc", "1"]


FROM scratch AS docserver
COPY --from=builder /doc-server /doc-server
COPY --from=builder /app/doc /doc

EXPOSE 8080

CMD ["/doc-server", "-port-doc", "8080"]