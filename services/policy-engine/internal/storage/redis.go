package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/trustchain/policy-engine/internal/domain"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr string, password string) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	return &RedisCache{client: rdb}
}

func (c *RedisCache) GetPolicyResponse(ctx context.Context, digest string) (*domain.PolicyResponse, error) {
	val, err := c.client.Get(ctx, "policy:"+digest).Result()
	if err == redis.Nil {
		return nil, nil // cache miss
	} else if err != nil {
		return nil, fmt.Errorf("redis get error: %w", err)
	}

	var resp domain.PolicyResponse
	if err := json.Unmarshal([]byte(val), &resp); err != nil {
		return nil, fmt.Errorf("unmarshaling cached response: %w", err)
	}
	resp.Cached = true // flag that it came from cache
	return &resp, nil
}

func (c *RedisCache) SetPolicyResponse(ctx context.Context, digest string, resp domain.PolicyResponse) error {
	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("marshaling response: %w", err)
	}
	
	// Cache for 5 minutes
	err = c.client.Set(ctx, "policy:"+digest, data, 5*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("redis set error: %w", err)
	}
	return nil
}
