package datastore

import (
	"context"
	"filerpc/internal/errors"

	"github.com/go-redis/redis/v8"
)

type RedisFileDataStore struct {
	Client *redis.Client
}

func NewRedisFileDataStore(client *redis.Client) *RedisFileDataStore {
	return &RedisFileDataStore{Client: client}
}

func (r *RedisFileDataStore) Save(ctx context.Context, key string, content []byte, hash string) error {
	err := r.Client.HMSet(ctx, key, map[string]interface{}{
		"content": content,
		"hash":    hash,
	}).Err()
	return err
}

func (r *RedisFileDataStore) Get(ctx context.Context, key string) (map[string]string, error) {
	return r.Client.HGetAll(ctx, key).Result()
}

func InitializeRedisClient(ctx context.Context, redisAddr string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, errors.ErrFailedToConnectRedis
	}

	return client, nil
}
