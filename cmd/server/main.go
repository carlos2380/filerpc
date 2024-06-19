package main

import (
	"flag"
	"net"

	"filerpc/internal/handler"
	pb "filerpc/internal/proto"

	log "filerpc/internal/logger"

	"google.golang.org/grpc"
)

func main() {

	network := flag.String("network", "tcp", "Network type to use (e.g., tcp, tcp4, tcp6, unix)")
	port := flag.String("port", "50051", "Port or address to listen on")
	flag.Parse()

	log.Logger.Info("Initializing module...")
	startGRPCServer(*network, *port)
}

func startGRPCServer(network, port string) {
	address := net.JoinHostPort("", port)
	lis, err := net.Listen(network, address)
	if err != nil {
		log.Logger.Fatalf("failed to listen: %v", err)
	}

	log.Logger.Infof("Server listening on %s %s...", network, address)
	s := grpc.NewServer()
	pb.RegisterFileServiceServer(s, &handler.Server{})
	if err := s.Serve(lis); err != nil {
		log.Logger.Fatalf("failed to serve: %v", err)
	}

}
