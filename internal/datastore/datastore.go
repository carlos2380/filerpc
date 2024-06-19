package datastore

import (
	"context"
)

type FileDataStore interface {
	Save(ctx context.Context, key string, content []byte, hash string) error
	Get(ctx context.Context, key string) (map[string]string, error)
}
