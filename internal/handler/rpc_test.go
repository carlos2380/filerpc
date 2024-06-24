package handler_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"filerpc/internal/datastore"
	"filerpc/internal/errors"
	fileutils "filerpc/internal/file"
	"filerpc/internal/handler"
	pb "filerpc/internal/proto"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

var redisClient *redis.Client

func setup() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func teardown() {
	redisClient.FlushAll(context.Background())
}

func createTestFile(path, content string) error {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func calculateHash(content string) string {
	hash := sha256.Sum256([]byte(content))
	return hex.EncodeToString(hash[:])
}

func TestReadFileRPCWithRedis(t *testing.T) {
	setup()
	defer teardown()

	fileContent := `{
		"reward": "100000"
	}`
	testFilePath := "core/1.0.0.json"
	err := createTestFile(testFilePath, fileContent)
	if err != nil {
		t.Fatalf("Error creating test file: %v", err)
	}
	defer os.RemoveAll(filepath.Dir(testFilePath))

	expectedHash := calculateHash(fileContent)

	dstore := datastore.NewRedisFileDataStore(redisClient)
	fileReader := fileutils.DefaultFileReader{}
	server := handler.NewServer(dstore, fileReader)

	tableTests := []struct {
		desc             string
		req              *pb.FileRequest
		expectedResponse *pb.FileResponse
		expectedError    error
	}{
		{
			desc: "Valid file and matching hash",
			req: &pb.FileRequest{
				Type:    "core",
				Version: "1.0.0",
				Hash:    expectedHash,
			},
			expectedResponse: &pb.FileResponse{
				Type:    "core",
				Version: "1.0.0",
				Hash:    expectedHash,
				Content: fileContent,
			},
			expectedError: nil,
		},
		{
			desc: "Valid file but non-matching hash",
			req: &pb.FileRequest{
				Type:    "core",
				Version: "1.0.0",
				Hash:    "invalidhash",
			},
			expectedResponse: &pb.FileResponse{
				Type:    "core",
				Version: "1.0.0",
				Hash:    expectedHash,
				Content: "",
			},
			expectedError: nil,
		},
		{
			desc: "File not found",
			req: &pb.FileRequest{
				Type:    "nonexistent",
				Version: "1.0.0",
				Hash:    "",
			},
			expectedResponse: nil,
			expectedError:    errors.ErrFileNotFound,
		},
	}

	for _, tt := range tableTests {
		t.Run(tt.desc, func(t *testing.T) {
			// Limpiar la base de datos antes de cada test
			teardown()
			setup()

			// Ejecutar el test
			response, err := server.ReadFile(context.Background(), tt.req)

			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponse.Type, response.Type)
				assert.Equal(t, tt.expectedResponse.Version, response.Version)
				assert.Equal(t, tt.expectedResponse.Hash, response.Hash)
				assert.Equal(t, tt.expectedResponse.Content, response.Content)

				// Comprobar el contenido en la base de datos
				result, err := redisClient.HGetAll(context.Background(), "core/1.0.0.json").Result()
				assert.NoError(t, err)
				if response.Content != "" {
					assert.Equal(t, response.Content, result["content"])
					assert.Equal(t, response.Hash, result["hash"])
				}
			}
		})
	}
}
