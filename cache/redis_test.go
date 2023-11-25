package cache

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestRedis(t *testing.T) {
	key := "field_key"
	keyValue := "value"
	keySetNX := "key_set_nx"

	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	cache := NewRedis[string](redisClient, "serviceName", "key", time.Minute)

	ctx := context.Background()

	value, err := cache.Get(ctx, key)
	require.ErrorIs(t, err, ErrNotFound)
	require.Empty(t, value)

	exists, err := cache.Exists(ctx, key)
	require.NoError(t, err)
	require.False(t, exists)

	require.NoError(t, cache.Set(ctx, key, keyValue))

	value, err = cache.Get(ctx, key)
	require.NoError(t, err)
	require.Equal(t, keyValue, value)

	exists, err = cache.Exists(ctx, key)
	require.NoError(t, err)
	require.True(t, exists)

	require.NoError(t, cache.Delete(ctx, key))

	exists, err = cache.Exists(ctx, key)
	require.NoError(t, err)
	require.False(t, exists)

	ok, err := cache.SetNX(context.Background(), keySetNX, "value")
	require.NoError(t, err)
	require.True(t, ok)

	ok, err = cache.SetNX(context.Background(), keySetNX, "value")
	require.NoError(t, err)
	require.False(t, ok)

	require.NoError(t, cache.Delete(ctx, keySetNX))
}

func TestRedisMap(t *testing.T) {
	key := "redis_map"
	keyValue := "value"
	valuesSlice := []string{"slice_value1", "slice_value2"}
	valuesMap := map[string]string{"key_map1": "value1", "key_map2": "value2"}

	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	cache := NewRedisMap[string](redisClient, "serviceName", "map_key", time.Minute)

	ctx := context.Background()

	require.NoError(t, cache.DeleteAll(ctx))

	value, err := cache.Get(ctx, key)
	require.ErrorIs(t, err, ErrNotFound)
	require.Empty(t, value)

	exists, err := cache.Exists(ctx, key)
	require.NoError(t, err)
	require.False(t, exists)

	require.NoError(t, cache.Set(ctx, key, keyValue))

	value, err = cache.Get(ctx, key)
	require.NoError(t, err)
	require.Equal(t, keyValue, value)

	exists, err = cache.Exists(ctx, key)
	require.NoError(t, err)
	require.True(t, exists)

	require.NoError(t, cache.Delete(ctx, key))

	exists, err = cache.Exists(ctx, key)
	require.NoError(t, err)
	require.False(t, exists)

	require.NoError(t, cache.SetAll(ctx, valuesSlice, func(v string) string { return v }))

	values, err := cache.GetAll(ctx)
	require.NoError(t, err)
	slices.Sort(values)
	require.Equal(t, valuesSlice, values)

	require.NoError(t, cache.DeleteAll(ctx))

	require.NoError(t, cache.SetAllMap(ctx, valuesMap))

	resultMap, err := cache.GetAllMap(ctx)
	require.NoError(t, err)
	require.Equal(t, valuesMap, resultMap)

	require.NoError(t, cache.DeleteAll(ctx))
}
