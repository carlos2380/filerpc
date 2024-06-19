package handler

import (
	"context"

	"filerpc/internal/errors"
	fileutils "filerpc/internal/file"
	log "filerpc/internal/logger"
	pb "filerpc/internal/proto"
	service "filerpc/internal/service"

	"github.com/go-redis/redis/v8"
)

type Server struct {
    pb.UnimplementedFileServiceServer
    RedisClient *redis.Client
}

// NewServer creates a new gRPC server with the provided Redis client
func NewServer(redisClient *redis.Client) *Server {
    return &Server{RedisClient: redisClient}
}


func (s *Server) ReadFile(ctx context.Context, req *pb.FileRequest) (*pb.FileResponse, error) {
	fileType, version, hash := getParamas(req)

	_, content, err := fileutils.ReadFile(fileType, version)
	if err != nil {
		log.Logger.Error("Error reading file: ", err)
		return nil, errors.ErrFileNotFound
	}

	contentHash := service.CalculateHash(content)

	if hash != contentHash {
		hash = ""
	}

	return &pb.FileResponse{
		Type:    fileType,
		Version: version,
		Hash:    hash,
		Content: content,
	}, nil
}

func getParamas(req *pb.FileRequest) (string, string, string) {

	defaultType := "core"
	defaultVersion := "1.0.0"

	fileType := req.GetType()
	if fileType == "" {
		fileType = defaultType
	}

	version := req.GetVersion()
	if version == "" {
		version = defaultVersion
	}

	hash := req.GetHash()

	return fileType, version, hash
}
