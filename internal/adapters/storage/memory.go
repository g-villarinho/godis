package storage

import (
	"context"
	"sync"
	"time"

	"github.com/g-villarinho/godis/internal/core/domain"
	"github.com/g-villarinho/godis/internal/core/ports"
)

type MemoryStorage struct {
	data        map[string]*domain.Item
	mu          sync.RWMutex
	stopCleanup chan struct{}
}

func NewMemoryStorage() ports.KeyValueRepository {
	return &MemoryStorage{
		data:        make(map[string]*domain.Item),
		stopCleanup: make(chan struct{}),
	}
}

func (m *MemoryStorage) Set(ctx context.Context, key string, value string) {
	if ctx.Err() != nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = &domain.Item{
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

func (m *MemoryStorage) Persist(ctx context.Context, key string, seconds int) bool {
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
	panic("unimplemented")
}

func (m *MemoryStorage) Exists(ctx context.Context, key string) bool {
	panic("unimplemented")
}

func (m *MemoryStorage) Size(ctx context.Context) int {
	panic("unimplemented")
}

func (m *MemoryStorage) StartCleanUp(intervalMs int64) {
	panic("unimplemented")
}

func (m *MemoryStorage) StopCleanUp() {
	panic("unimplemented")
}
