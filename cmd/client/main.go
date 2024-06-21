package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	pb "filerpc/internal/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func makeGRPCRequest(client pb.FileServiceClient, version string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.FileRequest{
		Version: version,
		Hash:    "68db232fd1980a3cbdc9cc714abe9a743ef1bcd24b9d351336951f0f15ed0b63",
	}

	res, err := client.ReadFile(ctx, req)
	if err != nil {
		log.Printf("could not get file: %v", err)
		return
	}

	log.Printf("Response: Type: %s, Version: %s, Hash: %s, Content: %s", res.GetType(), res.GetVersion(), res.GetHash(), res.GetContent())
}

func main() {
	cs := flag.String("c", "1", "Number of threads to use")
	ns := flag.String("nc", "1", "Transactions per thread")
	flag.Parse()

	c, err := strconv.Atoi(*cs)
	if err != nil {
		log.Fatal(err)
	}
	n, err := strconv.Atoi(*ns)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewFileServiceClient(conn)

	wg := sync.WaitGroup{}
	t0 := time.Now()

	for i := 0; i < c; i++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()
			for j := 0; j < n; j++ {
				version := fmt.Sprintf("1.%d.%d", threadID, j)
				makeGRPCRequest(client, version)
			}
		}(i)
	}

	wg.Wait()
	tf := time.Since(t0)
	tps := float64(c) * float64(n) / tf.Seconds()
	fmt.Println("SEC:", tf.Seconds())
	fmt.Println("TPS:", tps)
}
