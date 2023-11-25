package cache

import (
	"context"
	"sync"
	"time"
)

type entry[T any] struct {
	value     T
	expiresAt time.Time
}

func (e entry[T]) isExpired() bool {
	return time.Now().After(e.expiresAt)
}

type Memory[T any] struct {
	mu     sync.RWMutex
	data   map[string]entry[T]
	ttl    time.Duration
	stopGC chan struct{}
}

func NewMemory[T any](ttl, cleanInterval time.Duration) *Memory[T] {
	c := &Memory[T]{
		data:   make(map[string]entry[T]),
		stopGC: make(chan struct{}),
		ttl:    ttl,
	}
	go c.runGC(cleanInterval)

	return c
}

func (c *Memory[T]) Get(_ context.Context, key string) (T, error) {
	c.mu.RLock()
	e, ok := c.data[key]
	c.mu.RUnlock()

	if !ok || e.isExpired() {
		var zero T
		if !ok {
			return zero, ErrNotFound
		}

		c.mu.Lock()
		if e, exists := c.data[key]; exists && e.isExpired() {
			delete(c.data, key)
		}
		c.mu.Unlock()

		return zero, ErrNotFound
	}

	return e.value, nil
}

func (c *Memory[T]) GetAll(_ context.Context) ([]T, error) {
	c.mu.RLock()
	snapshot := make(map[string]entry[T], len(c.data))
	for k, e := range c.data {
		snapshot[k] = e
	}
	c.mu.RUnlock()

	if len(snapshot) == 0 {
		return nil, ErrNotFound
	}

	result := make([]T, 0, len(snapshot))
	var expired []string

	for k, e := range snapshot {
		if e.isExpired() {
			expired = append(expired, k)
		} else {
			result = append(result, e.value)
		}
	}

	if len(expired) > 0 {
		c.mu.Lock()
		for _, k := range expired {
			if e, exists := c.data[k]; exists && e.isExpired() {
				delete(c.data, k)
			}
		}
		c.mu.Unlock()
	}

	return result, nil
}

func (c *Memory[T]) GetAllMap(_ context.Context) (map[string]T, error) {
	c.mu.RLock()
	snapshot := make(map[string]entry[T], len(c.data))
	for k, e := range c.data {
		snapshot[k] = e
	}
	c.mu.RUnlock()

	if len(snapshot) == 0 {
		return nil, ErrNotFound
	}

	result := make(map[string]T)
	var expired []string

	for k, e := range snapshot {
		if e.isExpired() {
			expired = append(expired, k)
		} else {
			result[k] = e.value
		}
	}

	if len(expired) > 0 {
		c.mu.Lock()
		for _, k := range expired {
			if e, exists := c.data[k]; exists && e.isExpired() {
				delete(c.data, k)
			}
		}
		c.mu.Unlock()
	}

	return result, nil
}

func (c *Memory[T]) Set(_ context.Context, key string, value T) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = entry[T]{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	}

	return nil
}

func (c *Memory[T]) SetAll(_ context.Context, values []T, keyFunc func(v T) string) error {
	expiresAt := time.Now().Add(c.ttl)

	c.mu.Lock()
	defer c.mu.Unlock()

	for i := range values {
		c.data[keyFunc(values[i])] = entry[T]{
			value:     values[i],
			expiresAt: expiresAt,
		}
	}

	return nil
}

func (c *Memory[T]) SetAllMap(_ context.Context, valuesMap map[string]T) error {
	expiresAt := time.Now().Add(c.ttl)

	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range valuesMap {
		c.data[k] = entry[T]{
			value:     v,
			expiresAt: expiresAt,
		}
	}

	return nil
}

func (c *Memory[T]) SetNX(ctx context.Context, key string, value T) (bool, error) {
	if _, err := c.Get(ctx, key); err == nil {
		return false, nil
	}

	if err := c.Set(ctx, key, value); err != nil {
		return false, err
	}

	return true, nil
}

func (c *Memory[T]) Exists(_ context.Context, key string) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.data[key]
	return ok, nil
}

func (c *Memory[T]) Delete(_ context.Context, keys ...string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, key := range keys {
		delete(c.data, key)
	}
	return nil
}

func (c *Memory[T]) DeleteAll(_ context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	clear(c.data)
	return nil
}

func (c *Memory[T]) Close() {
	close(c.stopGC)
}

func (c *Memory[T]) runGC(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.evict()
		case <-c.stopGC:
			return
		}
	}
}

func (c *Memory[T]) evict() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, e := range c.data {
		if e.isExpired() {
			delete(c.data, k)
		}
	}
}
