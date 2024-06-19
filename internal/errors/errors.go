package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrFileNotFound         = status.New(codes.NotFound, "file not found").Err()
	ErrReadFile             = status.New(codes.Internal, "error on read file").Err()
	ErrFailedToConnectRedis = status.New(codes.Unavailable, "failed to connect to Redis").Err()
	ErrFailedToListen       = status.New(codes.Internal, "failed to listen").Err()
	ErrFailedToServe        = status.New(codes.Internal, "failed to serve").Err()
)
