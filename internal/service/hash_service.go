package service

import (
	"crypto/sha256"
	"encoding/hex"
)

func CalculateHash(content []byte) string {
	hash := sha256.Sum256(content)
	return hex.EncodeToString(hash[:])
}
