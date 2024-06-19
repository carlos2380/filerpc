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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, "localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewFileServiceClient(conn)

	req := &pb.FileRequest{
		Hash: "",
	}

	res, err := client.ReadFile(ctx, req)
	if err != nil {
		log.Fatalf("could not get file: %v", err)
	}

	log.Printf("Response: Type: %s, Version: %s, Hash: %s, Content: %s", res.GetType(), res.GetVersion(), res.GetHash(), res.GetContent())
}
