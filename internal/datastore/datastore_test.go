package datastore_test

import (
	"context"
	"testing"

	"filerpc/internal/datastore"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func setupRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func teardownRedis(client *redis.Client) {
	client.FlushAll(context.Background())
	client.Close()
}

func TestRedisFileDataStore(t *testing.T) {
	client := setupRedis()
	defer teardownRedis(client)

	dstore := datastore.NewRedisFileDataStore(client)

	ctx := context.Background()
	key := "testkey"
	content := []byte("testcontent")
	hash := "testhash"

	err := dstore.Save(ctx, key, content, hash)
	assert.NoError(t, err)

	result, err := dstore.Get(ctx, key)
	assert.NoError(t, err)

	savedContent := result["content"]
	savedHash := result["hash"]

	assert.Equal(t, string(content), savedContent)
	assert.Equal(t, hash, savedHash)
}
