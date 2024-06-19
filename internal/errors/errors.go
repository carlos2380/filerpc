package errors

import (
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

var (
    ErrFileNotFound = status.New(codes.NotFound, "file not found").Err()
)
