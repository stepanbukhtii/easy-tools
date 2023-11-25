package cache

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMemory(t *testing.T) {
	key := "key"
	key2 := "key2"
	keyValue := "value"
	keyValue2 := "value2"
	ctx := context.Background()

	t.Run("Get, Set", func(t *testing.T) {
		cache := NewMemory[string](time.Minute, time.Second)

		value, err := cache.Get(ctx, key)
		require.ErrorIs(t, err, ErrNotFound)
		require.Empty(t, value)

		require.NoError(t, cache.Set(ctx, key, keyValue))

		value, err = cache.Get(ctx, key)
		require.NoError(t, err)
		require.Equal(t, keyValue, value)
	})

	t.Run("GetAll, SetAll", func(t *testing.T) {
		cache := NewMemory[string](time.Minute, time.Second)

		values, err := cache.GetAll(ctx)
		require.ErrorIs(t, err, ErrNotFound)
		require.Empty(t, values)

		data := []string{keyValue, keyValue2}
		keyFunc := func(value string) string {
			if value == keyValue {
				return key
			}
			return key2
		}
		require.NoError(t, cache.SetAll(ctx, data, keyFunc))

		values, err = cache.GetAll(ctx)
		require.NoError(t, err)
		slices.Sort(values)
		require.Equal(t, []string{keyValue, keyValue2}, values)
	})

	t.Run("GetAllMap, SetAllMap", func(t *testing.T) {
		cache := NewMemory[string](time.Minute, time.Second)

		valuesMap, err := cache.GetAllMap(ctx)
		require.ErrorIs(t, err, ErrNotFound)
		require.Empty(t, valuesMap)

		data := map[string]string{key: keyValue, key2: keyValue2}
		require.NoError(t, cache.SetAllMap(ctx, data))

		require.NoError(t, cache.Set(ctx, key, keyValue))
		require.NoError(t, cache.Set(ctx, key2, keyValue2))

		valuesMap, err = cache.GetAllMap(ctx)
		require.NoError(t, err)
		expectedValues := map[string]string{
			key:  keyValue,
			key2: keyValue2,
		}
		require.Equal(t, expectedValues, valuesMap)
	})

	t.Run("SetNX", func(t *testing.T) {
		cache := NewMemory[string](time.Minute, time.Second)

		ok, err := cache.SetNX(context.Background(), key, "value")
		require.NoError(t, err)
		require.True(t, ok)

		ok, err = cache.SetNX(context.Background(), key, "value")
		require.NoError(t, err)
		require.False(t, ok)
	})

	t.Run("Exists", func(t *testing.T) {
		cache := NewMemory[string](time.Minute, time.Second)

		value, err := cache.Exists(ctx, key)
		require.NoError(t, err)
		require.False(t, value)

		require.NoError(t, cache.Set(ctx, key, keyValue))

		value, err = cache.Exists(ctx, key)
		require.NoError(t, err)
		require.True(t, value)
	})

	t.Run("Delete, DeleteAll", func(t *testing.T) {
		cache := NewMemory[string](time.Minute, time.Second)

		require.NoError(t, cache.Set(ctx, key, keyValue))
		require.NoError(t, cache.Delete(ctx, key))

		value, err := cache.Exists(ctx, key)
		require.NoError(t, err)
		require.False(t, value)

		require.NoError(t, cache.Set(ctx, key, keyValue))
		require.NoError(t, cache.Set(ctx, key2, keyValue2))
		require.NoError(t, cache.DeleteAll(ctx))

		value, err = cache.Exists(ctx, key)
		require.NoError(t, err)
		require.False(t, value)

		value, err = cache.Exists(ctx, key2)
		require.NoError(t, err)
		require.False(t, value)
	})
}
