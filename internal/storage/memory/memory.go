package storage

import "sync"

type MemStorage struct {
	links      map[string]string
	linksMutex sync.RWMutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		links: make(map[string]string),
	}
}

func (m *MemStorage) Set(key string, value string) {
	m.linksMutex.Lock()
	defer m.linksMutex.Unlock()
	m.links[key] = value
}

func (m *MemStorage) Get(key string) (string, bool) {
	m.linksMutex.RLock()
	defer m.linksMutex.RUnlock()
	val, ok := m.links[key]
	return val, ok
}
