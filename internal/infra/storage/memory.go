package storage

import (
	"context"
	"path/filepath"
	"sync"
	"time"

	"github.com/g-villarinho/godis/internal/core/domain/entity"
	"github.com/g-villarinho/godis/internal/core/ports"
)

type MemoryStorage struct {
	data        map[string]*entity.Item
	mu          sync.RWMutex
	stopCleanup chan struct{}
}

func NewMemoryStorage() ports.KeyValueRepository {
	return &MemoryStorage{
		data:        make(map[string]*entity.Item),
		stopCleanup: make(chan struct{}),
	}
}

func (m *MemoryStorage) Set(ctx context.Context, key string, value string) {
	if ctx.Err() != nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = &entity.Item{
		Value:     value,
		ExpiresAt: nil,
	}
}

func (m *MemoryStorage) Get(ctx context.Context, key string) (string, bool) {
	if ctx.Err() != nil {
		return "", false
	}
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, exists := m.data[key]
	if !exists {
		return "", false
	}

	if item.IsExpired(time.Now().Unix()) {
		return "", false
	}

	return item.Value, true
}

func (m *MemoryStorage) Del(ctx context.Context, key string) int {
	if ctx.Err() != nil {
		return -1
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.data[key]; !exists {
		return -1
	}

	delete(m.data, key)
	return 0
}

func (m *MemoryStorage) Expire(ctx context.Context, key string, seconds int) bool {
	if ctx.Err() != nil {
		return false
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	item, exists := m.data[key]
	if !exists {
		return false
	}

	expiresAt := time.Now().Unix() + int64(seconds)
	item.ExpiresAt = &expiresAt
	return true
}

func (m *MemoryStorage) TTL(ctx context.Context, key string) int64 {
	if ctx.Err() != nil {
		return -1
	}
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, exists := m.data[key]
	if !exists {
		return -1
	}

	if item.ExpiresAt == nil {
		return -1
	}

	now := time.Now().Unix()
	remaining := *item.ExpiresAt - now

	if remaining <= 0 {
		return 0
	}

	return remaining
}

func (m *MemoryStorage) Persist(ctx context.Context, key string) bool {
	if ctx.Err() != nil {
		return false
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	item, exists := m.data[key]
	if !exists {
		return false
	}

	item.ExpiresAt = nil
	return true
}

func (m *MemoryStorage) Keys(ctx context.Context, pattern string) []string {
	if ctx.Err() != nil {
		return []string{}
	}
	m.mu.RLock()
	defer m.mu.RUnlock()

	now := time.Now().Unix()
	var matches []string
	for key, item := range m.data {
		if item.IsExpired(now) {
			continue
		}

		if matchPattern(key, pattern) {
			matches = append(matches, key)
		}
	}

	return matches
}

func (m *MemoryStorage) Exists(ctx context.Context, key string) bool {
	if ctx.Err() != nil {
		return false
	}
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, exists := m.data[key]
	if !exists {
		return false
	}

	return !item.IsExpired(time.Now().Unix())
}

func (m *MemoryStorage) Size(ctx context.Context) int {
	if ctx.Err() != nil {
		return 0
	}
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.data)
}

func (m *MemoryStorage) StartCleanUp(intervalMs int64) {
	interval := time.Duration(intervalMs) * time.Millisecond

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				m.cleanupExpired()
			case <-m.stopCleanup:
				return
			}
		}
	}()
}

func (m *MemoryStorage) StopCleanUp() {
	close(m.stopCleanup)
}

func (m *MemoryStorage) cleanupExpired() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now().Unix()

	for key, item := range m.data {
		if item.IsExpired(now) {
			delete(m.data, key)
		}
	}
}

func matchPattern(key, pattern string) bool {
	if pattern == "*" {
		return true
	}

	matched, err := filepath.Match(pattern, key)
	if err != nil {
		return key == pattern
	}

	return matched
}
