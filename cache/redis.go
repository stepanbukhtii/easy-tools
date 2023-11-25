package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
)

type Redis[T any] struct {
	client      *redis.Client
	serviceName string
	keyPrefix   string
	ttl         time.Duration
}

func NewRedis[T any](client *redis.Client, serviceName, keyPrefix string, ttl time.Duration) Cache[T] {
	return &Redis[T]{
		client:      client,
		serviceName: serviceName,
		keyPrefix:   keyPrefix,
		ttl:         ttl,
	}
}

func (c *Redis[T]) Get(ctx context.Context, key string) (T, error) {
	var data T

	result := c.client.Get(ctx, c.key(key))

	if result.Err() != nil {
		if errors.Is(result.Err(), redis.Nil) {
			return data, ErrNotFound
		}
		return data, result.Err()
	}

	dataBytes, err := result.Bytes()
	if err != nil {
		return data, err
	}

	if err = json.Unmarshal(dataBytes, &data); err != nil {
		return data, err
	}

	return data, nil
}

func (c *Redis[T]) Set(ctx context.Context, key string, value T) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, c.key(key), data, c.ttl).Err()
}

func (c *Redis[T]) SetNX(ctx context.Context, key string, value T) (bool, error) {
	data, err := jsoniter.Marshal(value)
	if err != nil {
		return false, err
	}

	args := redis.SetArgs{
		Mode: "NX",
		TTL:  c.ttl,
	}

	if err = c.client.SetArgs(ctx, c.key(key), data, args).Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (c *Redis[T]) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.client.Exists(ctx, c.key(key)).Result()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (c *Redis[T]) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	redisKeys := lo.Map(keys, func(k string, _ int) string { return c.key(k) })

	return c.client.Del(ctx, redisKeys...).Err()
}

func (c *Redis[T]) key(key string) string {
	return fmt.Sprintf("%s::%s::%s", c.serviceName, c.keyPrefix, key)
}
