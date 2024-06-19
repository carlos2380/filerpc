package fileutils

import (
	"os"
	"path/filepath"

	"filerpc/internal/errors"
	log "filerpc/internal/logger"
)

func ReadFile(fileType, version string) (string, []byte, error) {
	filePath := filepath.Join(fileType, version+".json")
	log.Logger.Debug("Reading file from path:", filePath)
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Logger.Error("File not found:", filePath)
			return "", nil, errors.ErrFileNotFound
		}
		log.Logger.Error("Error reading file:", err)
		return "", nil, errors.ErrReadFile
	}
	log.Logger.Debug("File read successfully")
	return filePath, data, nil
}
