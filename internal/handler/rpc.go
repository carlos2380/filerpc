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
	DataStore  datastore.FileDataStore
	FileReader fileutils.FileReader
}

func NewServer(dstore datastore.FileDataStore, freader fileutils.FileReader) *Server {
	return &Server{DataStore: dstore, FileReader: freader}
}

func (s *Server) ReadFile(ctx context.Context, req *pb.FileRequest) (*pb.FileResponse, error) {
	fileType, version, hash := getParamas(req)

	filePath, content, err := s.FileReader.ReadFile(fileType, version)
	if err != nil {
		log.Logger.Error("Error reading file: ", err)
		return nil, errors.ErrFileNotFound
	}

	response := &pb.FileResponse{
		Type:    fileType,
		Version: version,
	}

	contentHash := service.CalculateHash(content)
	response.Hash = contentHash

	if hash == contentHash {
		if err := s.DataStore.Save(ctx, filePath, content, contentHash); err != nil {
			log.Logger.Error("Error saving to datastore: ", err)
		}
		response.Content = string(content)
	}

	return response, nil
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
