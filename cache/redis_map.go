package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisMap[T any] struct {
	client      *redis.Client
	serviceName string
	keyPrefix   string
	ttl         time.Duration
}

func NewRedisMap[T any](client *redis.Client, serviceName, keyPrefix string, ttl time.Duration) MapCache[T] {
	return &RedisMap[T]{
		client:      client,
		serviceName: serviceName,
		keyPrefix:   keyPrefix,
		ttl:         ttl,
	}
}

func (c *RedisMap[T]) Get(ctx context.Context, key string) (T, error) {
	var data T

	result := c.client.HGet(ctx, c.key(), key)

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

func (c *RedisMap[T]) GetAll(ctx context.Context) ([]T, error) {
	result, err := c.client.HVals(ctx, c.key()).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	data := make([]T, len(result))
	for i := range result {
		if err = json.Unmarshal([]byte(result[i]), &data[i]); err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (c *RedisMap[T]) GetAllMap(ctx context.Context) (map[string]T, error) {
	result, err := c.client.HGetAll(ctx, c.key()).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	data := make(map[string]T, len(result))
	for key := range result {
		var value T
		if err = json.Unmarshal([]byte(result[key]), &value); err != nil {
			return nil, err
		}
		data[key] = value
	}

	return data, nil
}

// Set not recommended
func (c *RedisMap[T]) Set(ctx context.Context, key string, value T) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.HSet(ctx, c.key(), key, data).Err()
}

func (c *RedisMap[T]) SetAll(ctx context.Context, values []T, keyFunc func(v T) string) error {
	cacheValues := make(map[string]string, len(values))
	for i := range values {
		data, err := json.Marshal(values[i])
		if err != nil {
			return err
		}
		cacheValues[keyFunc(values[i])] = string(data)
	}

	if err := c.client.HSet(ctx, c.key(), cacheValues).Err(); err != nil {
		return err
	}

	return c.client.Expire(ctx, c.key(), c.ttl).Err()
}

func (c *RedisMap[T]) SetAllMap(ctx context.Context, valuesMap map[string]T) error {
	cacheValues := make(map[string]string, len(valuesMap))
	for key := range valuesMap {
		data, err := json.Marshal(valuesMap[key])
		if err != nil {
			return err
		}
		cacheValues[key] = string(data)
	}

	if err := c.client.HSet(ctx, c.key(), cacheValues).Err(); err != nil {
		return err
	}

	return c.client.Expire(ctx, c.key(), c.ttl).Err()
}

func (c *RedisMap[T]) Exists(ctx context.Context, key string) (bool, error) {
	return c.client.HExists(ctx, c.key(), key).Result()
}

func (c *RedisMap[T]) Delete(ctx context.Context, keys ...string) error {
	return c.client.HDel(ctx, c.key(), keys...).Err()
}

func (c *RedisMap[T]) DeleteAll(ctx context.Context) error {
	return c.client.Del(ctx, c.key()).Err()
}

func (c *RedisMap[T]) key() string {
	return fmt.Sprintf("%s::%s", c.serviceName, c.keyPrefix)
}
