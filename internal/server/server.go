package server

import (
	"context"
	"fmt"
	"net"

	"filerpc/internal/datastore"
	"filerpc/internal/handler"
	log "filerpc/internal/logger"
	"filerpc/internal/proto"

	"google.golang.org/grpc"
)

func StartGRPCServer(ctx context.Context, network, port, dbAddr string) error {

	redisClient, err := datastore.InitializeRedisClient(ctx, dbAddr)
	if err != nil {
		return fmt.Errorf("failed to initialize Redis client: %w", err)
	}

	log.Logger.Info("Connected to Redis")

	dstore := datastore.NewRedisFileDataStore(redisClient)
	srv := handler.NewServer(dstore)

	address := fmt.Sprintf(":%s", port)
	lis, err := net.Listen(network, address)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	log.Logger.Infof("Server listening on %s %s...", network, address)
	grpcServer := grpc.NewServer()
	proto.RegisterFileServiceServer(grpcServer, srv)

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	log.Logger.Info("Module initialized successfully")
	return nil
}
