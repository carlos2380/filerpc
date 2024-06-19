package main

import (
	"context"
	"flag"
	"fmt"
	"net"

	"filerpc/internal/handler"
	pb "filerpc/internal/proto"

	log "filerpc/internal/logger"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
)

var redisClient *redis.Client

func main() {

	network := flag.String("network", "tcp", "Network type to use (e.g., tcp, tcp4, tcp6, unix)")
	port := flag.String("port", "50051", "Port or address to listen on")
	redisAddr := flag.String("redis-addr", "redis:6379", "Address of the Redis server")
	flag.Parse()

	redisClient := redis.NewClient(&redis.Options{
		Addr: *redisAddr,
	})

	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Logger.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Logger.Info("Connected to Redis")

	startGRPCServer(*network, *port, redisClient)
}

func startGRPCServer(network, port string, redisClient *redis.Client) {
	log.Logger.Info("Initializing module...")
	address := fmt.Sprintf(":%s", port)
	lis, err := net.Listen(network, address)
	if err != nil {
		log.Logger.Fatalf("failed to listen: %v", err)
	}

	log.Logger.Infof("Server listening on %s %s...", network, address)
	s := grpc.NewServer()
	pb.RegisterFileServiceServer(s, handler.NewServer(redisClient))
	if err := s.Serve(lis); err != nil {
		log.Logger.Fatalf("failed to serve: %v", err)
	}

	log.Logger.Info("Module initialized successfully")
}
