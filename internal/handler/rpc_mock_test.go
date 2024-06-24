package handler_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"filerpc/internal/datastore"
	"filerpc/internal/errors"
	fileutils "filerpc/internal/file"
	"filerpc/internal/handler"
	pb "filerpc/internal/proto"
	service "filerpc/internal/service"
)

func TestReadFileRPCWithMocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDataStore := datastore.NewMockFileDataStore(ctrl)
	mockFileReader := fileutils.NewMockFileReader(ctrl)
	server := handler.NewServer(mockDataStore, mockFileReader)

	fileContent := `{"reward": "100000"}`
	fileType := "core"
	version := "1.0.0"
	filePath := "core/1.0.0.json"
	expectedHash := service.CalculateHash([]byte(fileContent))

	tableTests := []struct {
		desc             string
		req              *pb.FileRequest
		mockSetup        func()
		expectedResponse *pb.FileResponse
		expectedError    error
	}{
		{
			desc: "Valid file and matching hash",
			req: &pb.FileRequest{
				Type:    fileType,
				Version: version,
				Hash:    expectedHash,
			},
			mockSetup: func() {
				mockFileReader.EXPECT().
					ReadFile(fileType, version).
					Return(filePath, []byte(fileContent), nil)

				mockDataStore.EXPECT().
					Save(gomock.Any(), filePath, []byte(fileContent), expectedHash).
					Return(nil)
			},
			expectedResponse: &pb.FileResponse{
				Type:    fileType,
				Version: version,
				Hash:    expectedHash,
				Content: fileContent,
			},
			expectedError: nil,
		},
		{
			desc: "Valid file but non-matching hash",
			req: &pb.FileRequest{
				Type:    fileType,
				Version: version,
				Hash:    "invalidhash",
			},
			mockSetup: func() {
				mockFileReader.EXPECT().
					ReadFile(fileType, version).
					Return(filePath, []byte(fileContent), nil)
			},
			expectedResponse: &pb.FileResponse{
				Type:    fileType,
				Version: version,
				Hash:    expectedHash,
				Content: "",
			},
			expectedError: nil,
		},
		{
			desc: "File not found",
			req: &pb.FileRequest{
				Type:    "nonexistent",
				Version: version,
				Hash:    "",
			},
			mockSetup: func() {
				mockFileReader.EXPECT().
					ReadFile("nonexistent", version).
					Return("", nil, errors.ErrFileNotFound)
			},
			expectedResponse: nil,
			expectedError:    errors.ErrFileNotFound,
		},
	}

	for _, tt := range tableTests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.mockSetup()

			response, err := server.ReadFile(context.Background(), tt.req)

			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponse.Type, response.Type)
				assert.Equal(t, tt.expectedResponse.Version, response.Version)
				assert.Equal(t, tt.expectedResponse.Hash, response.Hash)
				assert.Equal(t, tt.expectedResponse.Content, response.Content)
			}
		})
	}
}
