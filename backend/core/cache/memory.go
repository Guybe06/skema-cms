package cache

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"
)

type memoryEntry struct {
	value     []byte
	expiresAt time.Time
}

type memoryCache struct {
	mu    sync.RWMutex
	store map[string]memoryEntry
}

func newMemoryCache() *memoryCache {
	return &memoryCache{store: make(map[string]memoryEntry)}
}

func (m *memoryCache) Set(_ context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.store[key] = memoryEntry{value: data, expiresAt: time.Now().Add(ttl)}
	return nil
}

func (m *memoryCache) Get(_ context.Context, key string, dest any) error {
	m.mu.RLock()
	entry, ok := m.store[key]
	m.mu.RUnlock()

	if !ok || time.Now().After(entry.expiresAt) {
		m.mu.Lock()
		delete(m.store, key)
		m.mu.Unlock()
		return ErrCacheMiss
	}

	return json.Unmarshal(entry.value, dest)
}

func (m *memoryCache) Del(_ context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.store, key)
	return nil
}

func (m *memoryCache) DelPattern(_ context.Context, pattern string) error {
	prefix := strings.TrimSuffix(pattern, "*")
	m.mu.Lock()
	defer m.mu.Unlock()

	for key := range m.store {
		if strings.HasPrefix(key, prefix) {
			delete(m.store, key)
		}
	}
	return nil
}
