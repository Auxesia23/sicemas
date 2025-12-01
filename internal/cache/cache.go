package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Set(ctx context.Context, key string, data any, ttl time.Duration) error
	Get(ctx context.Context, key string, dest any) error
	Delete(ctx context.Context, key string) error
}

type cacheImpl struct {
	RDB *redis.Client
}

func NewCache(rdb *redis.Client) Cache {
	return &cacheImpl{
		RDB: rdb,
	}
}

func (r *cacheImpl) Set(ctx context.Context, key string, data any, ttl time.Duration) error {
	bytesData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if err := r.RDB.Set(ctx, key, bytesData, ttl).Err(); err != nil {
		return err
	}
	return nil
}

func (r *cacheImpl) Get(ctx context.Context, key string, dest any) error {
	val, err := r.RDB.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), &dest)
}

func (r *cacheImpl) Delete(ctx context.Context, key string) error {
	if err := r.RDB.Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}
