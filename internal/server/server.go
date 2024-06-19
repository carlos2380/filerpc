package server

import (
	"context"
	"fmt"
	"net"

	"filerpc/internal/datastore"
	"filerpc/internal/errors"
	"filerpc/internal/handler"
	log "filerpc/internal/logger"
	"filerpc/internal/proto"

	"google.golang.org/grpc"
)

func StartGRPCServer(ctx context.Context, network, port, dbAddr string) error {

	redisClient, err := datastore.InitializeRedisClient(ctx, dbAddr)
	if err != nil {
		return errors.ErrFailedToConnectRedis
	}

	log.Logger.Info("Connected to Redis")

	dstore := datastore.NewRedisFileDataStore(redisClient)
	srv := handler.NewServer(dstore)

	address := fmt.Sprintf(":%s", port)
	lis, err := net.Listen(network, address)
	if err != nil {
		return errors.ErrFailedToListen
	}

	log.Logger.Infof("Server listening on %s %s...", network, address)
	grpcServer := grpc.NewServer()
	proto.RegisterFileServiceServer(grpcServer, srv)

	if err := grpcServer.Serve(lis); err != nil {
		return errors.ErrFailedToServe
	}

	log.Logger.Info("Module initialized successfully")
	return nil
}
