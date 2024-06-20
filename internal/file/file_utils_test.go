package fileutils_test

import (
	"filerpc/internal/errors"
	fileutils "filerpc/internal/file"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestFile(t *testing.T, path, content string) {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	assert.NoError(t, err)

	err = os.WriteFile(path, []byte(content), 0644)
	assert.NoError(t, err)
}

func TestDefaultFileReader_ReadFile(t *testing.T) {
	reader := fileutils.DefaultFileReader{}

	t.Run("File exists and can be read", func(t *testing.T) {
		testFilePath := "testdir/testfile.json"
		testContent := `{"reward": "100000"}`
		createTestFile(t, testFilePath, testContent)
		defer os.RemoveAll("testdir")

		filePath, data, err := reader.ReadFile("testdir", "testfile")
		assert.NoError(t, err)
		assert.Equal(t, testFilePath, filePath)
		assert.Equal(t, []byte(testContent), data)
	})

	t.Run("File does not exist", func(t *testing.T) {
		filePath, data, err := reader.ReadFile("nonexistentdir", "nonexistentfile")
		assert.ErrorIs(t, err, errors.ErrFileNotFound)
		assert.Empty(t, filePath)
		assert.Nil(t, data)
	})

	t.Run("Error reading file", func(t *testing.T) {
		testFilePath := "testdir/unreadablefile.json"
		createTestFile(t, testFilePath, "")
		defer os.RemoveAll("testdir")

		err := os.Chmod(testFilePath, 0000)
		assert.NoError(t, err)
		defer os.Chmod(testFilePath, 0644)

		filePath, data, err := reader.ReadFile("testdir", "unreadablefile")
		assert.ErrorIs(t, err, errors.ErrReadFile)
		assert.Empty(t, filePath)
		assert.Nil(t, data)
	})
}
