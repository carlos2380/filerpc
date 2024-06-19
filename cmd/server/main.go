package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "filerpc/internal/logger"
	"filerpc/internal/server"
)

func main() {
	network := flag.String("network", "tcp", "Network type to use (e.g., tcp, tcp4, tcp6, unix)")
	port := flag.String("port", "50051", "Port or address to listen on")
	dbAddr := flag.String("redis-addr", "redis:6379", "Address of the Redis server")
	flag.Parse()

	log.Logger.Info("Initializing module...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.StartGRPCServer(ctx, *network, *port, *dbAddr); err != nil {
			log.Logger.Fatalf("failed to start server: %v", err)
		}
	}()

	<-sig
	log.Logger.Info("Shutting down server...")
	cancel()

	time.Sleep(2 * time.Second)
	log.Logger.Info("Server stopped")
}
