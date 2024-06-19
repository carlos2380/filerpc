package main

import (
	"net"

	"filerpc/internal/handler"
	pb "filerpc/internal/proto"

	log "filerpc/internal/logger"

	"google.golang.org/grpc"
)

func main() {
	log.Logger.Info("Initializing module...")
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Logger.Fatalf("failed to listen: %v", err)
	}
	log.Logger.Info("Server listening port 50051...")
	s := grpc.NewServer()
	pb.RegisterFileServiceServer(s, &handler.Server{})
	if err := s.Serve(lis); err != nil {
		log.Logger.Fatalf("failed to serve: %v", err)
	}
	log.Logger.Info("Module initialized successfully")
}
