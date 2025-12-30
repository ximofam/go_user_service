package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisCacheService struct {
	rdb *redis.Client
}

func NewCacheService(rdb *redis.Client) *redisCacheService {
	return &redisCacheService{rdb: rdb}
}

func (s *redisCacheService) Get(ctx context.Context, key string, dest any) error {
	data, err := s.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}

func (s *redisCacheService) Set(ctx context.Context, key string, value any, tll time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return s.rdb.Set(ctx, key, data, tll).Err()
}

func (s *redisCacheService) Del(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	return s.rdb.Del(ctx, keys...).Err()
}
