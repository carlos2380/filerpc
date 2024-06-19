package main

import (
	"context"
	"log"
	"time"

	pb "filerpc/internal/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Dial the gRPC server
	conn, err := grpc.DialContext(ctx, "localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create a new gRPC client
	client := pb.NewFileServiceClient(conn)

	// Prepare the request
	req := &pb.FileRequest{
		Hash: "",
	}

	// Call the ReadFile method
	res, err := client.ReadFile(ctx, req)
	if err != nil {
		log.Fatalf("could not get file: %v", err)
	}

	// Print the response
	log.Printf("Response: Type: %s, Version: %s, Hash: %s, Content: %s", res.GetType(), res.GetVersion(), res.GetHash(), res.GetContent())
}
