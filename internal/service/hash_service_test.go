package service_test

import (
	"filerpc/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateHash(t *testing.T) {
	content := []byte("Hello World")
	hash := service.CalculateHash(content)

	expectedHash := "a591a6d40bf420404a011733cfb7b190d62c65bf0bcda32b57b277d9ad9f146e"
	assert.Equal(t, expectedHash, hash, "Hashes should match")
}
