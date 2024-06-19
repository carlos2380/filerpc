package handler

import (
	"context"

	"filerpc/internal/datastore"
	"filerpc/internal/errors"
	fileutils "filerpc/internal/file"
	log "filerpc/internal/logger"
	pb "filerpc/internal/proto"
	service "filerpc/internal/service"
)

type Server struct {
	pb.UnimplementedFileServiceServer
	DataStore datastore.FileDataStore
}

// NewServer creates a new gRPC server with the provided DataStore
func NewServer(ds datastore.FileDataStore) *Server {
	return &Server{DataStore: ds}
}

func (s *Server) ReadFile(ctx context.Context, req *pb.FileRequest) (*pb.FileResponse, error) {
	fileType, version, hash := getParamas(req)

	filePath, content, err := fileutils.ReadFile(fileType, version)
	if err != nil {
		log.Logger.Error("Error reading file: ", err)
		return nil, errors.ErrFileNotFound
	}

	contentHash := service.CalculateHash(content)

	if hash == contentHash {
		if err := s.DataStore.Save(ctx, filePath, content, contentHash); err != nil {
			log.Logger.Error("Error saving to datastore: ", err)
		}
	} else {
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
