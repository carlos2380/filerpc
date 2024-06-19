package fileutils

import (
	"fmt"
	"os"
	"path/filepath"

	log "filerpc/internal/logger"
)

func ReadFile(fileType, version string) (string, []byte, error) {
	filePath := filepath.Join(fileType, version+".json")
	log.Logger.Debug("Reading file from path:", filePath)
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Logger.Error("File not found:", filePath)
			return "", nil, fmt.Errorf("file not found: %s", filePath)
		}
		log.Logger.Error("Error reading file:", err)
		return "", nil, fmt.Errorf("error reading file: %w", err)
	}
	log.Logger.Debug("File read successfully")
	return filePath, data, nil
}
