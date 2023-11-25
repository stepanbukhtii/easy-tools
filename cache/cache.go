package cache

import (
	"context"
	"errors"
	"time"
)

var DefaultTTL = time.Hour

var (
	ErrNotFound = errors.New("not found")
)

type Cache[T any] interface {
	Get(ctx context.Context, key string) (T, error)
	Set(ctx context.Context, key string, value T) error
	SetNX(ctx context.Context, key string, value T) (bool, error)
	Exists(ctx context.Context, key string) (bool, error)
	Delete(ctx context.Context, keys ...string) error
}

type MapCache[T any] interface {
	Get(ctx context.Context, key string) (T, error)
	GetAll(ctx context.Context) ([]T, error)
	GetAllMap(ctx context.Context) (map[string]T, error)
	Set(ctx context.Context, key string, value T) error
	SetAll(ctx context.Context, values []T, keyFunc func(v T) string) error
	SetAllMap(ctx context.Context, valuesMap map[string]T) error
	Exists(ctx context.Context, key string) (bool, error)
	Delete(ctx context.Context, keys ...string) error
	DeleteAll(ctx context.Context) error
}
